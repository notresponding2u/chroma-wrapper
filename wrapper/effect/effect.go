package effect

const Static = "CHROMA_STATIC"
const Custom = "CHROMA_CUSTOM"

type Effect struct {
	Effect string `json:"effect"`
	Param  Param  `json:"param"`
}

type Param struct {
	Color int64 `json:"color"`
}

type Identifier struct {
	Id string `json:"id"`
}
type List struct {
	Ids []string `json:"ids"`
}

type KeyboardGrid struct {
	Effect string       `json:"effect"`
	Param  [6][22]int64 `json:"param"`
}
