package vaccine

import (
	"github.com/isd-sgcu/rnkm65-gateway/src/app/dto"
	"github.com/rs/zerolog/log"
	"net/http"
)

type Service struct {
	userService IUserService
	client      IClient
}

type IUserService interface {
	Verify(string) (bool, *dto.ResponseErr)
}

type IClient interface {
	Verify(*dto.VaccineRequest, *dto.VaccineResponse) error
}

func NewService(userService IUserService, client IClient) *Service {
	return &Service{
		userService: userService,
		client:      client,
	}
}

func (s *Service) Verify(hcert string, userId string) (*dto.VaccineResponse, *dto.ResponseErr) {
	res := &dto.VaccineResponse{}

	err := s.client.Verify(&dto.VaccineRequest{
		HCert:     hcert,
		StudentId: userId,
	}, res)

	if err != nil {

		log.Error().
			Err(err).
			Str("service", "vaccine").
			Str("module", "verify").
			Msg("Cannot verify the QR code")

		return nil, &dto.ResponseErr{
			StatusCode: http.StatusBadRequest,
			Message:    "Cannot verify the qr code",
			Data:       nil,
		}
	}

	ok, errRes := s.userService.Verify(res.Uid)
	if errRes != nil {

		log.Error().Interface("error", errRes).
			Str("service", "vaccine").
			Str("module", "verify").
			Msg("Error while update the user status")

		return nil, errRes
	}

	if !ok {
		return nil, &dto.ResponseErr{
			StatusCode: http.StatusNotFound,
			Message:    "Not found user",
			Data:       nil,
		}
	}

	return res, nil
}
