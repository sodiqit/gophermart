package order

import (
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

	w.WriteHeader(http.StatusCreated)
}
