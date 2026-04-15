package grpc

import (
	"context"

	paymentpb "github.com/AlikhanF2006/ap2-protos-gen/payment"

	"payment-service/internal/usecase"
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

func (h *PaymentHandler) ListPayments(ctx context.Context, req *paymentpb.ListPaymentsRequest) (*paymentpb.ListPaymentsResponse, error) {
	payments, err := h.uc.ListPayments(req.Status)
	if err != nil {
		return nil, err
	}

	resp := &paymentpb.ListPaymentsResponse{
		Payments: make([]*paymentpb.PaymentItem, 0, len(payments)),
	}

	for _, p := range payments {
		resp.Payments = append(resp.Payments, &paymentpb.PaymentItem{
			Id:            p.ID,
			OrderId:       p.OrderID,
			TransactionId: p.TransactionID,
			Amount:        p.Amount,
			Status:        p.Status,
		})
	}

	return resp, nil
}
