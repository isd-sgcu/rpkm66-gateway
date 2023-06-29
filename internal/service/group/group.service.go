package group

import (
	"context"
	"net/http"
	"time"

	"github.com/isd-sgcu/rpkm66-gateway/internal/dto"
	"github.com/isd-sgcu/rpkm66-gateway/proto"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
				log.Error().
					Err(errRes).
					Str("service", "group").
					Str("module", "find one").
					Str("user_id", id).
					Msg("Not found")
				return nil, &dto.ResponseErr{
					StatusCode: http.StatusNotFound,
					Message:    st.Message(),
					Data:       nil,
				}
			case codes.InvalidArgument:
				log.Error().
					Err(errRes).
					Str("service", "group").
					Str("module", "find one").
					Str("user_id", id).
					Msg("Invalid user id")
				return nil, &dto.ResponseErr{
					StatusCode: http.StatusBadRequest,
					Message:    st.Message(),
					Data:       nil,
				}
			default:

				log.Error().
					Err(errRes).
					Str("service", "group").
					Str("module", "find one").
					Str("user_id", id).
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
			Str("module", "find one").
			Str("user_id", id).
			Msg("Error while connecting to service")

		return nil, &dto.ResponseErr{
			StatusCode: http.StatusServiceUnavailable,
			Message:    "Service is down",
			Data:       nil,
		}
	}

	log.Info().
		Str("service", "group").
		Str("module", "find one").
		Str("user_id", id).
		Msg("Find group success")
	return res.Group, nil
}

