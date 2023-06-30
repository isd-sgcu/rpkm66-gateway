package file

import (
	fileConst "github.com/isd-sgcu/rpkm66-gateway/constant/file"
	"github.com/isd-sgcu/rpkm66-gateway/internal/dto"
	"github.com/isd-sgcu/rpkm66-gateway/internal/service/file"
	proto "github.com/isd-sgcu/rpkm66-go-proto/rpkm66/file/file/v1"
)

type Service interface {
	Upload(*dto.DecomposedFile, string, fileConst.Tag, fileConst.Type) (string, *dto.ResponseErr)
}

func NewService(client proto.FileServiceClient) Service {
	return file.NewService(client)
}
