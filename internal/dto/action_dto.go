package dto

import "time"

type ActionResponse struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	IsActive  bool      `json:"is_active"`
	HexCode   string    `json:"hex_code"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AvailableTicketActionResponse struct {
	ActionName     string  `json:"action_name"`
	HexCode        *string `json:"hex_code"`
	RequiresReason bool    `json:"requires_reason"`
	ReasonLabel    *string `json:"reason_label"`
	RequiresFile   bool    `json:"requires_file"`
}

type TransitionDetail struct {
	RequiredActorRole string
	Action            AvailableTicketActionResponse
}
