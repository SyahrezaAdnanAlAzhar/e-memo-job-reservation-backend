package model

import (
	"database/sql"
	"time"
)

type Ticket struct {
	ID                  int           `json:"id"`
	Requestor           string        `json:"requestor"`
	DepartmentTargetID  int           `json:"department_target_id"`
	PhysicalLocationID  sql.NullInt64 `json:"physical_location_id"`
	SpecifiedLocationID sql.NullInt64 `json:"specified_location_id"`
	Description         string        `json:"description"`
	TicketPriority      int           `json:"ticket_priority"`
	SupportFile         []string      `json:"support_file"`
	Version             int           `json:"version"`
	Deadline            sql.NullTime  `json:"deadline"`
	CreatedAt           time.Time     `json:"created_at"`
	UpdatedAt           time.Time     `json:"updated_at"`
}
