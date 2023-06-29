package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/isd-sgcu/rpkm66-gateway/src/app/handler/estamp"
)

func (r *FiberRouter) GetEstamp(path string, h func(ctx estamp.IContext)) {
	r.estamp.Get(path, func(c *fiber.Ctx) error {
		h(NewFiberCtx(c))
		return nil
	})
}
