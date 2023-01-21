package parser

type Parser interface {
	Unmarshal(data []byte) error
	Query() string
}
