package auth

var ExcludePath = map[string]struct{}{
	"POST /auth/verify":       {},
	"POST /auth/refreshToken": {},
	"GET /baan":               {},
	"GET /baan/:id":           {},
}

var MapPath2Phase = map[string][]string{
	"PUT /user":                 {"register", "eventDay", "eStamp"},
	"PATCH /user":               {"register", "select", "eventDay", "eStamp"},
	"POST /vaccine/verify":      {"register"},
	"PUT /file/upload":          {"register"},
	"GET /group/:token":         {"select"},
	"POST /group/:token":        {"select"},
	"DELETE /group/leave":       {"select"},
	"PUT /group":                {"select"},
	"GET /baan":                 {"select"},
	"GET /baan/:id":             {"select"},
	"DELETE /group/members/:id": {"select"},
	"POST /qr/checkin/verify":   {"eventDay", "eStamp"},
	"POST /qr/checkin/confirm":  {"eventDay", "eStamp"},
	"POST /qr/estamp/verify":    {"eStamp"},
	"POST /qr/estamp/confirm":   {"eStamp"},
	"GET /estamp":               {"eStamp"},
	"GET /estamp/:id":           {"eStamp"},
}
