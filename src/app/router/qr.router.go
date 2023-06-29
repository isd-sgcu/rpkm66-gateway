package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/isd-sgcu/rpkm66-gateway/src/pkg/rctx"
)

func (r *FiberRouter) PostQr(path string, h func(rctx.Context)) {
	r.qr.Post(path, func(c *fiber.Ctx) error {
		h(NewFiberCtx(c))
		return nil
	})
}
