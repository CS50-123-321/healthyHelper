package api

type ProgresRequest struct {
	TeleID int `json:"tele_id" binding:"required"`
}
