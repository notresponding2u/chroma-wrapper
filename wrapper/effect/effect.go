package effect

const Static = "CHROMA_STATIC"
const Custom = "CHROMA_CUSTOM2"

type Effect struct {
	Effect string `json:"effect"`
	Param  Param  `json:"param"`
}

type EffectResponse struct {
	Result int64 `json:"result"`
}

type Param struct {
	Color int64 `json:"color"`
}
