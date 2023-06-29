package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/isd-sgcu/rpkm66-gateway/src/app/handler/baan"
)

func (r *FiberRouter) GetBaan(path string, h func(ctx baan.IContext)) {
	r.baan.Get(path, func(c *fiber.Ctx) error {
		h(NewFiberCtx(c))
		return nil
	})
}
