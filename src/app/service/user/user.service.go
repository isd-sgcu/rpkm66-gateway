package user

import (
	"context"
	"github.com/isd-sgcu/rnkm65-gateway/src/app/dto"
	"github.com/isd-sgcu/rnkm65-gateway/src/proto"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"time"
)

type Service struct {
	client proto.UserServiceClient
}

func NewService(client proto.UserServiceClient) *Service {
	return &Service{
		client: client,
	}
}

func (s *Service) FindOne(id string) (result *proto.User, err *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, errRes := s.client.FindOne(ctx, &proto.FindOneUserRequest{Id: id})
	if errRes != nil {
		st, ok := status.FromError(errRes)
		if ok {
			switch st.Code() {
			case codes.NotFound:
				return nil, &dto.ResponseErr{
					StatusCode: http.StatusNotFound,
					Message:    "User not found",
					Data:       nil,
				}
			default:

				log.Error().
					Err(errRes).
					Str("service", "user").
					Str("module", "findOne").
					Msg("Error while connecting to service")

				return nil, &dto.ResponseErr{
					StatusCode: http.StatusServiceUnavailable,
					Message:    "Service is down",
					Data:       nil,
				}
			}
		}

		log.Error().
			Err(errRes).
			Str("service", "user").
			Str("module", "findOne").
			Msg("Error while connecting to service")

		return nil, &dto.ResponseErr{
			StatusCode: http.StatusServiceUnavailable,
			Message:    "Service is down",
			Data:       nil,
		}
	}

	return res.User, nil
}

func (s *Service) Create(in *dto.UserDto) (result *proto.User, err *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	usrDto := &proto.User{
		Title:           in.Title,
		Firstname:       in.Firstname,
		Lastname:        in.Lastname,
		Nickname:        in.Nickname,
		Phone:           in.Phone,
		LineID:          in.LineID,
		Email:           in.Email,
		AllergyFood:     in.AllergyFood,
		FoodRestriction: in.FoodRestriction,
		AllergyMedicine: in.AllergyMedicine,
		Disease:         in.Disease,
		CanSelectBaan:   *in.CanSelectBaan,
	}

	res, errRes := s.client.Create(ctx, &proto.CreateUserRequest{User: usrDto})
	if errRes != nil {

		log.Error().
			Err(errRes).
			Str("service", "user").
			Str("module", "create").
			Msg("Error while connecting to service")

		return nil, &dto.ResponseErr{
			StatusCode: http.StatusServiceUnavailable,
			Message:    "Service is down",
			Data:       nil,
		}
	}

	return res.User, nil
}

func (s *Service) Update(id string, in *dto.UpdateUserDto) (result *proto.User, err *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	usrReq := &proto.UpdateUserRequest{
		Id:              id,
		Title:           in.Title,
		Firstname:       in.Firstname,
		Lastname:        in.Lastname,
		Nickname:        in.Nickname,
		Phone:           in.Phone,
		LineID:          in.LineID,
		Email:           in.Email,
		AllergyFood:     in.AllergyFood,
		FoodRestriction: in.FoodRestriction,
		AllergyMedicine: in.AllergyMedicine,
		Disease:         in.Disease,
	}

	res, errRes := s.client.Update(ctx, usrReq)
	if errRes != nil {
		st, ok := status.FromError(errRes)
		if ok {
			switch st.Code() {
			case codes.NotFound:
				return nil, &dto.ResponseErr{
					StatusCode: http.StatusNotFound,
					Message:    st.Message(),
					Data:       nil,
				}
			default:

				log.Error().
					Err(errRes).
					Str("service", "user").
					Str("module", "update").
					Msg("Error while connecting to service")

				return nil, &dto.ResponseErr{
					StatusCode: http.StatusServiceUnavailable,
					Message:    "Service is down",
					Data:       nil,
				}
			}
		}

		log.Error().
			Err(errRes).
			Str("service", "user").
			Str("module", "update").
			Msg("Error while connecting to service")

		return nil, &dto.ResponseErr{
			StatusCode: http.StatusServiceUnavailable,
			Message:    "Service is down",
			Data:       nil,
		}
	}

	return res.User, nil
}

func (s *Service) Verify(studentId string) (result bool, err *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Info().
		Str("service", "user").
		Str("module", "verify").
		Str("student_id", studentId).
		Msg("Trying to verify the user")

	res, errRes := s.client.Verify(ctx, &proto.VerifyUserRequest{StudentId: studentId})
	if errRes != nil {
		return false, &dto.ResponseErr{
			StatusCode: http.StatusNotFound,
			Message:    "User not found",
			Data:       nil,
		}
	}

	log.Info().
		Str("service", "user").
		Str("module", "verify").
		Str("student_id", studentId).
		Msg("Verified the user")

	return res.Success, nil
}

