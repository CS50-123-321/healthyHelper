package api
type ProgresRequest struct {
    TeleID int `form:"teleID" binding:"required"`
}