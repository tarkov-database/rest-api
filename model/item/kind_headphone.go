package item

const (
	KindHeadphone Kind = "headphone"
)

type Headphone struct {
	Item `json:",inline" bson:",inline"`

	AmbientVolume  float64    `json:"ambientVol" bson:"ambientVol"`
	DryVolume      float64    `json:"dryVol" bson:"dryVol"`
	Distortion     float64    `json:"distortion" bson:"distortion"`
	HighPassFilter HighPass   `json:"hpf" bson:"hpf"`
	Compressor     Compressor `json:"compressor" bson:"compressor"`
}

type HighPass struct {
	CutoffFrequency float64 `json:"cutoffFreq" bson:"cutoffFreq"`
	Resonance       float64 `json:"resonance" bson:"resonance"`
}

type Compressor struct {
	Attack    float64 `json:"attack" bson:"attack"`
	Gain      float64 `json:"gain" bson:"gain"`
	Release   float64 `json:"release" bson:"release"`
	Treshhold float64 `json:"treshhold" bson:"treshhold"`
	Volume    float64 `json:"volume" bson:"volume"`
}
