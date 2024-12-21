package model

//easyjson:json
type InternalAnalyseRequest struct {
	InputFile string `json:"input_file"`
}

//easyjson:json
type InternalAnalyseResult struct {
	Result  bool   `json:"result"`
	Predict string `json:"predict"`
}
