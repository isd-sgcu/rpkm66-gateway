package vaccine

import (
	"github.com/isd-sgcu/rnkm65-gateway/src/app/dto"
	"github.com/isd-sgcu/rnkm65-gateway/src/proto"
	"github.com/rs/zerolog/log"
	"net/http"
)

type Service struct {
	userService IUserService
	client      IClient
}

type IUserService interface {
	FindOne(string) (*proto.User, *dto.ResponseErr)
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
	user, errRes := s.userService.FindOne(userId)
	if errRes != nil {

		log.Error().
			Interface("error", errRes).
			Str("service", "vaccine").
			Str("module", "verify").
			Str("user_id", userId).
			Msg("Error while verifying the user")

		return nil, errRes
	}

	res := &dto.VaccineResponse{}

	err := s.client.Verify(&dto.VaccineRequest{
		HCert:     hcert,
		StudentId: user.StudentID,
	}, res)

	if err != nil {

		log.Error().
			Err(err).
			Str("service", "vaccine").
			Str("module", "verify").
			Str("student_id", user.StudentID).
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
			Str("student_id", user.StudentID).
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

	log.Info().
		Str("service", "vaccine").
		Str("module", "verify").
		Str("student_id", user.StudentID).
		Msg("Verified the vaccine cert")

	return res, nil
}
