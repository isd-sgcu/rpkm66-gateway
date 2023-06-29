package utils

import (
	"errors"
	"net/http"
	"testing"

	"github.com/isd-sgcu/rpkm66-gateway/internal/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UtilTest struct {
	suite.Suite
}

func TestUtil(t *testing.T) {
	suite.Run(t, new(UtilTest))
}

func (u *UtilTest) SetupTest() {
}

func (u *UtilTest) TestIsExistedTrue() {
	e := map[string]struct{}{
		"A": {},
		"B": {},
	}

	ok := IsExisted(e, "A")

	assert.True(u.T(), ok)
}

func (u *UtilTest) TestIsExistedFalse() {
	e := map[string]struct{}{
		"A": {},
		"B": {},
	}

	ok := IsExisted(e, "C")

	assert.False(u.T(), ok)
}

func (u *UtilTest) TestFormatPathWithID() {
	testFormatPathWithID(u.T(), "POST", "/user/1", []string{"1"}, "POST /user/:id")
	testFormatPathWithID(u.T(), "GET", "/user/join/1", []string{"1"}, "GET /user/join/:id")
	testFormatPathWithID(u.T(), "PUT", "/group/1/join", []string{"1"}, "PUT /group/:id/join")
	testFormatPathWithID(u.T(), "GET", "/estamp/1", []string{"1"}, "GET /estamp/:id")
	testFormatPathWithID(u.T(), "DELETE", "/group/1/kick/2", []string{"1", "2"}, "DELETE /group/:id/kick/:id")
	testFormatPathWithID(u.T(), "DELETE", "/group", []string{"1", "2"}, "DELETE /group")
}

func testFormatPathWithID(t *testing.T, method string, path string, keys []string, want string) {
	actual := FormatPath(method, path, keys)

	assert.Equal(t, want, actual)
}

func (u *UtilTest) TestGetIntFromStr() {
	var nilSlice []string

	testFindIntFromStr(u.T(), "/user/1", "/", []string{"1"})
	testFindIntFromStr(u.T(), "/user/join/1", "/", []string{"1"})
	testFindIntFromStr(u.T(), "/user/2f434b27-0ecf-47a2-9a61-901fa3303aca", "/", nilSlice)
	testFindIntFromStr(u.T(), "/user", "/", nilSlice)
}

func testFindIntFromStr(t *testing.T, s string, sep string, want []string) {
	actual := FindIntFromStr(s, sep)

	assert.Equal(t, want, actual)
}

func (u *UtilTest) TestFindUUIDFromString() {
	var nilSlice []string

	testFindUUIDFromStr(u.T(), "/user/2f434b27-0ecf-47a2-9a61-901fa3303aca", "/", []string{"2f434b27-0ecf-47a2-9a61-901fa3303aca"})
	testFindUUIDFromStr(u.T(), "/group/join/2f434b27-0ecf-47a2-9a61-901fa3303aca", "/", []string{"2f434b27-0ecf-47a2-9a61-901fa3303aca"})
	testFindUUIDFromStr(u.T(), "/group/97c1504c-50d1-413c-8f47-921f4b2919c1/kick/2f434b27-0ecf-47a2-9a61-901fa3303aca", "/", []string{"97c1504c-50d1-413c-8f47-921f4b2919c1", "2f434b27-0ecf-47a2-9a61-901fa3303aca"})
	testFindUUIDFromStr(u.T(), "/user", "/", nilSlice)
}

func testFindUUIDFromStr(t *testing.T, s string, sep string, want []string) {
	actual := FindUUIDFromStr(s, sep)

	assert.Equal(t, want, actual)
}

func (u *UtilTest) TestFindUUIDFromStringNotFound() {
	var want []string

	actual := FindUUIDFromStr("/user", "/")

	assert.Equal(u.T(), want, actual)
}

