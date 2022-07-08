package shan3

type UserDTO interface {
	GetId() string
	GetName() string
	GetGroupId() string
	GetGroupName() string
	GetInviterId() string
	GetInviterName() string
	GetRole() string
	GetDealerId() string
}
