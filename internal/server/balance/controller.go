package balance

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sodiqit/gophermart/internal/logger"
	"github.com/sodiqit/gophermart/internal/server/auth"
	"github.com/sodiqit/gophermart/internal/utils"
	"github.com/sodiqit/gophermart/pkg/luhn"
)

type BalanceController struct {
	logger         logger.Logger
	tokenService   auth.TokenService
	balanceService BalanceService
}

func (c *BalanceController) Connect(r *chi.Mux, basePath string) {
	r.Group(func(r chi.Router) {
		r.Use(auth.JWTAuth(c.tokenService))

		r.Get(fmt.Sprintf("%suser/balance", basePath), c.handleGetUserBalance)
		r.Get(fmt.Sprintf("%suser/withdrawals", basePath), c.handleGetUserWithdrawals)
		r.With(middleware.AllowContentType("application/json")).Post(fmt.Sprintf("%suser/balance/withdraw", basePath), c.handleWithdraw)
	})
}

func (c *BalanceController) handleGetUserBalance(w http.ResponseWriter, r *http.Request) {
	op := "balanceController.handleGetUserBalance"

	logger := c.logger.With("op", op)

	user := auth.ExtractUserFromContext(r.Context())

	balance, err := c.balanceService.GetTotalBalance(r.Context(), user.ID)

	if err != nil {
		logger.Errorw("", "err", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	result, err := json.Marshal(balance)

	if err != nil {
		logger.Errorw("error while serialize balance", "err", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(result)
}

func (c *BalanceController) handleGetUserWithdrawals(w http.ResponseWriter, r *http.Request) {
	op := "balanceController.handleGetUserWithdrawals"

	logger := c.logger.With("op", op)

	user := auth.ExtractUserFromContext(r.Context())

	withdrawals, err := c.balanceService.GetWithdrawals(r.Context(), user.ID)

	if err != nil {
		logger.Errorw("", "err", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if len(withdrawals) == 0 {
		http.Error(w, "", http.StatusNoContent)
		return
	}

	result, err := json.Marshal(withdrawals)

	if err != nil {
		logger.Errorw("error while serialize withdrawals", "err", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(result)
}

func (c *BalanceController) handleWithdraw(w http.ResponseWriter, r *http.Request) {
	op := "balanceController.handleWithdraw"

	logger := c.logger.With("op", op)

	user := auth.ExtractUserFromContext(r.Context())

	var dto WithdrawRequestDTO

	err := utils.ValidateJSONBody(r.Context(), r.Body, &dto)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	isValidLuhnString := luhn.ValidateString(dto.OrderID)

	if !isValidLuhnString {
		http.Error(w, "Invalid order number", http.StatusUnprocessableEntity)
		return
	}

	err = c.balanceService.Withdraw(r.Context(), user.ID, dto.OrderID, dto.Sum)

	if err != nil && errors.Is(err, ErrInsufficientFunds) {
		http.Error(w, "Not have enough funds", http.StatusPaymentRequired)
		return
	}

	if err != nil {
		logger.Errorw("error while register withdraw", "err", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func NewController(logger logger.Logger, tokenService auth.TokenService, balanceService BalanceService) *BalanceController {
	return &BalanceController{
		logger,
		tokenService,
		balanceService,
	}
}

// func mapUploadResultToHttpAnswer(w http.ResponseWriter, err error, logger logger.Logger) {
// 	if errors.Is(err, ErrUserAlreadyUploadOrder) {
// 		w.WriteHeader(http.StatusOK)
// 		return
// 	}

// 	if errors.Is(err, ErrOrderAlreadyUploadByAnotherUser) {
// 		http.Error(w, "Order already uploaded by another user", http.StatusConflict)
// 		return
// 	}

// 	if err != nil {
// 		logger.Errorw("unexpected error while upload order", "err", err.Error())
// 		http.Error(w, "Internal server error", http.StatusInternalServerError)
// 		return
// 	}

// 	w.WriteHeader(http.StatusCreated)
// }
