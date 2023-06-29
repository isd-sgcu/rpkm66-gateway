package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/isd-sgcu/rpkm66-gateway/src/app/handler/vaccine"
)

func (r *FiberRouter) PostVaccine(path string, h func(ctx vaccine.IContext)) {
	r.vaccine.Post(path, func(c *fiber.Ctx) error {
		h(NewFiberCtx(c))
		return nil
	})
}
