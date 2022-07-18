package dto

type GroupDto struct {
	ID       string     `json:"id" validate:"uuid_optional"`
	LeaderID string     `json:"leader_id" validate:"required"`
	Token    string     `json:"token" validate:"required"`
	Members  []*UserDto `json:"members" validate:"required"`
}

type JoinGroupRequest struct {
	IsLeader bool `json:"is_leader" validate:"required"`
	Members  int  `json:"members" validate:"required"`
}
