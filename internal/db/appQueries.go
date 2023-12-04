package db

import (
	"backend/models"
	"context"
	"encoding/json"
	"time"
)

// GetUIConfig retrieves the UI config.
func (db *DB) GetBusinessUIConfigByID(business_id int) (*models.UIConfig, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbConnectTimeout*time.Second)
	defer cancel()

	query := `SELECT
				value
			FROM
				business_config
			WHERE
				business_id = $1 AND key = 'ui-config'`

	var jsonData []byte

	row := db.QueryRowContext(ctx, query, business_id)

	err := row.Scan(&jsonData)
	if err != nil {
		return nil, err
	}

	var uiConfig models.UIConfig
	if err := json.Unmarshal(jsonData, &uiConfig); err != nil {
		return nil, err
	}

	return &uiConfig, nil
}

// GetBookingConfig retrieves the booking config.
func (db *DB) GetBusinessBookingConfigByID(business_id int) (*models.BookingConfig, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbConnectTimeout*time.Second)
	defer cancel()

	query := `SELECT
				value
			FROM
				business_config
			WHERE
				business_id = $1 AND key = 'booking-config'`

	var jsonData []byte

	row := db.QueryRowContext(ctx, query, business_id)

	err := row.Scan(&jsonData)
	if err != nil {
		return nil, err
	}

	var bookingConfig models.BookingConfig
	if err := json.Unmarshal(jsonData, &bookingConfig); err != nil {
		return nil, err
	}

	return &bookingConfig, nil
}
