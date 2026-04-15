package repository

import (
	"database/sql"
	"errors"

	"payment-service/internal/domain"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{
		db: db,
	}
}

func (r *PostgresRepository) Create(payment *domain.Payment) error {
	query := `
		INSERT INTO payments (id, order_id, transaction_id, amount, status)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.Exec(
		query,
		payment.ID,
		payment.OrderID,
		payment.TransactionID,
		payment.Amount,
		payment.Status,
	)

	return err
}

func (r *PostgresRepository) GetByOrderID(orderID string) (*domain.Payment, error) {
	query := `
		SELECT id, order_id, transaction_id, amount, status
		FROM payments
		WHERE order_id = $1
	`

	row := r.db.QueryRow(query, orderID)

	var payment domain.Payment
	err := row.Scan(
		&payment.ID,
		&payment.OrderID,
		&payment.TransactionID,
		&payment.Amount,
		&payment.Status,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &payment, nil
}
func (r *PostgresRepository) ListByStatus(status string) ([]*domain.Payment, error) {
	var (
		rows *sql.Rows
		err  error
	)

	if status == "" {
		rows, err = r.db.Query(`
			SELECT id, order_id, transaction_id, amount, status
			FROM payments
		`)
	} else {
		rows, err = r.db.Query(`
			SELECT id, order_id, transaction_id, amount, status
			FROM payments
			WHERE status = $1
		`, status)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []*domain.Payment

	for rows.Next() {
		var p domain.Payment
		if err := rows.Scan(
			&p.ID,
			&p.OrderID,
			&p.TransactionID,
			&p.Amount,
			&p.Status,
		); err != nil {
			return nil, err
		}

		payments = append(payments, &p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return payments, nil
}
