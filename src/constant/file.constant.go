package constant

var AllowImageContentType = map[string]struct{}{
	"image/jpeg": {},
	"image/jpg":  {},
	"image/png":  {},
	"image/gif":  {},
}

type FileType string

const (
	File  FileType = "file"
	Image          = "image"
)
