package model

//go:generate ../../.tool/enumer -type AnalyseStatus

type AnalyseStatus uint8

const (
	Created AnalyseStatus = iota
	Success
	Failed
)
