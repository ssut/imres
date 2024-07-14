package imres

const (
	tiffTagImageWidth  = 256
	tiffTagImageHeight = 257
)

var (
	ispeSignature = string([]byte{0x00, 0x00, 0x00, 0x14, 0x69, 0x73, 0x70, 0x65})
)
