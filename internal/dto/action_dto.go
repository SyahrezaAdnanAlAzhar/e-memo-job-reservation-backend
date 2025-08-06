package dto

type ActionResponse struct {
	ActionName     string  `json:"action_name"`
	HexCode        *string `json:"hex_code"`
	RequiresReason bool    `json:"requires_reason"`
	ReasonLabel    *string `json:"reason_label"`
	RequiresFile   bool    `json:"requires_file"`
}
