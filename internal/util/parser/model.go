package parser

type Message struct {
	Query string `json:"query,omitempty"`
	PosX  int32  `json:"pos_x,omitempty"`
	PosY  int32  `json:"pos_y,omitempty"`
}
