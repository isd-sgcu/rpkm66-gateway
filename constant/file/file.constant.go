package file

var AllowContentType = map[string]struct{}{
	"image/jpeg": {},
	"image/jpg":  {},
	"image/png":  {},
	"image/gif":  {},
}

type Type int

const (
	UnknownType Type = 0
	File             = 1
	Image            = 2
)

type Tag int

const (
	UnknownTag Tag = 0
	Profile        = 1
	Baan           = 2
)
