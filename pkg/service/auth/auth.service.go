package auth

import (
	"github.com/isd-sgcu/rpkm66-gateway/internal/dto"
	"github.com/isd-sgcu/rpkm66-gateway/internal/service/auth"
	proto "github.com/isd-sgcu/rpkm66-go-proto/rpkm66/auth/auth/v1"
)

type Service interface {
	VerifyTicket(string) (*proto.Credential, *dto.ResponseErr)
	Validate(string) (*dto.TokenPayloadAuth, *dto.ResponseErr)
	RefreshToken(string) (*proto.Credential, *dto.ResponseErr)
}

func NewService(client proto.AuthServiceClient) Service {
	return auth.NewService(client)
}
