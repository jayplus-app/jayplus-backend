package db

import (
	"backend/models"
	"context"
	"time"
)

// RecordPayment records a payment.
func (db *DB) RecordPayment(payment *models.Payment) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbConnectTimeout*time.Second)
	defer cancel()

	query := `INSERT INTO payments (
				amount,
				description
			) VALUES (
				$1,
				$2
			) RETURNING id`

	err := db.QueryRowContext(
		ctx,
		query,
		payment.Amount,
		payment.Description,
	).Scan(&payment.ID)
	if err != nil {
		return err
	}

	return nil
}

// GetTransactionByBookingID retrieves a transaction by Booking ID.
func (db *DB) GetTransactionByBookingID(transactionID int) (*models.Transaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbConnectTimeout*time.Second)
	defer cancel()

	query := `SELECT
				id,
				booking_id,
				amount,
				gateway,
				status,
				notes
			FROM
				transactions
			WHERE
				id = $1`

	row := db.QueryRowContext(ctx, query, transactionID)

	var transaction models.Transaction

	err := row.Scan(
		&transaction.ID,
		&transaction.BookingID,
		&transaction.Amount,
		&transaction.Gateway,
		&transaction.Status,
		&transaction.Notes,
	)
	if err != nil {
		return nil, err
	}

	return &transaction, nil
}
