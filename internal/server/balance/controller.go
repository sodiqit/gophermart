package balance

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sodiqit/gophermart/internal/logger"
	"github.com/sodiqit/gophermart/internal/server/auth"
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
