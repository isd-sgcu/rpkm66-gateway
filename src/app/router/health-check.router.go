package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/isd-sgcu/rpkm66-gateway/src/pkg/rctx"
)

func (r *FiberRouter) GetHealthCheck(path string, h func(ctx rctx.Context)) {
	r.Get(path, func(c *fiber.Ctx) error {
		h(NewFiberCtx(c))
		return nil
	})
}
