package file

import (
	"context"
	"time"

	"github.com/isd-sgcu/rpkm66-gateway/constant/file"
	"github.com/isd-sgcu/rpkm66-gateway/internal/dto"
	"github.com/isd-sgcu/rpkm66-gateway/internal/utils"
	"github.com/isd-sgcu/rpkm66-gateway/proto"
	"github.com/rs/zerolog/log"
)

type serviceImpl struct {
	client proto.FileServiceClient
}

func NewService(client proto.FileServiceClient) *serviceImpl {
	return &serviceImpl{client: client}
}

func (s *serviceImpl) Upload(file *dto.DecomposedFile, userId string, tag file.Tag, fileType file.Type) (string, *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := s.client.Upload(ctx, &proto.UploadRequest{
		Filename: file.Filename,
		Data:     file.Data,
		UserId:   userId,
		Tag:      int32(tag),
		Type:     int32(fileType),
	})

	if err != nil {
		log.Error().
			Err(err).
			Str("service", "file").
			Str("method", "upload file").
			Msg("Error while upload file")

		return "", utils.ServiceErrorHandler(err)
	}

	return res.Url, nil
}
