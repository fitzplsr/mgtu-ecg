package model

import (
	"github.com/fitzplsr/mgtu-ecg/pkg/messages"
)

//easyjson:json
type ErrorResponse struct {
	Error messages.Message `json:"error"`
}
