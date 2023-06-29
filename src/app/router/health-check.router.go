package router

import (
	"github.com/gofiber/fiber/v2"
	health_check "github.com/isd-sgcu/rpkm66-gateway/src/app/handler/health-check"
)

func (r *FiberRouter) GetHealthCheck(path string, h func(ctx health_check.IContext)) {
	r.Get(path, func(c *fiber.Ctx) error {
		h(NewFiberCtx(c))
		return nil
	})
}
