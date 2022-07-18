package group

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
	client proto.GroupServiceClient
}

func NewService(client proto.GroupServiceClient) *Service {
	return &Service{
		client: client,
	}
}

func (s *Service) FindOne(id string) (result *proto.Group, err *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, errRes := s.client.FindOne(ctx, &proto.FindOneGroupRequest{UserId: id})
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
			case codes.InvalidArgument:
				return nil, &dto.ResponseErr{
					StatusCode: http.StatusBadRequest,
					Message:    st.Message(),
					Data:       nil,
				}
			default:

				log.Error().
					Err(errRes).
					Str("service", "group").
					Str("module", "findByToken").
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
			Str("service", "group").
			Str("module", "findByToken").
			Msg("Error while connecting to service")

		return nil, &dto.ResponseErr{
			StatusCode: http.StatusServiceUnavailable,
			Message:    "Service is down",
			Data:       nil,
		}
	}

	return res.Group, nil
}

func (s *Service) FindByToken(token string) (result *proto.Group, err *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, errRes := s.client.FindByToken(ctx, &proto.FindByTokenGroupRequest{Token: token})
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
					Str("service", "group").
					Str("module", "findByToken").
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
			Str("service", "group").
			Str("module", "findByToken").
			Msg("Error while connecting to service")

		return nil, &dto.ResponseErr{
			StatusCode: http.StatusServiceUnavailable,
			Message:    "Service is down",
			Data:       nil,
		}
	}

	return res.Group, nil
}

func (s *Service) Create(id string) (result *proto.Group, err *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, errRes := s.client.Create(ctx, &proto.CreateGroupRequest{UserId: id})
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
			case codes.InvalidArgument:
				return nil, &dto.ResponseErr{
					StatusCode: http.StatusBadRequest,
					Message:    st.Message(),
					Data:       nil,
				}
			case codes.Internal:
				return nil, &dto.ResponseErr{
					StatusCode: http.StatusInternalServerError,
					Message:    st.Message(),
					Data:       nil,
				}
			default:
				log.Error().
					Err(errRes).
					Str("service", "group").
					Str("module", "findByToken").
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
			Str("service", "group").
			Str("module", "create").
			Msg("Error while connecting to service")

		return nil, &dto.ResponseErr{
			StatusCode: http.StatusServiceUnavailable,
			Message:    "Service is down",
			Data:       nil,
		}
	}

	return res.Group, nil
}

func (s *Service) Update(in *dto.GroupDto, leaderId string) (result *proto.Group, err *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var grpMembers []*proto.User

	for _, mem := range in.Members {
		usr := &proto.User{
			Id:              mem.ID,
			Title:           mem.Title,
			Firstname:       mem.Firstname,
			Lastname:        mem.Lastname,
			Nickname:        mem.Nickname,
			Phone:           mem.Phone,
			LineID:          mem.LineID,
			Email:           mem.Email,
			AllergyFood:     mem.AllergyFood,
			FoodRestriction: mem.FoodRestriction,
			AllergyMedicine: mem.AllergyMedicine,
			Disease:         mem.Disease,
			CanSelectBaan:   *mem.CanSelectBaan,
			GroupId:         mem.GroupId,
		}
		grpMembers = append(grpMembers, usr)
	}

	grpDto := &proto.Group{
		Id:       in.ID,
		LeaderID: in.LeaderID,
		Token:    in.Token,
		Members:  grpMembers,
	}

	res, errRes := s.client.Update(ctx, &proto.UpdateGroupRequest{Group: grpDto, LeaderId: leaderId})
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
			case codes.InvalidArgument:
				return nil, &dto.ResponseErr{
					StatusCode: http.StatusBadRequest,
					Message:    st.Message(),
					Data:       nil,
				}
			default:

				log.Error().
					Err(errRes).
					Str("service", "group").
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
			Str("service", "group").
			Str("module", "update").
			Msg("Error while connecting to service")

		return nil, &dto.ResponseErr{
			StatusCode: http.StatusServiceUnavailable,
			Message:    "Service is down",
			Data:       nil,
		}
	}

	return res.Group, nil
}

