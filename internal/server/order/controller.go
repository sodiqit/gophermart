package order

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sodiqit/gophermart/internal/logger"
	"github.com/sodiqit/gophermart/internal/server/auth"
	"github.com/sodiqit/gophermart/pkg/luhn"
)

type OrderController struct {
	logger       logger.Logger
	tokenService auth.TokenService
	orderService OrderService
}

func (c *OrderController) Route() *chi.Mux {
	r := chi.NewRouter()

	r.Use(auth.JWTAuth(c.tokenService))

	r.With(middleware.AllowContentType("text/plain")).Post("/", c.handleUploadOrder)
	r.Get("/", c.handleGetUserList)

	return r
}

func (c *OrderController) handleUploadOrder(w http.ResponseWriter, r *http.Request) {
	op := "orderController.handleUploadOrder"

	logger := c.logger.With("op", op)

	user := auth.ExtractUserFromContext(r.Context())

	b, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	body := string(b)

	isValidLuhnString := luhn.ValidateString(body)

	if !isValidLuhnString {
		http.Error(w, "Invalid order number", http.StatusUnprocessableEntity)
		return
	}

	err = c.orderService.Upload(r.Context(), user.ID, body)
	mapUploadResultToHttpAnswer(w, err, logger)
}

func (c *OrderController) handleGetUserList(w http.ResponseWriter, r *http.Request) {
	op := "orderController.handleGetUserList"

	logger := c.logger.With("op", op)

	user := auth.ExtractUserFromContext(r.Context())

	orders, err := c.orderService.GetUserOrders(r.Context(), user.ID)

	if err != nil {
		logger.Errorw("error while get user order list", "err", err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	result, err := json.Marshal(orders)

	if err != nil {
		logger.Errorw("error while serialize to json", "err", err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")

	if len(orders) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Write(result)
}

func NewController(logger logger.Logger, tokenService auth.TokenService, orderService OrderService) *OrderController {
	return &OrderController{
		logger,
		tokenService,
		orderService,
	}
}

func mapUploadResultToHttpAnswer(w http.ResponseWriter, err error, logger logger.Logger) {
	if errors.Is(err, ErrUserAlreadyUploadOrder) {
		w.WriteHeader(http.StatusOK)
		return
	}

	if errors.Is(err, ErrOrderAlreadyUploadByAnotherUser) {
		http.Error(w, "Order already uploaded by another user", http.StatusConflict)
		return
	}

	if err != nil {
		logger.Errorw("unexpected error while upload order", "err", err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
