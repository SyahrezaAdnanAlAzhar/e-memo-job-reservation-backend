package dto

type AssignPICRequest struct {
	PicJobNPK string `json:"pic_job_npk" binding:"required"`
}

type ReorderJobsRequest struct {
	DepartmentTargetID int   `json:"department_target_id" binding:"required"`
	OrderedJobIDs      []int `json:"ordered_job_ids" binding:"required"`
}