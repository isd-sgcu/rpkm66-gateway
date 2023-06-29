package group

import (
	"context"
	"time"

	"github.com/isd-sgcu/rpkm66-gateway/internal/dto"
	"github.com/isd-sgcu/rpkm66-gateway/internal/utils"
	"github.com/isd-sgcu/rpkm66-gateway/proto"
	"github.com/rs/zerolog/log"
)

type serviceImpl struct {
	client proto.GroupServiceClient
}

func NewService(client proto.GroupServiceClient) *serviceImpl {
	return &serviceImpl{
		client: client,
	}
}

func (s *serviceImpl) FindOne(id string) (*proto.Group, *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := s.client.FindOne(ctx, &proto.FindOneGroupRequest{UserId: id})
	if err != nil {
		log.Error().
			Err(err).
			Str("service", "group").
			Str("method", "find one").
			Msg("Error while find one")

		return nil, utils.ServiceErrorHandler(err)
	}

	return res.Group, nil
}

func (s *serviceImpl) FindByToken(token string) (*proto.FindByTokenGroupResponse, *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := s.client.FindByToken(ctx, &proto.FindByTokenGroupRequest{Token: token})
	if err != nil {
		log.Error().
			Err(err).
			Str("service", "group").
			Str("method", "find by token").
			Msg("Error while find by token")

		return nil, utils.ServiceErrorHandler(err)
	}
	return res, nil
}

func (s *serviceImpl) Update(in *dto.GroupDto, leaderId string) (*proto.Group, *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var grpMembers []*proto.UserInfo

	for _, mem := range in.Members {
		usr := &proto.UserInfo{
			Id:        mem.ID,
			Firstname: mem.Firstname,
			Lastname:  mem.Lastname,
			ImageUrl:  mem.ImageUrl,
		}
		grpMembers = append(grpMembers, usr)
	}

	grpDto := &proto.Group{
		Id:       in.ID,
		LeaderID: in.LeaderID,
		Token:    in.Token,
		Members:  grpMembers,
	}

	res, err := s.client.Update(ctx, &proto.UpdateGroupRequest{Group: grpDto, LeaderId: leaderId})
	if err != nil {
		log.Error().
			Err(err).
			Str("service", "group").
			Str("method", "find by token").
			Msg("Error while find by token")

		return nil, utils.ServiceErrorHandler(err)
	}

	return res.Group, nil
}

func (s *serviceImpl) Join(token string, userId string) (*proto.Group, *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := s.client.Join(ctx, &proto.JoinGroupRequest{Token: token, UserId: userId})
	if err != nil {
		log.Error().
			Err(err).
			Str("service", "group").
			Str("method", "join group").
			Msg("Error while join group")

		return nil, utils.ServiceErrorHandler(err)
	}

	return res.Group, nil
}

func (s *serviceImpl) DeleteMember(userId string, leaderId string) (*proto.Group, *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := s.client.DeleteMember(ctx, &proto.DeleteMemberGroupRequest{UserId: userId, LeaderId: leaderId})
	if err != nil {
		log.Error().
			Err(err).
			Str("service", "group").
			Str("method", "delete member").
			Msg("Error while delete member")

		return nil, utils.ServiceErrorHandler(err)
	}

	return res.Group, nil
}

func (s *serviceImpl) Leave(userId string) (*proto.Group, *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := s.client.Leave(ctx, &proto.LeaveGroupRequest{UserId: userId})
	if err != nil {
		log.Error().
			Err(err).
			Str("service", "group").
			Str("method", "leave group").
			Msg("Error while leave group")

		return nil, utils.ServiceErrorHandler(err)
	}

	return res.Group, nil
}

func (s *serviceImpl) SelectBaan(userId string, baanIds []string) (bool, *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := s.client.SelectBaan(ctx, &proto.SelectBaanRequest{
		UserId: userId,
		Baans:  baanIds,
	})
	if err != nil {
		log.Error().
			Err(err).
			Str("service", "group").
			Str("method", "select baan").
			Msg("Error while select baan")

		return false, utils.ServiceErrorHandler(err)
	}

	return res.Success, nil
}
