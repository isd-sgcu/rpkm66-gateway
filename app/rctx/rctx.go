package rctx

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/isd-sgcu/rpkm66-gateway/app/dto"
	"github.com/isd-sgcu/rpkm66-gateway/app/utils"
)

type GinCtx struct {
	Ctx *gin.Context
}

func (c *GinCtx) UserID() string {
	v, exist := c.Ctx.Get("UserId")
	if exist {
		return v.(string)
	} else {
		return ""
	}
}

func (c *GinCtx) Role() string {
	v, exist := c.Ctx.Get("Role")
	if exist {
		return v.(string)
	} else {
		return ""
	}
}

func (c *GinCtx) Bind(v interface{}) error {
	return c.Ctx.BindJSON(v)
}

func (c *GinCtx) JSON(statusCode int, v interface{}) {
	c.Ctx.JSON(statusCode, v)
}

func (c *GinCtx) ID() (id string, err error) {
	id = c.Param("id")

	_, err = uuid.Parse(id)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (c *GinCtx) Param(key string) string {
	value := c.Ctx.Param(key)

	return value
}

func (c *GinCtx) Token() string {
	token := c.Ctx.GetHeader("Authorization")

	if token == "" {
		return ""
	}

	if !strings.HasPrefix(token, "Bearer ") {
		return ""
	}
	token = strings.Replace(token, "Bearer ", "", 1)

	return token
}

func (c *GinCtx) StoreValue(k string, v string) {
	c.Ctx.Set(k, v)
}

func (c *GinCtx) Next() {
	c.Ctx.Next()
}

func (c *GinCtx) Query(key string) string {
	return c.Ctx.Query(key)
}

func (c *GinCtx) File(key string, allowContent map[string]struct{}, maxSize int64) (*dto.DecomposedFile, error) {
	file, err := c.Ctx.FormFile(key)
	if err != nil {
		return nil, err
	}

	if !utils.IsExisted(allowContent, file.Header["Content-Type"][0]) {
		return nil, errors.New("Not allow content")
	}

	if file.Size > maxSize {
		return nil, errors.New(fmt.Sprintf("Max file size is %v", maxSize))
	}
	content, err := file.Open()
	if err != nil {
		return nil, errors.New("Cannot read file")
	}

	defer content.Close()

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, content); err != nil {
		return nil, err
	}

	return &dto.DecomposedFile{
		Filename: file.Filename,
		Data:     buf.Bytes(),
	}, nil
}

func (c *GinCtx) GetFormData(key string) string {
	return c.Ctx.PostForm(key)
}
