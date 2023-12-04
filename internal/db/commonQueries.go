package db

import (
	"backend/models"
	"context"
	"time"
)

// GetBusinessByBusinessName retrieves a business by business name.
func (db *DB) GetBusinessByBusinessName(businessName string) (*models.Business, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbConnectTimeout*time.Second)
	defer cancel()

	query := `SELECT * FROM businesses WHERE business_name = $1`

	var business models.Business

	row := db.QueryRowContext(ctx, query, businessName)

	err := row.Scan(
		&business.ID,
		&business.Name,
		&business.BusinessName,
		&business.Timezone,
		&business.CreatedAt,
		&business.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &business, nil
}
