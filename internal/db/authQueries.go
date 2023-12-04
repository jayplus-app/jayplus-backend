package db

import (
	"backend/models"
	"context"
	"encoding/json"
	"time"
)

// GetUserByEmail retrieves a user by email.
func (db *DB) GetUserByEmail(email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbConnectTimeout*time.Second)
	defer cancel()

	quey := `SELECT 
				*
			FROM 
				users
			WHERE
				email = $1`

	var user models.User

	row := db.QueryRowContext(ctx, quey, email)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.PhoneNumber,
		&user.HashedPassword,
		&user.FirstName,
		&user.LastName,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetUserByID retrieves a user by ID.
func (db *DB) GetUserByID(id int) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbConnectTimeout*time.Second)
	defer cancel()

	quey := `SELECT 
				*
			FROM
				users
			WHERE
				id = $1`

	var user models.User

	row := db.QueryRowContext(ctx, quey, id)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.PhoneNumber,
		&user.HashedPassword,
		&user.FirstName,
		&user.LastName,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// IsUserInBusiness checks if a user is in a business.
func (db *DB) IsUserInBusiness(userID, businessID int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbConnectTimeout*time.Second)
	defer cancel()

	query := `SELECT EXISTS(SELECT 1 FROM business_users WHERE user_id = $1 AND business_id = $2)`

	var exists bool
	err := db.QueryRowContext(ctx, query, userID, businessID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

// GetRoleNameByBusinessIDAndUserID retrieves a role name by business ID and user ID.
func (db *DB) GetRoleByBusinessIDAndUserID(businessID, userID int) (*models.Role, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbConnectTimeout*time.Second)
	defer cancel()

	query := `SELECT
				r.id,
				r.name
			FROM
				roles r
			INNER JOIN
				business_users bu ON bu.role_id = r.id
			WHERE
				bu.business_id = $1 AND
				bu.user_id = $2`

	row := db.QueryRowContext(ctx, query, businessID, userID)

	var role models.Role

	err := row.Scan(
		&role.ID,
		&role.Name,
	)
	if err != nil {
		return nil, err
	}

	return &role, nil
}

// GetBusinessHoursByBusinessID retrieves business hours by business ID.
func (db *DB) GetBusinessHoursByBusinessID(businessID int) (*models.BusinessHours, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbConnectTimeout*time.Second)
	defer cancel()

	query := `SELECT
            value
        FROM
            business_config
        WHERE
            business_id = $1 AND key = 'business-hours'`

	var jsonData []byte

	row := db.QueryRowContext(ctx, query, businessID)

	err := row.Scan(&jsonData)
	if err != nil {
		return nil, err
	}

	var businessHours models.BusinessHours
	if err := json.Unmarshal(jsonData, &businessHours); err != nil {
		return nil, err
	}

	return &businessHours, nil
}
