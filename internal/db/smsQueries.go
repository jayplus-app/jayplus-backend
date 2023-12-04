package db

import (
	"backend/models"
	"context"
	"time"
)

// RecordSMS records a sms.
func (db *DB) RecordSMS(sms *models.SMS) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbConnectTimeout*time.Second)
	defer cancel()

	query := `INSERT INTO sms (
				phone,
				message
			) VALUES (
				$1,
				$2
			) RETURNING id`

	err := db.QueryRowContext(
		ctx,
		query,
		sms.ToNumber,
		sms.Content,
	).Scan(&sms.ID)
	if err != nil {
		return err
	}

	return nil
}
