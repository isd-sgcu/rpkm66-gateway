package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/isd-sgcu/rnkm65-gateway/src/app/handler/health-check"
)

func (r *FiberRouter) GetHealthCheck(path string, h func(ctx health_check.Context)) {
	r.Get(path, func(c *fiber.Ctx) error {
		h(NewFiberCtx(c))
		return nil
	})
}
