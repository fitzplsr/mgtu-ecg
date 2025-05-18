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

//easyjson:json
type ConvertRequest struct {
	InputFile string `json:"input_file"`
}

//easyjson:json
type ConvertResult struct {
	Channels []Channel `json:"channels"`
	VECG     VECG      `json:"vector_ecg_xyz"`
}

//easyjson:json
type Channel struct {
	Label           string    `json:"label"`
	Dimension       string    `json:"dimension"`
	SampleFrequency float64   `json:"sample_frequency"`
	PhysicalMax     float64   `json:"physical_max"`
	PhysicalMin     float64   `json:"physical_min"`
	DigitalMax      int       `json:"digital_max"`
	DigitalMin      int       `json:"digital_min"`
	Prefilter       string    `json:"prefilter"`
	Transducer      string    `json:"transducer"`
	Signal          []float64 `json:"signal"`
}

//easyjson:json
type VECG struct {
	X []float64 `json:"x"`
	Y []float64 `json:"y"`
	Z []float64 `json:"z"`
}
