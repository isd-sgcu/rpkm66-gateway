package route

import (
	"github.com/isd-sgcu/rpkm66-gateway/src/constant/auth"
)

const (
	Get    = "GET"
	Post   = "POST"
	Patch  = "PATCH"
	Put    = "PUT"
	Delete = "DELETE"
)

type RouteData struct {
	Path       string
	Method     string
	AllowPerms map[string]struct{}
	Phases     map[string]struct{}
	Debug      bool
}

const Any = "*"

var AllowAll = map[string]struct{}{
	Any: {},
}

var Authenticated = map[string]struct{}{
	auth.ADMIN:       {},
	auth.BAAN_STAFF:  {},
	auth.EVENT_STAFF: {},
	auth.USER:        {},
}

var AllRole = map[string]struct{}{
	auth.ADMIN:           {},
	auth.BAAN_STAFF:      {},
	auth.EVENT_STAFF:     {},
	auth.USER:            {},
	auth.UNAUTHENTICATED: {},
}

var Routes = map[string]RouteData{
	"GET /": {
		Path:       "/",
		Method:     Get,
		AllowPerms: Authenticated,
		Phases:     AllowAll,
		Debug:      false,
	},
	"PUT /user": {
		Method:     Put,
		Path:       "/user",
		AllowPerms: Authenticated,
		Phases:     AllowAll,
		Debug:      false,
	},
	"PATCH /user": {
		Method:     Patch,
		Path:       "/user",
		AllowPerms: Authenticated,
		Phases:     AllowAll,
		Debug:      false,
	},
	"GET /auth/me": {
		Method:     Get,
		Path:       "/auth/me",
		AllowPerms: Authenticated,
		Phases:     AllowAll,
		Debug:      false,
	},
	"POST /auth/verify": {
		Method:     Post,
		Path:       "/auth/verify",
		AllowPerms: AllRole,
		Phases:     AllowAll,
		Debug:      false,
	},
	"POST /auth/refreshToken": {
		Method:     Post,
		Path:       "/auth/refreshToken",
		AllowPerms: AllRole,
		Phases:     AllowAll,
		Debug:      false,
	},
	"PUT /file/upload": {
		Method:     Post,
		Path:       "/file/upload",
		AllowPerms: Authenticated,
		Phases:     AllowAll,
		Debug:      false,
	},
	"GET /baan": {
		Method:     Get,
		Path:       "/baan",
		AllowPerms: AllRole,
		Phases:     AllowAll,
		Debug:      false,
	},
	"GET /baan/:id": {
		Method:     Get,
		Path:       "/baan/:id",
		AllowPerms: AllRole,
		Phases:     AllowAll,
		Debug:      false,
	},
	"GET /group": {
		Method:     Get,
		Path:       "/group",
		AllowPerms: Authenticated,
		Phases:     AllowAll,
		Debug:      false,
	},
	"GET /group/:token": {
		Method:     Get,
		Path:       "/group/:token",
		AllowPerms: Authenticated,
		Phases:     AllowAll,
		Debug:      false,
	},
	"POST /group/:token": {
		Method:     Post,
		Path:       "/group/:token",
		AllowPerms: Authenticated,
		Phases:     AllowAll,
		Debug:      false,
	},
	"DELETE /group/leave": {
		Method:     Delete,
		Path:       "/group/leave",
		AllowPerms: Authenticated,
		Phases:     AllowAll,
		Debug:      false,
	},
	"PUT /group/select": {
		Method:     Put,
		Path:       "/group/select",
		AllowPerms: Authenticated,
		Phases:     AllowAll,
		Debug:      false,
	},
	"DELETE /members/:member_id": {
		Method:     Delete,
		Path:       "/members/:member_id",
		AllowPerms: Authenticated,
		Phases:     AllowAll,
		Debug:      false,
	},
	"POST /qr/checkin/verify": {
		Method:     Post,
		Path:       "/qr/checkin/verify",
		AllowPerms: Authenticated,
		Phases:     AllowAll,
		Debug:      false,
	},
	"POST /qr/checkin/confirm": {
		Method:     Post,
		Path:       "/qr/checkin/confirm",
		AllowPerms: Authenticated,
		Phases:     AllowAll,
		Debug:      false,
	},
	"POST /qr/estamp/verify": {
		Method:     Post,
		Path:       "/qr/estamp/verify",
		AllowPerms: Authenticated,
		Phases:     AllowAll,
		Debug:      false,
	},
	"POST /qr/estamp/confirm": {
		Method:     Post,
		Path:       "/qr/estamp/confirm",
		AllowPerms: Authenticated,
		Phases:     AllowAll,
		Debug:      false,
	},
	"POST /qr/ticket": {
		Method:     Post,
		Path:       "/qr/ticket",
		AllowPerms: Authenticated,
		Phases:     AllowAll,
		Debug:      false,
	},
	"GET /estamp/:id": {
		Method:     Get,
		Path:       "/estamp/:id",
		AllowPerms: Authenticated,
		Phases:     AllowAll,
		Debug:      false,
	},
	"GET /estamp": {
		Method:     Get,
		Path:       "/estamp",
		AllowPerms: AllRole,
		Phases:     AllowAll,
		Debug:      false,
	},
	"GET /estamp/user": {
		Method:     Get,
		Path:       "/estamp/user",
		AllowPerms: Authenticated,
		Phases:     AllowAll,
		Debug:      false,
	},
	"GET /user/:id": {
		Method:     Get,
		Path:       "/user/:id",
		AllowPerms: Authenticated,
		Phases:     AllowAll,
		Debug:      true,
	},
	"POST /user": {
		Method:     Post,
		Path:       "/user",
		AllowPerms: Authenticated,
		Phases:     AllowAll,
		Debug:      true,
	},
	"DELETE /user/:id": {
		Method:     Delete,
		Path:       "/user/:id",
		AllowPerms: Authenticated,
		Phases:     AllowAll,
		Debug:      true,
	},
}
