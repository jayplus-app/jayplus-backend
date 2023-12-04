package models

import "time"

type RolePermission struct {
	RoleID       int       `json:"roleId"`
	PermissionID int       `json:"permissionId"`
	CreatedAt    time.Time `json:"createdAt"`
}
