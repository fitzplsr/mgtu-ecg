package model

//go:generate ../../.tool/enumer -type FileFormat

type FileFormat uint8

const (
	EDF FileFormat = iota + 1
)
