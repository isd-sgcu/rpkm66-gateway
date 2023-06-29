package user

import (
	"context"
	"time"

	"github.com/isd-sgcu/rpkm66-gateway/internal/dto"
	"github.com/isd-sgcu/rpkm66-gateway/internal/utils"
	"github.com/isd-sgcu/rpkm66-gateway/proto"
	"github.com/rs/zerolog/log"
)

type serviceImpl struct {
	client proto.UserServiceClient
}

func NewService(client proto.UserServiceClient) *serviceImpl {
	return &serviceImpl{
		client: client,
	}
}

func (s *serviceImpl) FindOne(id string) (*proto.User, *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := s.client.FindOne(ctx, &proto.FindOneUserRequest{Id: id})
	if err != nil {
		log.Error().
			Err(err).
			Str("service", "user").
			Str("method", "find one user").
			Msg("Error while find one user")

		return nil, utils.ServiceErrorHandler(err)
	}

	return res.User, nil
}

func (s *serviceImpl) Create(in *dto.UserDto) (*proto.User, *dto.ResponseErr) {
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

	res, err := s.client.Create(ctx, &proto.CreateUserRequest{User: usrDto})
	if err != nil {
		log.Error().
			Err(err).
			Str("service", "user").
			Str("method", "create user").
			Msg("Error while create user")

		return nil, utils.ServiceErrorHandler(err)
	}

	return res.User, nil
}

func (s *serviceImpl) Update(id string, in *dto.UpdateUserDto) (*proto.User, *dto.ResponseErr) {
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

	res, err := s.client.Update(ctx, usrReq)
	if err != nil {
		log.Error().
			Err(err).
			Str("service", "user").
			Str("method", "update user").
			Msg("Error while update user")

		return nil, utils.ServiceErrorHandler(err)
	}

	return res.User, nil
}

func (s *serviceImpl) Verify(studentId string, verifyType string) (bool, *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Info().
		Str("service", "user").
		Str("module", "verify").
		Str("type", verifyType).
		Str("student_id", studentId).
		Msg("Trying to verify the user")

	res, err := s.client.Verify(ctx, &proto.VerifyUserRequest{StudentId: studentId, VerifyType: verifyType})
	if err != nil {
		log.Error().
			Err(err).
			Str("service", "user").
			Str("method", "verify").
			Msg("Error while verify")

		return false, utils.ServiceErrorHandler(err)
	}

	log.Info().
		Str("service", "user").
		Str("module", "verify").
		Str("type", verifyType).
		Str("student_id", studentId).
		Msg("Verified the user")

	return res.Success, nil
}

func (s *serviceImpl) CreateOrUpdate(in *dto.UserDto) (*proto.User, *dto.ResponseErr) {
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

	res, err := s.client.CreateOrUpdate(ctx, &proto.CreateOrUpdateUserRequest{User: usrDto})
	if err != nil {
		log.Error().
			Err(err).
			Str("service", "user").
			Str("method", "upsert user").
			Msg("Error while upsert user")

		return nil, utils.ServiceErrorHandler(err)
	}

	return res.User, nil
}

func (s *serviceImpl) Delete(id string) (bool, *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := s.client.Delete(ctx, &proto.DeleteUserRequest{Id: id})
	if err != nil {
		log.Error().
			Err(err).
			Str("service", "user").
			Str("method", "delete user").
			Msg("Error while delete user")

		return false, utils.ServiceErrorHandler(err)
	}

	return res.Success, nil
}

func (s *serviceImpl) GetUserEstamp(userid string) (*proto.GetUserEstampResponse, *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := s.client.GetUserEstamp(ctx, &proto.GetUserEstampRequest{
		UId: userid,
	})

	if err != nil {
		log.Error().
			Err(err).
			Str("service", "user").
			Str("method", "get user estamp").
			Msg("Error while get user estamp")

		return nil, utils.ServiceErrorHandler(err)
	}

	return res, nil
}

func (s *serviceImpl) ConfirmEstamp(uid string, eid string) (*proto.ConfirmEstampResponse, *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := s.client.ConfirmEstamp(ctx, &proto.ConfirmEstampRequest{
		UId: uid,
		EId: eid,
	})

	if err != nil {
		log.Error().
			Err(err).
			Str("service", "user").
			Str("method", "confirm estamp").
			Msg("Error while confirm estamp")

		return nil, utils.ServiceErrorHandler(err)
	}

	return res, nil
}
