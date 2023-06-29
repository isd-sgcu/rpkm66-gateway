package utils

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/isd-sgcu/rpkm66-gateway/internal/dto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func IsExisted(e map[string]struct{}, key string) bool {
	_, ok := e[key]
	if ok {
		return true
	}
	return false
}

func FormatPath(method string, path string, keys []string) string {
	for _, key := range keys {
		path = strings.Replace(path, key, ":id", 1)
	}

	return fmt.Sprintf("%v %v", method, path)
}

func FindIntFromStr(s string, sep string) []string {
	spliteds := strings.Split(s, sep)

	var result []string

	for _, splited := range spliteds {
		_, err := strconv.Atoi(splited)
		if err == nil {
			result = append(result, splited)
		}
	}

	return result
}

func FindUUIDFromStr(s string, sep string) []string {
	spliteds := strings.Split(s, sep)

	var result []string

	for _, splited := range spliteds {
		_, err := uuid.Parse(splited)
		if err == nil {
			result = append(result, splited)
		}
	}

	return result
}

func MergeStringSlice(s1 []string, s2 []string) []string {
	return append(s1, s2...)
}

func FindIDFromPath(path string) []string {
	uuids := FindUUIDFromStr(path, "/")
	ids := FindIntFromStr(path, "/")

	return MergeStringSlice(ids, uuids)
}

func BoolAdr(b bool) *bool {
	return &b
}

func ServiceErrorHandler(err error) *dto.ResponseErr {
	st, ok := status.FromError(err)
	if !ok {
		return &dto.ResponseErr{
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal Error",
		}
	}
	switch st.Code() {
	case codes.InvalidArgument:
		return &dto.ResponseErr{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid Input",
		}
	case codes.Unauthenticated:
		return &dto.ResponseErr{
			StatusCode: http.StatusUnauthorized,
			Message:    "Unauthorized",
		}
	case codes.PermissionDenied:
		return &dto.ResponseErr{
			StatusCode: http.StatusForbidden,
			Message:    "Forbidden",
		}
	case codes.NotFound:
		return &dto.ResponseErr{
			StatusCode: http.StatusNotFound,
			Message:    "Not Found",
		}
	case codes.AlreadyExists:
		return &dto.ResponseErr{
			StatusCode: http.StatusConflict,
			Message:    "Duplicated Entity",
		}
	case codes.Internal:
		return &dto.ResponseErr{
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal Error",
		}
	default:
		return &dto.ResponseErr{
			StatusCode: http.StatusServiceUnavailable,
			Message:    "Service is down",
		}
	}
}
