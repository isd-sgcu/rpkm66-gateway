package user

import (
	"context"
	"github.com/isd-sgcu/rnkm65-gateway/src/dto"
	"github.com/isd-sgcu/rnkm65-gateway/src/proto"
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
					Message:    st.Message(),
					Data:       nil,
				}
			default:
				return nil, &dto.ResponseErr{
					StatusCode: http.StatusServiceUnavailable,
					Message:    "Service is down",
					Data:       nil,
				}
			}
		}
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
		Firstname:             in.Firstname,
		Lastname:              in.Lastname,
		Nickname:              in.Nickname,
		StudentID:             in.StudentID,
		Faculty:               in.Faculty,
		Year:                  in.Year,
		Phone:                 in.Phone,
		LineID:                in.LineID,
		Email:                 in.Email,
		AllergyFood:           in.AllergyFood,
		FoodRestriction:       in.FoodRestriction,
		AllergyMedicine:       in.AllergyMedicine,
		Disease:               in.Disease,
		VaccineCertificateUrl: in.VaccineCertificateUrl,
		ImageUrl:              in.ImageUrl,
	}

	res, errRes := s.client.Create(ctx, &proto.CreateUserRequest{User: usrDto})
	if errRes != nil {
		return nil, &dto.ResponseErr{
			StatusCode: http.StatusServiceUnavailable,
			Message:    "Service is down",
			Data:       nil,
		}
	}

	return res.User, nil
}

func (s *Service) Update(id string, in *dto.UserDto) (result *proto.User, err *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	cancel()

	usrDto := &proto.User{
		Id:                    id,
		Firstname:             in.Firstname,
		Lastname:              in.Lastname,
		Nickname:              in.Nickname,
		StudentID:             in.StudentID,
		Faculty:               in.Faculty,
		Year:                  in.Year,
		Phone:                 in.Phone,
		LineID:                in.LineID,
		Email:                 in.Email,
		AllergyFood:           in.AllergyFood,
		FoodRestriction:       in.FoodRestriction,
		AllergyMedicine:       in.AllergyMedicine,
		Disease:               in.Disease,
		VaccineCertificateUrl: in.VaccineCertificateUrl,
		ImageUrl:              in.ImageUrl,
	}

	res, errRes := s.client.Update(ctx, &proto.UpdateUserRequest{User: usrDto})
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
				return nil, &dto.ResponseErr{
					StatusCode: http.StatusServiceUnavailable,
					Message:    "Service is down",
					Data:       nil,
				}
			}
		}
		return nil, &dto.ResponseErr{
			StatusCode: http.StatusServiceUnavailable,
			Message:    "Service is down",
			Data:       nil,
		}
	}

	return res.User, nil
}

func (s *Service) CreateOrUpdate(in *dto.UserDto) (result *proto.User, err *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	cancel()

	usrDto := &proto.User{
		Id:                    in.ID,
		Firstname:             in.Firstname,
		Lastname:              in.Lastname,
		Nickname:              in.Nickname,
		StudentID:             in.StudentID,
		Faculty:               in.Faculty,
		Year:                  in.Year,
		Phone:                 in.Phone,
		LineID:                in.LineID,
		Email:                 in.Email,
		AllergyFood:           in.AllergyFood,
		FoodRestriction:       in.FoodRestriction,
		AllergyMedicine:       in.AllergyMedicine,
		Disease:               in.Disease,
		VaccineCertificateUrl: in.VaccineCertificateUrl,
		ImageUrl:              in.ImageUrl,
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
				return nil, &dto.ResponseErr{
					StatusCode: http.StatusServiceUnavailable,
					Message:    "Service is down",
					Data:       nil,
				}
			}
		}
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
	cancel()

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
				return false, &dto.ResponseErr{
					StatusCode: http.StatusServiceUnavailable,
					Message:    "Service is down",
					Data:       nil,
				}
			}
		}
		return false, &dto.ResponseErr{
			StatusCode: http.StatusServiceUnavailable,
			Message:    "Service is down",
			Data:       nil,
		}
	}

	return res.Success, nil
}
