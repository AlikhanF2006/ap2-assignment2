package grpc

import (
	"context"

	"payment-service/internal/usecase"

	paymentpb "github.com/AlikhanF2006/ap2-protos-gen/payment"
)

type PaymentHandler struct {
	paymentpb.UnimplementedPaymentServiceServer
	uc *usecase.PaymentUsecase
}

func NewPaymentHandler(uc *usecase.PaymentUsecase) *PaymentHandler {
	return &PaymentHandler{uc: uc}
}

func (h *PaymentHandler) ProcessPayment(ctx context.Context, req *paymentpb.PaymentRequest) (*paymentpb.PaymentResponse, error) {

	payment, err := h.uc.CreatePayment(req.OrderId, int64(req.Amount))
	if err != nil {
		return &paymentpb.PaymentResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &paymentpb.PaymentResponse{
		Success:       true,
		TransactionId: payment.TransactionID,
		Message:       payment.Status,
	}, nil
}
