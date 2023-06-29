package group

import (
	"github.com/isd-sgcu/rpkm66-gateway/internal/dto"
	"github.com/isd-sgcu/rpkm66-gateway/internal/service/group"
	"github.com/isd-sgcu/rpkm66-gateway/proto"
)

type Service interface {
	FindOne(string) (*proto.Group, *dto.ResponseErr)
	FindByToken(string) (*proto.FindByTokenGroupResponse, *dto.ResponseErr)
	Update(*dto.GroupDto, string) (*proto.Group, *dto.ResponseErr)
	Join(string, string) (*proto.Group, *dto.ResponseErr)
	DeleteMember(string, string) (*proto.Group, *dto.ResponseErr)
	Leave(string) (*proto.Group, *dto.ResponseErr)
	SelectBaan(string, []string) (bool, *dto.ResponseErr)
}

func NewService(client proto.GroupServiceClient) Service {
	return group.NewService(client)
}
