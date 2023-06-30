package checkin

import (
	"github.com/isd-sgcu/rpkm66-gateway/internal/dto"
	"github.com/isd-sgcu/rpkm66-gateway/internal/service/checkin"
	proto "github.com/isd-sgcu/rpkm66-go-proto/rpkm66/backend/checkin/v1"
)

type Service interface {
	CheckinVerify(string, int) (*proto.CheckinVerifyResponse, *dto.ResponseErr)
	CheckinConfirm(token string) (*proto.CheckinConfirmResponse, *dto.ResponseErr)
}

func NewService(client proto.CheckinServiceClient) Service {
	return checkin.NewService(client)
}
