package order

import (
	"context"
	"errors"

	"github.com/bufbuild/protovalidate-go"
	proto "github.com/sodiqit/gophermart/gen/proto/order/v1"
	"github.com/sodiqit/gophermart/internal/logger"
	"github.com/sodiqit/gophermart/internal/server/auth"
	"github.com/sodiqit/gophermart/internal/server/repository"
	"github.com/sodiqit/gophermart/pkg/luhn"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type OrderServer struct {
	proto.UnimplementedOrderServiceServer
	logger       logger.Logger
	orderService OrderService
	validator    *protovalidate.Validator
}

func (s *OrderServer) Upload(ctx context.Context, in *proto.UploadRequest) (*proto.UploadResponse, error) {
	var response proto.UploadResponse

	logger := s.logger.With("op", proto.OrderService_Upload_FullMethodName)

	err := s.validator.Validate(in)

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	user := auth.ExtractUserFromContext(ctx)

	isValidLuhnString := luhn.ValidateString(in.OrderId)

	if !isValidLuhnString {
		return nil, status.Error(codes.InvalidArgument, "invalid order id")
	}

	err = s.orderService.Upload(ctx, user.ID, in.OrderId)

	if errors.Is(err, ErrUserAlreadyUploadOrder) || err == nil {
		return &response, nil
	}

	return nil, mapUploadServiceError(err, logger)
}

func (s *OrderServer) GetList(ctx context.Context, in *proto.GetListRequest) (*proto.GetListResponse, error) {
	var response proto.GetListResponse

	logger := s.logger.With("op", proto.OrderService_GetList_FullMethodName)

	err := s.validator.Validate(in)

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	user := auth.ExtractUserFromContext(ctx)

	orders, err := s.orderService.GetUserOrders(ctx, user.ID)

	if err != nil {
		logger.Errorw("failed to get orders", "err", err)
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	result := make([]*proto.Order, 0, len(orders))

	for _, order := range orders {
		result = append(result, &proto.Order{
			Number:     order.ID,
			Accrual:    order.Accrual,
			UploadedAt: timestamppb.New(order.CreatedAt),
			Status:     mapOrderStatusToProto(order.Status),
		})
	}

	response.Orders = result

	return &response, nil
}

func mapOrderStatusToProto(status string) proto.Order_OrderStatus {
	switch status {
	case repository.OrderStatusNew:
		return proto.Order_NEW
	case repository.OrderStatusProcessing:
		return proto.Order_PROCESSING
	case repository.OrderStatusInvalid:
		return proto.Order_INVALID
	case repository.OrderStatusProcessed:
		return proto.Order_PROCESSED
	}

	panic("invalid order status")
}

func mapUploadServiceError(err error, logger logger.Logger) error {
	code := codes.Internal
	msg := "Internal server error"

	if errors.Is(err, ErrOrderAlreadyUploadByAnotherUser) {
		code = codes.AlreadyExists
		msg = err.Error()
	}

	return status.Error(code, msg)
}

func NewOrderServer(logger logger.Logger, orderService OrderService) *OrderServer {
	v, err := protovalidate.New()
	if err != nil {
		panic(err)
	}
	return &OrderServer{
		logger:       logger,
		orderService: orderService,
		validator:    v,
	}
}