func (s *Service) FindByToken(token string) (result *proto.FindByTokenGroupResponse, err *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, errRes := s.client.FindByToken(ctx, &proto.FindByTokenGroupRequest{Token: token})
	if errRes != nil {
		st, ok := status.FromError(errRes)
		if ok {
			switch st.Code() {
			case codes.NotFound:
				log.Error().
					Err(errRes).
					Str("service", "group").
					Str("module", "find by token").
					Str("token", token).
					Msg("Not found group")
				return nil, &dto.ResponseErr{
					StatusCode: http.StatusNotFound,
					Message:    st.Message(),
					Data:       nil,
				}
			default:

				log.Error().
					Err(errRes).
					Str("service", "group").
					Str("module", "find by token").
					Str("token", token).
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
			Str("module", "find by token").
			Str("token", token).
			Msg("Error while connecting to service")

		return nil, &dto.ResponseErr{
			StatusCode: http.StatusServiceUnavailable,
			Message:    "Service is down",
			Data:       nil,
		}
	}

	log.Info().
		Str("service", "group").
		Str("module", "find by token").
		Str("token", token).
		Msg("Find group success")
	return res, nil
}

func (s *Service) Update(in *dto.GroupDto, leaderId string) (result *proto.Group, err *dto.ResponseErr) {
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

	res, errRes := s.client.Update(ctx, &proto.UpdateGroupRequest{Group: grpDto, LeaderId: leaderId})
	if errRes != nil {
		st, ok := status.FromError(errRes)
		if ok {
			switch st.Code() {
			case codes.NotFound:
				log.Error().
					Err(errRes).
					Str("service", "group").
					Str("module", "update").
					Str("user_id", leaderId).
					Msg("Not found")
				return nil, &dto.ResponseErr{
					StatusCode: http.StatusNotFound,
					Message:    st.Message(),
					Data:       nil,
				}
			case codes.InvalidArgument:
				log.Error().
					Err(errRes).
					Str("service", "group").
					Str("module", "update").
					Str("user_id", leaderId).
					Msg("Invalid user id")
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
					Str("user_id", leaderId).
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
			Str("user_id", leaderId).
			Msg("Error while connecting to service")

		return nil, &dto.ResponseErr{
			StatusCode: http.StatusServiceUnavailable,
			Message:    "Service is down",
			Data:       nil,
		}
	}

	log.Info().
		Str("service", "group").
		Str("module", "update").
		Str("user_id", leaderId).
		Msg("Update group success")
	return res.Group, nil
}

func (s *Service) Join(token string, userId string) (result *proto.Group, err *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, errRes := s.client.Join(ctx, &proto.JoinGroupRequest{Token: token, UserId: userId})
	if errRes != nil {
		st, ok := status.FromError(errRes)
		if ok {
			switch st.Code() {
			case codes.PermissionDenied:
				log.Error().
					Err(errRes).
					Str("service", "group").
					Str("module", "join").
					Str("user_id", userId).
					Msg("Not allowed")
				return nil, &dto.ResponseErr{
					StatusCode: http.StatusForbidden,
					Message:    st.Message(),
					Data:       nil,
				}
			case codes.InvalidArgument:
				log.Error().
					Err(errRes).
					Str("service", "group").
					Str("module", "join").
					Str("user_id", userId).
					Msg("Invalid user id")
				return nil, &dto.ResponseErr{
					StatusCode: http.StatusBadRequest,
					Message:    st.Message(),
					Data:       nil,
				}
			case codes.NotFound:
				log.Error().
					Err(errRes).
					Str("service", "group").
					Str("module", "join").
					Str("user_id", userId).
					Msg("Not found")
				return nil, &dto.ResponseErr{
					StatusCode: http.StatusNotFound,
					Message:    st.Message(),
					Data:       nil,
				}
			default:
				log.Error().
					Err(errRes).
					Str("service", "group").
					Str("module", "join").
					Str("user_id", userId).
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
			Str("module", "join").
			Str("user_id", userId).
			Msg("Error while connecting to service")

		return nil, &dto.ResponseErr{
			StatusCode: http.StatusServiceUnavailable,
			Message:    "Service is down",
			Data:       nil,
		}
	}

	log.Info().
		Str("service", "group").
		Str("module", "join").
		Str("user_id", userId).
		Msg("Join group success")
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
				log.Error().
					Err(errRes).
					Str("service", "group").
					Str("module", "delete member").
					Str("user_id", leaderId).
					Msg("Invalid user id")
				return nil, &dto.ResponseErr{
					StatusCode: http.StatusBadRequest,
					Message:    st.Message(),
					Data:       nil,
				}
			case codes.NotFound:
				log.Error().
					Err(errRes).
					Str("service", "group").
					Str("module", "delete member").
					Str("user_id", leaderId).
					Msg("Not found")
				return nil, &dto.ResponseErr{
					StatusCode: http.StatusNotFound,
					Message:    st.Message(),
					Data:       nil,
				}
			case codes.PermissionDenied:
				log.Error().
					Err(errRes).
					Str("service", "group").
					Str("module", "delete member").
					Str("user_id", leaderId).
					Msg("Not allowed")
				return nil, &dto.ResponseErr{
					StatusCode: http.StatusForbidden,
					Message:    st.Message(),
					Data:       nil,
				}
			default:
				log.Error().
					Err(errRes).
					Str("service", "group").
					Str("module", "delete member").
					Str("user_id", leaderId).
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
			Str("module", "delete member").
			Str("user_id", leaderId).
			Msg("Error while connecting to service")

		return nil, &dto.ResponseErr{
			StatusCode: http.StatusServiceUnavailable,
			Message:    "Service is down",
			Data:       nil,
		}
	}

	log.Info().
		Str("service", "group").
		Str("module", "delete member").
		Str("user_id", leaderId).
		Msg("Delete member success")
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
				log.Error().
					Err(errRes).
					Str("service", "group").
					Str("module", "leave").
					Str("user_id", userId).
					Msg("Not allowed")
				return nil, &dto.ResponseErr{
					StatusCode: http.StatusForbidden,
					Message:    st.Message(),
					Data:       nil,
				}
			case codes.InvalidArgument:
				log.Error().
					Err(errRes).
					Str("service", "group").
					Str("module", "leave").
					Str("user_id", userId).
					Msg("Invalid user id")
				return nil, &dto.ResponseErr{
					StatusCode: http.StatusBadRequest,
					Message:    st.Message(),
					Data:       nil,
				}
			case codes.NotFound:
				log.Error().
					Err(errRes).
					Str("service", "group").
					Str("module", "leave").
					Str("user_id", userId).
					Msg("Not found")
				return nil, &dto.ResponseErr{
					StatusCode: http.StatusNotFound,
					Message:    st.Message(),
					Data:       nil,
				}
			case codes.Internal:
				log.Error().
					Err(errRes).
					Str("service", "group").
					Str("module", "leave").
					Str("user_id", userId).
					Msg("Fail to create group")
				return nil, &dto.ResponseErr{
					StatusCode: http.StatusInternalServerError,
					Message:    st.Message(),
					Data:       nil,
				}
			default:
				log.Error().
					Err(errRes).
					Str("service", "group").
					Str("module", "leave").
					Str("user_id", userId).
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
			Str("user_id", userId).
			Msg("Error while connecting to service")

		return nil, &dto.ResponseErr{
			StatusCode: http.StatusServiceUnavailable,
			Message:    "Service is down",
			Data:       nil,
		}
	}

	log.Info().
		Str("service", "group").
		Str("module", "leave").
		Str("user_id", userId).
		Msg("Leave group success")
	return res.Group, nil
}

func (s *Service) SelectBaan(userId string, baanIds []string) (bool, *dto.ResponseErr) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := s.client.SelectBaan(ctx, &proto.SelectBaanRequest{
		UserId: userId,
		Baans:  baanIds,
	})
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.InvalidArgument:
				log.Error().
					Err(err).
					Str("service", "group").
					Str("module", "leave").
					Str("user_id", userId).
					Msg("invalid baan ids")
				return false, &dto.ResponseErr{
					StatusCode: http.StatusBadRequest,
					Message:    "Invalid numbers of baan or bann is duplicated",
					Data:       nil,
				}

			case codes.PermissionDenied:
				log.Error().
					Err(err).
					Str("service", "group").
					Str("module", "leave").
					Str("user_id", userId).
					Msg("invalid baan ids")
				return false, &dto.ResponseErr{
					StatusCode: http.StatusForbidden,
					Message:    "Forbidden Action",
					Data:       nil,
				}

			case codes.NotFound:
				log.Error().
					Err(err).
					Str("service", "group").
					Str("module", "leave").
					Str("user_id", userId).
					Msg("not found baan")
				return false, &dto.ResponseErr{
					StatusCode: http.StatusNotFound,
					Message:    "Not found baan",
					Data:       nil,
				}
			case codes.Internal:
				log.Error().
					Err(err).
					Str("service", "group").
					Str("module", "leave").
					Str("user_id", userId).
					Msg("Internal error in backend")
				return false, &dto.ResponseErr{
					StatusCode: http.StatusInternalServerError,
					Message:    "Service is down",
					Data:       nil,
				}
			default:
				log.Error().
					Err(err).
					Str("service", "group").
					Str("module", "leave").
					Str("user_id", userId).
					Msg("Error while connecting to service")

				return false, &dto.ResponseErr{
					StatusCode: http.StatusServiceUnavailable,
					Message:    "Service is down",
					Data:       nil,
				}
			}
		}
		log.Error().
			Err(err).
			Str("service", "group").
			Str("module", "delete").
			Str("user_id", userId).
			Msg("Error while connecting to service")

		return false, &dto.ResponseErr{
			StatusCode: http.StatusServiceUnavailable,
			Message:    "Service is down",
			Data:       nil,
		}
	}

	return res.Success, nil
}
