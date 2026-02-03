package dto

import "github.com/zetsux/gin-gorm-api-starter/support/base"

// UserMaintenanceRequest represents a complex, highly customizable operation
// over a filtered set of users. It embeds UserGetsRequest so callers can
// reuse the existing filtering, search and pagination options.
type UserMaintenanceRequest struct {
	// Reuse existing filters / search / pagination for selecting target users.
	UserGetsRequest

	// Business rules / operations:

	// InactiveDays: only process users whose last update was at least this many days ago.
	// If 0 or negative, inactivity is ignored.
	InactiveDays int `json:"inactive_days" form:"inactive_days"`

	// NewRole: if non-empty, set this role for all selected users.
	NewRole string `json:"new_role" form:"new_role"`

	// ClearPicture: if true, clear profile pictures for processed users.
	ClearPicture bool `json:"clear_picture" form:"clear_picture"`
}

// UserMaintenanceResult summarizes what happened to a single user.
type UserMaintenanceResult struct {
	UserID         string `json:"user_id"`
	OldRole        string `json:"old_role,omitempty"`
	NewRole        string `json:"new_role,omitempty"`
	PictureCleared bool   `json:"picture_cleared"`
}

// UserMaintenanceResponse summarizes the overall maintenance run.
type UserMaintenanceResponse struct {
	TotalSelected       int                     `json:"total_selected"`
	TotalProcessed      int                     `json:"total_processed"`
	RoleChangedCount    int                     `json:"role_changed_count"`
	PictureClearedCount int                     `json:"picture_cleared_count"`
	Details             []UserMaintenanceResult `json:"details"`
	base.PaginationResponse
}
