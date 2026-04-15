package grpc

import (
	"time"

	orderpb "github.com/AlikhanF2006/ap2-protos-gen/order"

	"order-service/internal/usecase"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type OrderHandler struct {
	orderpb.UnimplementedOrderServiceServer
	usecase *usecase.OrderUsecase
}

func NewOrderHandler(u *usecase.OrderUsecase) *OrderHandler {
	return &OrderHandler{usecase: u}
}

func (h *OrderHandler) SubscribeToOrderUpdates(
	req *orderpb.OrderRequest,
	stream orderpb.OrderService_SubscribeToOrderUpdatesServer,
) error {
	if req.GetOrderId() == "" {
		return status.Error(codes.InvalidArgument, "order_id is required")
	}

	orderID := req.GetOrderId()

	order, err := h.usecase.GetOrderByID(orderID)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to get order: %v", err)
	}
	if order == nil {
		return status.Error(codes.NotFound, "order not found")
	}

	lastStatus := order.Status

	if err := stream.Send(&orderpb.OrderStatusUpdate{
		OrderId:   order.ID,
		Status:    order.Status,
		UpdatedAt: timestamppb.Now(),
	}); err != nil {
		return status.Errorf(codes.Internal, "failed to send initial update: %v", err)
	}

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-stream.Context().Done():
			return nil

		case <-ticker.C:
			order, err := h.usecase.GetOrderByID(orderID)
			if err != nil {
				return status.Errorf(codes.Internal, "failed to get order: %v", err)
			}
			if order == nil {
				return status.Error(codes.NotFound, "order not found")
			}

			if order.Status != lastStatus {
				lastStatus = order.Status

				if err := stream.Send(&orderpb.OrderStatusUpdate{
					OrderId:   order.ID,
					Status:    order.Status,
					UpdatedAt: timestamppb.Now(),
				}); err != nil {
					return status.Errorf(codes.Internal, "failed to send update: %v", err)
				}
			}
		}
	}
}
