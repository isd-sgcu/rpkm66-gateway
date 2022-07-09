package constant

var AuthExcludePath = map[string]struct{}{
	"POST /auth/verify":       {},
	"POST /auth/refreshToken": {},
}

var MapPath2Phase = map[string][]string{
	"GET /user":                 {"register", "select", "eventDay", "eStamp"},
	"PUT /user":                 {"register", "eventDay"},
	"PUT /file/image":           {"register"},
	"GET /group":                {"select"},
	"GET /group/:token":         {"select"},
	"PUT /group":                {"select"},
	"POST /group/:token":        {"select"},
	"DELETE /group/members/:id": {"select"},
	"POST /qr/checkin/verify":   {"eventDay", "eStamp"},
	"POST /qr/checkin/confirm":  {"eventDay", "eStamp"},
	"POST /qr/estamp/verify":    {"eStamp"},
	"POST /qr/estamp/confirm":   {"eStamp"},
	"GET /estamp":               {"eStamp"},
	"GET /estamp/:id":           {"eStamp"},
}
