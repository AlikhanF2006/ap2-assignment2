package provider

import (
	"context"
	"errors"
	"log"
	"math/rand"
	"time"
)

type SimulatedEmailSender struct {
	latencyMS   int
	failureRate int
}

func NewSimulatedEmailSender(latencyMS int, failureRate int) *SimulatedEmailSender {
	return &SimulatedEmailSender{
		latencyMS:   latencyMS,
		failureRate: failureRate,
	}
}

func (s *SimulatedEmailSender) Send(ctx context.Context, to string, subject string, body string) error {
	log.Printf("[Provider] Simulating external email provider latency=%dms", s.latencyMS)

	select {
	case <-time.After(time.Duration(s.latencyMS) * time.Millisecond):
	case <-ctx.Done():
		return ctx.Err()
	}

	if rand.Intn(100) < s.failureRate {
		return errors.New("simulated external provider failure")
	}

	log.Printf("[Provider] Email sent successfully to=%s subject=%s body=%s", to, subject, body)
	return nil
}
