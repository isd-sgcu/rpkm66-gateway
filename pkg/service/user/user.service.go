package user

import (
	"github.com/isd-sgcu/rpkm66-gateway/internal/dto"
	"github.com/isd-sgcu/rpkm66-gateway/internal/service/user"
	proto "github.com/isd-sgcu/rpkm66-go-proto/rpkm66/backend/user/v1"
)

type Service interface {
	FindOne(string) (*proto.User, *dto.ResponseErr)
	Create(*dto.UserDto) (*proto.User, *dto.ResponseErr)
	Update(string, *dto.UpdateUserDto) (*proto.User, *dto.ResponseErr)
	CreateOrUpdate(*dto.UserDto) (*proto.User, *dto.ResponseErr)
	Delete(string) (bool, *dto.ResponseErr)
	Verify(string, string) (bool, *dto.ResponseErr)
	GetUserEstamp(string) (*proto.GetUserEstampResponse, *dto.ResponseErr)
	ConfirmEstamp(string, string) (*proto.ConfirmEstampResponse, *dto.ResponseErr)
}

func NewService(client proto.UserServiceClient) Service {
	return user.NewService(client)
}