func (s *Service) CreateOrUpdate(in *dto.UserDto) (result *proto.User, err *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	usrDto := &proto.User{
		Id:              in.ID,
		Title:           in.Title,
		Firstname:       in.Firstname,
		Lastname:        in.Lastname,
		Nickname:        in.Nickname,
		Phone:           in.Phone,
		LineID:          in.LineID,
		Email:           in.Email,
		AllergyFood:     in.AllergyFood,
		FoodRestriction: in.FoodRestriction,
		AllergyMedicine: in.AllergyMedicine,
		Disease:         in.Disease,
		CanSelectBaan:   *in.CanSelectBaan,
	}

	res, errRes := s.client.CreateOrUpdate(ctx, &proto.CreateOrUpdateUserRequest{User: usrDto})
	if errRes != nil {
		st, ok := status.FromError(errRes)
		if ok {
			switch st.Code() {
			case codes.NotFound:
				return nil, &dto.ResponseErr{
					StatusCode: http.StatusNotFound,
					Message:    st.Message(),
					Data:       nil,
				}
			default:

				log.Error().
					Err(errRes).
					Str("service", "user").
					Str("module", "create and update").
					Str("student_id", usrDto.StudentID).
					Msg("Error while connecting to service")

				return nil, &dto.ResponseErr{
					StatusCode: http.StatusServiceUnavailable,
					Message:    "Service is down",
					Data:       nil,
				}
			}
		}

		log.Error().
			Err(errRes).
			Str("service", "user").
			Str("module", "create and update").
			Str("student_id", usrDto.StudentID).
			Msg("Error while connecting to service")

		return nil, &dto.ResponseErr{
			StatusCode: http.StatusServiceUnavailable,
			Message:    "Service is down",
			Data:       nil,
		}
	}

	return res.User, nil
}

func (s *Service) Delete(id string) (result bool, err *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, errRes := s.client.Delete(ctx, &proto.DeleteUserRequest{Id: id})
	if errRes != nil {
		st, ok := status.FromError(errRes)
		if ok {
			switch st.Code() {
			case codes.NotFound:
				return false, &dto.ResponseErr{
					StatusCode: http.StatusNotFound,
					Message:    st.Message(),
					Data:       nil,
				}
			default:

				log.Error().
					Err(errRes).
					Str("service", "user").
					Str("module", "delete").
					Msg("Error while connecting to service")

				return false, &dto.ResponseErr{
					StatusCode: http.StatusServiceUnavailable,
					Message:    "Service is down",
					Data:       nil,
				}
			}
		}

		log.Error().
			Err(errRes).
			Str("service", "user").
			Str("module", "delete").
			Msg("Error while connecting to service")

		return false, &dto.ResponseErr{
			StatusCode: http.StatusServiceUnavailable,
			Message:    "Service is down",
			Data:       nil,
		}
	}

	return res.Success, nil
}

func (s *Service) GetUserEstamp(userid string) (*proto.GetUserEstampResponse, *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := s.client.GetUserEstamp(ctx, &proto.GetUserEstampRequest{
		UId: userid,
	})

	if err != nil {
		st, ok := status.FromError(err)

		if !ok {
			log.Error().
				Err(err).
				Str("service", "estamp").
				Str("module", "find_by_id").
				Msg("\"Error parsing\" error")
			return nil, &dto.ResponseErr{
				StatusCode: http.StatusInternalServerError,
				Message:    "Internal Server Error",
				Data:       nil,
			}
		}

		switch st.Code() {
		case codes.Unavailable:
			log.Error().
				Err(err).
				Str("service", "checkin").
				Str("module", "find_user_estamp").
				Msg("Service is down")
			return nil, &dto.ResponseErr{
				StatusCode: http.StatusServiceUnavailable,
				Message:    "Service is down",
				Data:       nil,
			}
		case codes.NotFound:
			return nil, &dto.ResponseErr{
				StatusCode: http.StatusNotFound,
				Message:    "User not found",
				Data:       nil,
			}
		case codes.PermissionDenied:
			return nil, &dto.ResponseErr{
				StatusCode: http.StatusForbidden,
				Message:    "Forbidden resource",
				Data:       nil,
			}
		default:
			log.Error().
				Err(err).
				Str("service", "checkin").
				Str("module", "find_user_estamp").
				Msg("Unhandled error")
			return nil, &dto.ResponseErr{
				StatusCode: http.StatusInternalServerError,
				Message:    "Internal Server Error",
				Data:       nil,
			}
		}
	}

	return res, nil
}

func (s *Service) ConfirmEstamp(uid string, eid string) (*proto.ConfirmEstampResponse, *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := s.client.ConfirmEstamp(ctx, &proto.ConfirmEstampRequest{
		UId: uid,
		EId: eid,
	})

	if err != nil {
		st, ok := status.FromError(err)

		if !ok {
			log.Error().
				Err(err).
				Str("service", "estamp").
				Str("module", "confirm_estamp").
				Msg("\"Error parsing\" error")
			return nil, &dto.ResponseErr{
				StatusCode: http.StatusInternalServerError,
				Message:    "Internal Server Error",
				Data:       nil,
			}
		}

		switch st.Code() {
		case codes.Unavailable:
			log.Error().
				Err(err).
				Str("service", "checkin").
				Str("module", "confirm_estamp").
				Msg("Service is down")
			return nil, &dto.ResponseErr{
				StatusCode: http.StatusServiceUnavailable,
				Message:    "Service is down",
				Data:       nil,
			}
		case codes.NotFound:
			return nil, &dto.ResponseErr{
				StatusCode: http.StatusNotFound,
				Message:    "User not found",
				Data:       nil,
			}
		case codes.PermissionDenied:
			return nil, &dto.ResponseErr{
				StatusCode: http.StatusForbidden,
				Message:    "Forbidden resource",
				Data:       nil,
			}
		default:
			log.Error().
				Err(err).
				Str("service", "checkin").
				Str("module", "confirm_estamp").
				Msg("Unhandled error")
			return nil, &dto.ResponseErr{
				StatusCode: http.StatusInternalServerError,
				Message:    "Internal Server Error",
				Data:       nil,
			}
		}
	}

	return res, nil
}
