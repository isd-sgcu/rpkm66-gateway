package constant

var AuthExcludePath = map[string]struct{}{
	"POST /auth/verify":       {},
	"POST /auth/refreshToken": {},
}

var MapPath2Phase = map[string][]string{
	"GET /user":                 []string{"register", "select", "eventDay", "eStamp"},
	"PUT /user":                 []string{"register", "eventDay"},
	"PUT /file/image":           []string{"register"},
	"GET /group":                []string{"select"},
	"GET /group/:token":         []string{"select"},
	"PUT /group":                []string{"select"},
	"POST /group/:token":        []string{"select"},
	"DELETE /group/members/:id": []string{"select"},
	"POST /qr/checkin/verify":   []string{"eventDay", "eStamp"},
	"POST /qr/checkin/confirm":  []string{"eventDay", "eStamp"},
	"POST /qr/estamp/verify":    []string{"eStamp"},
	"POST /qr/estamp/confirm":   []string{"eStamp"},
	"GET /estamp":               []string{"eStamp"},
	"GET /estamp/:id":           []string{"eStamp"},
}
