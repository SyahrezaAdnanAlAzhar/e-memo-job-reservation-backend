package model

import (
	"database/sql"
	"time"
)

type Job struct {
	ID          int            `json:"id"`
	TicketID    int            `json:"ticket_id"`
	PicJob      sql.NullString `json:"pic_job"`
	JobPriority int            `json:"job_priority"`
	ReportFile  []string       `json:"report_file"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}