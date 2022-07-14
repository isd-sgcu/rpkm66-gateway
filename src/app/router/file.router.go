package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/isd-sgcu/rnkm65-gateway/src/app/handler/file"
)

func (r *FiberRouter) PostFile(path string, h func(ctx file.IContext)) {
	r.file.Post(path, func(c *fiber.Ctx) error {
		h(NewFiberCtx(c))
		return nil
	})
}
