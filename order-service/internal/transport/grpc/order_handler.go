package grpc

import (
	"time"

	"order-service/internal/usecase"

	orderpb "github.com/AlikhanF2006/ap2-protos-gen/order"

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

	for {
		order, err := h.usecase.GetOrderByID(orderID)
		if err != nil {
			return status.Errorf(codes.Internal, "failed to get order: %v", err)
		}

		update := &orderpb.OrderStatusUpdate{
			OrderId:   order.ID,
			Status:    order.Status,
			UpdatedAt: timestamppb.Now(),
		}

		if err := stream.Send(update); err != nil {
			return err
		}

		time.Sleep(2 * time.Second)
	}
}
