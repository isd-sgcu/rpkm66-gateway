package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/isd-sgcu/rpkm66-gateway/src/pkg/rctx"
)

func (r *FiberRouter) GetGroup(path string, h func(ctx rctx.Context)) {
	r.group.Get(path, func(c *fiber.Ctx) error {
		h(NewFiberCtx(c))
		return nil
	})
}

func (r *FiberRouter) PostGroup(path string, h func(ctx rctx.Context)) {
	r.group.Post(path, func(c *fiber.Ctx) error {
		h(NewFiberCtx(c))
		return nil
	})
}

func (r *FiberRouter) PutGroup(path string, h func(ctx rctx.Context)) {
	r.group.Put(path, func(c *fiber.Ctx) error {
		h(NewFiberCtx(c))
		return nil
	})
}

func (r *FiberRouter) DeleteGroup(path string, h func(ctx rctx.Context)) {
	r.group.Delete(path, func(c *fiber.Ctx) error {
		h(NewFiberCtx(c))
		return nil
	})
}
