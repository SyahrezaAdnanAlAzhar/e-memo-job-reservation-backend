package model

import (
	"database/sql"
	"time"

	"github.com/lib/pq"
)

type Job struct {
	ID                   int            `json:"id"`
	TicketID             int            `json:"ticket_id"`
	PicJob               sql.NullString `json:"pic_job"`
	AssignedDepartmentID int            `json:"assigned_department_id"`
	JobPriority          int            `json:"job_priority"`
	ReportFile           pq.StringArray `json:"report_file"`
	CreatedAt            time.Time      `json:"created_at"`
	UpdatedAt            time.Time      `json:"updated_at"`
}
