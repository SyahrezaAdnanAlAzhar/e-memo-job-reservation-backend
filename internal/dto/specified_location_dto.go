package dto

type CreateSpecifiedLocationRequest struct {
	PhysicalLocationID int    `json:"physical_location_id" binding:"required,gt=0"`
	Name               string `json:"name" binding:"required"`
}