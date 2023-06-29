package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/isd-sgcu/rpkm66-gateway/src/interfaces/qr"
)

func (r *FiberRouter) PostQr(path string, h func(qr.IContext)) {
	r.qr.Post(path, func(c *fiber.Ctx) error {
		h(NewFiberCtx(c))
		return nil
	})
}
