package constant

var AuthExcludePath = map[string]struct{}{
	"POST /auth/verify":       {},
	"POST /auth/refreshToken": {},
}
