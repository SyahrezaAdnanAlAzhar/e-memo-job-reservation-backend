package dto

type AssignPICRequest struct {
	PicJobNPK string `json:"pic_job_npk" binding:"required"`
}

type ReorderJobsRequest struct {
	DepartmentTargetID int              `json:"department_target_id" binding:"required"`
	Items              []ReorderJobItem `json:"items" binding:"required,min=1"`
}

type ReorderJobItem struct {
	JobID   int `json:"job_id" binding:"required"`
	Version int `json:"version" binding:"required"`
}
