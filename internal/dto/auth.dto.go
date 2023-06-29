package dto

type TokenPayloadAuth struct {
	UserId string `json:"user_id"`
	Role   string `json:"role"`
}

type VerifyTicket struct {
	Ticket string `json:"ticket" validate:"required"`
}

type Validate struct {
	Token string `json:"token" validate:"jwt"`
}

type RedeemNewToken struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}
