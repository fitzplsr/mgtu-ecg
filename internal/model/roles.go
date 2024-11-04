package model

//go:generate ../../.tool/enumer -type Role

type Role uint8

const (
	RoleAnonymous Role = iota + 1
	RoleAdmin
	RoleUser
)
