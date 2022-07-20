package qr

import validate "github.com/isd-sgcu/rnkm65-gateway/src/app/validator"

type Handler struct {
	checkinService ICheckinService
	validate       *validate.DtoValidator
	// estampService  IEstampService
}

func NewHandler(checkinService ICheckinService, v *validate.DtoValidator) *Handler {
	return &Handler{
		checkinService: checkinService,
		validate:       v,
	}
}