func (s *Service) Join(token string, userId string, isLeader bool, members int) (result *proto.Group, err *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, errRes := s.client.Join(ctx, &proto.JoinGroupRequest{Token: token, UserId: userId, IsLeader: isLeader, Members: int32(members)})
	if errRes != nil {
		st, ok := status.FromError(errRes)
		if ok {
			switch st.Code() {
			case codes.PermissionDenied:
				return nil, &dto.ResponseErr{
					StatusCode: http.StatusForbidden,
					Message:    st.Message(),
					Data:       nil,
				}
			case codes.InvalidArgument:
				return nil, &dto.ResponseErr{
					StatusCode: http.StatusBadRequest,
					Message:    st.Message(),
					Data:       nil,
				}
			case codes.NotFound:
				return nil, &dto.ResponseErr{
					StatusCode: http.StatusNotFound,
					Message:    st.Message(),
					Data:       nil,
				}
			default:
				log.Error().
					Err(errRes).
					Str("service", "group").
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
			Str("service", "group").
			Str("module", "delete").
			Msg("Error while connecting to service")

		return nil, &dto.ResponseErr{
			StatusCode: http.StatusServiceUnavailable,
			Message:    "Service is down",
			Data:       nil,
		}
	}

	return res.Group, nil
}

func (s *Service) DeleteMember(userId string, leaderId string) (result *proto.Group, err *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, errRes := s.client.DeleteMember(ctx, &proto.DeleteMemberGroupRequest{UserId: userId, LeaderId: leaderId})
	if errRes != nil {
		st, ok := status.FromError(errRes)
		if ok {
			switch st.Code() {
			case codes.InvalidArgument:
				return nil, &dto.ResponseErr{
					StatusCode: http.StatusBadRequest,
					Message:    st.Message(),
					Data:       nil,
				}
			case codes.NotFound:
				return nil, &dto.ResponseErr{
					StatusCode: http.StatusNotFound,
					Message:    st.Message(),
					Data:       nil,
				}
			case codes.PermissionDenied:
				return nil, &dto.ResponseErr{
					StatusCode: http.StatusForbidden,
					Message:    st.Message(),
					Data:       nil,
				}
			default:
				log.Error().
					Err(errRes).
					Str("service", "group").
					Str("module", "delete").
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
			Str("service", "group").
			Str("module", "delete").
			Msg("Error while connecting to service")

		return nil, &dto.ResponseErr{
			StatusCode: http.StatusServiceUnavailable,
			Message:    "Service is down",
			Data:       nil,
		}
	}

	return res.Group, nil
}

func (s *Service) Leave(userId string) (result *proto.Group, err *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, errRes := s.client.Leave(ctx, &proto.LeaveGroupRequest{UserId: userId})
	if errRes != nil {
		st, ok := status.FromError(errRes)
		if ok {
			switch st.Code() {
			case codes.PermissionDenied:
				return nil, &dto.ResponseErr{
					StatusCode: http.StatusForbidden,
					Message:    st.Message(),
					Data:       nil,
				}
			case codes.InvalidArgument:
				return nil, &dto.ResponseErr{
					StatusCode: http.StatusBadRequest,
					Message:    st.Message(),
					Data:       nil,
				}
			case codes.NotFound:
				return nil, &dto.ResponseErr{
					StatusCode: http.StatusNotFound,
					Message:    st.Message(),
					Data:       nil,
				}
			case codes.Internal:
				return nil, &dto.ResponseErr{
					StatusCode: http.StatusInternalServerError,
					Message:    st.Message(),
					Data:       nil,
				}
			default:
				log.Error().
					Err(errRes).
					Str("service", "group").
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
			Str("service", "group").
			Str("module", "delete").
			Msg("Error while connecting to service")

		return nil, &dto.ResponseErr{
			StatusCode: http.StatusServiceUnavailable,
			Message:    "Service is down",
			Data:       nil,
		}
	}

	return res.Group, nil
}
