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

// handleGetUserBalance godoc
//
//	@Summary		get balance
//	@Description	get total user balance
//	@Tags			balance
//
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Success		200	{object}	dtos.Balance
//	@Failure		401
//	@Failure		500
//	@Router			/api/user/balance [get]
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

// handleGetUserBalance godoc
//
//	@Summary		get withdrawals
//	@Description	get user withdrawals
//	@Tags			balance
//
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Success		200	{array}	dtos.Withdraw
//	@Failure		401
//	@Failure		500
//	@Router			/api/user/withdrawals [get]
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

// handleWithdraw godoc
//
//	@Summary		create withdraw
//	@Description	Process new withdraw request
//	@Tags			balance
//
//	@Param			body	body	WithdrawRequestDTO	true	"Withdraw body"
//	@Security		ApiKeyAuth
//	@Accept			json
//	@Produce		json
//	@Success		200
//	@Failure		400
//	@Failure		401
//	@Failure		402	string	true	"Not enough balance"
//	@Failure		422	string	true	"Not correct order number"
//	@Failure		500
//	@Router			/api/user/balance/withdraw [post]
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
