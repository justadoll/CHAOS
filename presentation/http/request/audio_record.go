package request

type StartRecordRequestForm struct {
	Address string `form:"address" binding:"required"`
	Seconds string `form:"seconds"  binding:"required"`
}
