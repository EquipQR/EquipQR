package utils

// ─────────────────────────────────────────────
// User-related requests
// ─────────────────────────────────────────────

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type CreateUserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=64"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=256"`
}

// ─────────────────────────────────────────────
// Business-related requests
// ─────────────────────────────────────────────

type CreateBusinessRequest struct {
	BusinessName string `json:"businessName" validate:"required,min=2,max=64"`
}

// ─────────────────────────────────────────────
// Equipment-related requests
// ─────────────────────────────────────────────

type CreateEquipmentRequest struct {
	BusinessID string `json:"business_id" validate:"required"`
	Status     string `json:"status" validate:"required,oneof='in service' 'not in service'"`
	Type       string `json:"type" validate:"required"`
	Location   string `json:"location"`
	MoreFields any    `json:"more_fields"`
}

// ─────────────────────────────────────────────
// Issue-related requests
// ─────────────────────────────────────────────

type CreateIssueRequest struct {
	Title       string `json:"title" validate:"required,min=3,max=128"`
	Description string `json:"description" validate:"required"`
	EquipmentID string `json:"equipment_id" validate:"required"`
	AssigneeID  string `json:"assignee_id" validate:"required"`
}
