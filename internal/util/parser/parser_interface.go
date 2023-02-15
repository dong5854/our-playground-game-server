package parser

type Parser interface {
	Unmarshal(data []byte) error
	Function() string
	Data() string
	Dx() float32
	Dy() float32
}
