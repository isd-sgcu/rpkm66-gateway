package dto

type DecomposedFile struct {
	Filename string
	Data     []byte
}

type FileResponse struct {
	Filename string `json:"filename" example:"file-example.jpg-6b86b273ff34fce19d6b804eff5a3f5747ada4eaa22f1d49c01e52ddb7875b4b"`
}
