package model

//go:generate ../../.tool/enumer -type AnalyseResult

type AnalyseResult uint8

const (
	Unspecified AnalyseResult = iota
	True
	False
)

func AnalyseResultFromBool(res bool) AnalyseResult {
	if res {
		return True
	}
	return False
}
