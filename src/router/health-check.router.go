package router

import "github.com/gofiber/fiber/v2"
import "github.com/samithiwat/rnkm65-gateway/src/handler"

func (r *FiberRouter) GetHealthCheck(path string, h func(ctx handler.HealthCheckContext)) {
	r.Get(path, func(c *fiber.Ctx) error {
		h(NewFiberCtx(c))
		return nil
	})
}
