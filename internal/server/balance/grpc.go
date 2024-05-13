package balance

import (
	"context"
	"errors"

	"github.com/bufbuild/protovalidate-go"
	proto "github.com/sodiqit/gophermart/gen/proto/balance/v1"
	"github.com/sodiqit/gophermart/internal/logger"
	"github.com/sodiqit/gophermart/internal/server/auth"
	"github.com/sodiqit/gophermart/pkg/luhn"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type BalanceServer struct {
	proto.UnimplementedBalanceServiceServer
	logger         logger.Logger
	balanceService BalanceService
	validator      *protovalidate.Validator
}

func (s *BalanceServer) GetBalance(ctx context.Context, in *proto.GetBalanceRequest) (*proto.GetBalanceResponse, error) {
	var response proto.GetBalanceResponse

	logger := s.logger.With("op", proto.BalanceService_GetBalance_FullMethodName)

	user := auth.ExtractUserFromContext(ctx)

	balance, err := s.balanceService.GetTotalBalance(ctx, user.ID)

	if err != nil {
		logger.Errorw("failed to get balance", "error", err)
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	response.Balance = &proto.Balance{
		Current:   balance.Current,
		Withdrawn: balance.Withdrawn,
	}

	return &response, nil
}

func (s *BalanceServer) GetWithdrawals(ctx context.Context, in *proto.GetWithdrawalsRequest) (*proto.GetWithdrawalsResponse, error) {
	var response proto.GetWithdrawalsResponse

	logger := s.logger.With("op", proto.BalanceService_GetWithdrawals_FullMethodName)

	user := auth.ExtractUserFromContext(ctx)

	withdrawals, err := s.balanceService.GetWithdrawals(ctx, user.ID)

	if err != nil {
		logger.Errorw("failed to get balance", "error", err)
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	result := make([]*proto.Withdrawal, 0, len(withdrawals))

	for _, withdraw := range withdrawals {
		result = append(result, &proto.Withdrawal{
			OrderId:     withdraw.OrderID,
			ProcessedAt: timestamppb.New(withdraw.ProcessedAt),
			Amount:      withdraw.Amount,
		})
	}

	response.Withdrawals = result

	return &response, nil
}

func (s *BalanceServer) Withdraw(ctx context.Context, in *proto.WithdrawRequest) (*proto.WithdrawResponse, error) {
	var response proto.WithdrawResponse

	logger := s.logger.With("op", proto.BalanceService_Withdraw_FullMethodName)

	user := auth.ExtractUserFromContext(ctx)

	isValidLuhnString := luhn.ValidateString(in.OrderId)

	if !isValidLuhnString {
		return nil, status.Error(codes.InvalidArgument, "Invalid order id")
	}

	err := s.balanceService.Withdraw(ctx, user.ID, in.OrderId, in.Sum)

	if err != nil && errors.Is(err, ErrInsufficientFunds) {
		return nil, status.Error(codes.InvalidArgument, "Not enough funds")
	}

	if err != nil {
		logger.Errorw("failed to withdraw", "error", err)
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	return &response, nil
}

func NewBalanceServer(logger logger.Logger, balanceService BalanceService) *BalanceServer {
	v, err := protovalidate.New()
	if err != nil {
		panic(err)
	}
	return &BalanceServer{
		logger:         logger,
		balanceService: balanceService,
		validator:      v,
	}
}
