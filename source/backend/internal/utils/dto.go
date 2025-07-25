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

	// Used for joining an existing business (optional)
	BusinessID string `json:"businessId" validate:"omitempty,uuid"`

	// Used for creating a new business
	BusinessName string `json:"businessName" validate:"omitempty,min=2,max=64"`
	BusinessType string `json:"businessType" validate:"omitempty,min=2,max=64"` // maps to Type

	BusinessEmail string `json:"businessEmail" validate:"omitempty,email"`
	Phone         string `json:"phone" validate:"omitempty,min=3,max=32"`
	CountryCode   string `json:"countryCode" validate:"omitempty,max=8"`
	CompanySize   string `json:"companySize" validate:"omitempty,max=64"`
	Country       string `json:"country" validate:"omitempty,max=64"`
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
	EquipmentID string `json:"equipmentId" validate:"required"`
}

// ─────────────────────────────────────────────
// Pending-related requests
// ─────────────────────────────────────────────

type ApproveJoinRequest struct {
	UserID     string `json:"userId"`
	BusinessID string `json:"businessId"`
}

type InviteParams struct {
	BusinessID string
	Token      string
	Email      string
	Expiry     string
	Signature  string
}