func (u *UtilTest) TestMergeStringSlice() {
	var nilSlice []string

	testMergeStringSlice(u.T(), []string{"1"}, []string{"2"}, []string{"1", "2"})
	testMergeStringSlice(u.T(), []string{"1"}, []string{"2", "3"}, []string{"1", "2", "3"})
	testMergeStringSlice(u.T(), []string{}, []string{"2"}, []string{"2"})
	testMergeStringSlice(u.T(), []string{"1"}, []string{}, []string{"1"})
	testMergeStringSlice(u.T(), nilSlice, nilSlice, nilSlice)
}

func testMergeStringSlice(t *testing.T, s1 []string, s2 []string, want []string) {
	actual := MergeStringSlice(s1, s2)

	assert.Equal(t, want, actual)
}

func (u *UtilTest) TestGetIDFromPath() {
	var nilSlice []string

	testGetIDFromPath(u.T(), "/user/2f434b27-0ecf-47a2-9a61-901fa3303aca", []string{"2f434b27-0ecf-47a2-9a61-901fa3303aca"})
	testGetIDFromPath(u.T(), "/user/1", []string{"1"})
	testGetIDFromPath(u.T(), "/user/1/2f434b27-0ecf-47a2-9a61-901fa3303aca", []string{"1", "2f434b27-0ecf-47a2-9a61-901fa3303aca"})
	testGetIDFromPath(u.T(), "/user/2f434b27-0ecf-47a2-9a61-901fa3303aca/2", []string{"2", "2f434b27-0ecf-47a2-9a61-901fa3303aca"})
	testGetIDFromPath(u.T(), "/user", nilSlice)
}

func testGetIDFromPath(t *testing.T, path string, want []string) {
	actual := FindIDFromPath(path)

	assert.Equal(t, want, actual)
}

func (u *UtilTest) TestServiceErrorHandlerBadRequest() {
	want := &dto.ResponseErr{
		StatusCode: http.StatusBadRequest,
		Message:    "Invalid Input",
	}

	err := status.Error(codes.InvalidArgument, "")

	newErr := ServiceErrorHandler(err)

	assert.Equal(u.T(), want, newErr)
}

func (u *UtilTest) TestServiceErrorHandlerUnauthorized() {
	want := &dto.ResponseErr{
		StatusCode: http.StatusUnauthorized,
		Message:    "Unauthorized",
	}

	err := status.Error(codes.Unauthenticated, "")

	newErr := ServiceErrorHandler(err)

	assert.Equal(u.T(), want, newErr)
}

func (u *UtilTest) TestServiceErrorHandlerForbidden() {
	want := &dto.ResponseErr{
		StatusCode: http.StatusForbidden,
		Message:    "Forbidden",
	}

	err := status.Error(codes.PermissionDenied, "")

	newErr := ServiceErrorHandler(err)

	assert.Equal(u.T(), want, newErr)
}

func (u *UtilTest) TestServiceErrorHandlerNotFound() {
	want := &dto.ResponseErr{
		StatusCode: http.StatusNotFound,
		Message:    "Not Found",
	}

	err := status.Error(codes.NotFound, "")

	newErr := ServiceErrorHandler(err)

	assert.Equal(u.T(), want, newErr)
}

func (u *UtilTest) TestServiceErrorHandlerDuplicated() {
	want := &dto.ResponseErr{
		StatusCode: http.StatusConflict,
		Message:    "Duplicated Entity",
	}

	err := status.Error(codes.AlreadyExists, "")

	newErr := ServiceErrorHandler(err)

	assert.Equal(u.T(), want, newErr)
}

func (u *UtilTest) TestServiceErrorHandlerInternal() {
	want := &dto.ResponseErr{
		StatusCode: http.StatusInternalServerError,
		Message:    "Internal Error",
	}

	err := status.Error(codes.Internal, "")

	newErr := ServiceErrorHandler(err)

	assert.Equal(u.T(), want, newErr)
}

func (u *UtilTest) TestServiceErrorHandlerNonGrpcError() {
	want := &dto.ResponseErr{
		StatusCode: http.StatusInternalServerError,
		Message:    "Internal Error",
	}

	err := errors.New("Other error")

	newErr := ServiceErrorHandler(err)

	assert.Equal(u.T(), want, newErr)
}
