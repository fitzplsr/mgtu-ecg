package model

const (
	DefaultLimit = 10
	MaxLimit     = 1000
)

//easyjson:json
type Filter struct {
	Search string `json:"search"`
	Offset uint64 `json:"offset"`
	Limit  uint64 `json:"limit"`
}

func (f *Filter) AlignLimit() {
	switch {
	case f.Limit == 0:
		f.Limit = DefaultLimit
	case f.Limit > MaxLimit:
		f.Limit = MaxLimit
	}
}
