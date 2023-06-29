package file

import (
	"context"
	"net/http"
	"time"

	"github.com/isd-sgcu/rpkm66-gateway/constant/file"
	"github.com/isd-sgcu/rpkm66-gateway/internal/dto"
	"github.com/isd-sgcu/rpkm66-gateway/proto"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.Unavailable:
				log.Error().
					Err(err).
					Str("service", "file").
					Str("module", "upload").
					Str("user_id", userId).
					Msg("Something wrong")
				return "", &dto.ResponseErr{
					StatusCode: http.StatusGatewayTimeout,
					Message:    "Connection timeout",
					Data:       nil,
				}

			default:
				log.Error().
					Err(err).
					Str("service", "file").
					Str("module", "upload").
					Str("user_id", userId).
					Msg("Error while connecting to service")
				return "", &dto.ResponseErr{
					StatusCode: http.StatusServiceUnavailable,
					Message:    "Service is down",
					Data:       nil,
				}
			}
		}

		log.Error().
			Err(err).
			Str("service", "file").
			Str("module", "upload").
			Str("user_id", userId).
			Msg("Error while connecting to service")

		return "", &dto.ResponseErr{
			StatusCode: http.StatusServiceUnavailable,
			Message:    "Service is down",
			Data:       nil,
		}

	}

	return res.Url, nil
}
