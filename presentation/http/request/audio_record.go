package request

type StartRecordRequestForm struct {
	Address string `form:"address" binding:"required"`
	seconds string `form:"seconds"  binding:"required"`
}
