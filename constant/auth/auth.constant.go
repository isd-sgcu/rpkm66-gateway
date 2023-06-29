package auth

import "github.com/isd-sgcu/rpkm66-gateway/constant/phase"

var ExcludePath = map[string]struct{}{
	"POST /auth/verify":       {},
	"POST /auth/refreshToken": {},
	"GET /baan":               {},
	"GET /baan/:id":           {},
	"GET /estamp":             {},
}

var MapPath2Phase = map[string][]string{
	"PUT /user":                 {phase.Register, phase.EventDay, phase.Estamp},
	"PATCH /user":               {phase.Register, phase.Select, phase.EventDay, phase.Estamp},
	"PUT /file/upload":          {phase.Register, phase.EventDay, phase.Estamp},
	"GET /group/:token":         {phase.Select},
	"POST /group/:token":        {phase.Select},
	"DELETE /group/leave":       {phase.Select},
	"PUT /group":                {phase.Select},
	"GET /baan":                 {phase.Select},
	"GET /baan/:id":             {phase.Select},
	"DELETE /group/members/:id": {phase.Select},
	"POST /qr/checkin/verify":   {phase.EventDay, phase.Estamp},
	"POST /qr/checkin/confirm":  {phase.EventDay, phase.Estamp},
	"POST /qr/estamp/verify":    {phase.Estamp},
	"POST /qr/estamp/confirm":   {phase.Estamp},
	"GET /estamp":               {phase.Estamp},
	"GET /estamp/:id":           {phase.Estamp},
}
