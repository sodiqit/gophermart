package accrual

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
	"golang.org/x/time/rate"
)

const tooManyRequestTemplate = "No more than %d requests per minute allowed"

var ErrOrderNotFound = errors.New("order not found")
var ErrOrderIncorrectResponseBody = errors.New("order response body not match")
var ErrRateLimit = errors.New("rate limit")

type OrderInfoDTO struct {
	OrderID string   `json:"order"`
	Accrual *float64 `json:"accrual,omitempty"`
	Status  string   `json:"status"`
}

type AccrualClient interface {
	GetOrderInfo(ctx context.Context, orderID string) (OrderInfoDTO, error)
}

type HTTPAccrualClient struct {
	endpointTemplate string
	limiter          *rate.Limiter
	limiterMutex     sync.Mutex
	httpClient       *resty.Client
	rl               int
}

func (c *HTTPAccrualClient) GetOrderInfo(ctx context.Context, orderID string) (OrderInfoDTO, error) {
	if c.limiter != nil {
		err := c.limiter.Wait(ctx)
		if err != nil {
			return OrderInfoDTO{}, err
		}
	}

	ctx, cancel := context.WithTimeout(ctx, time.Duration(5)*time.Second)

	defer cancel()

	request := c.httpClient.R()
	response, err := request.SetResult(OrderInfoDTO{}).SetContext(ctx).Get(fmt.Sprintf(c.endpointTemplate, orderID))

	if err != nil {
		return OrderInfoDTO{}, err
	}

	if response.StatusCode() == http.StatusTooManyRequests {
		var rl int
		_, err = fmt.Sscanf(response.String(), tooManyRequestTemplate, &rl)

		if err != nil {
			return OrderInfoDTO{}, err
		}

		c.limiterMutex.Lock()

		if c.rl != rl || c.limiter == nil {
			c.limiter = rate.NewLimiter(rate.Every(time.Minute/time.Duration(rl)), rl)
			c.rl = rl
		}

		c.limiterMutex.Unlock()
		return OrderInfoDTO{}, ErrRateLimit
	}

	if response.StatusCode() >= http.StatusBadRequest {
		return OrderInfoDTO{}, fmt.Errorf("failed get order info %s. StatusCode: %d. URL: %s", response.String(), response.StatusCode(), response.Request.URL)
	}

	if response.StatusCode() == http.StatusNoContent {
		return OrderInfoDTO{}, fmt.Errorf("%w: %s", ErrOrderNotFound, orderID)
	}

	result, ok := response.Result().(*OrderInfoDTO)

	if !ok {
		return OrderInfoDTO{}, fmt.Errorf("%w. Got: %s. Status code: %d", ErrOrderIncorrectResponseBody, response.String(), response.StatusCode())
	}

	return *result, nil
}

func NewHTTPAccrualClient(endpointTemplate string) *HTTPAccrualClient {
	return &HTTPAccrualClient{
		endpointTemplate: endpointTemplate,
		limiter:          nil,
		httpClient:       resty.New(),
		limiterMutex:     sync.Mutex{},
	}
}
