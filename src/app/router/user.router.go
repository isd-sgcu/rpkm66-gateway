package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/isd-sgcu/rpkm66-gateway/src/pkg/rctx"
)

func (r *FiberRouter) GetUser(path string, h func(ctx rctx.Context)) {
	r.user.Get(path, func(c *fiber.Ctx) error {
		h(NewFiberCtx(c))
		return nil
	})
}

func (r *FiberRouter) PostUser(path string, h func(ctx rctx.Context)) {
	r.user.Post(path, func(c *fiber.Ctx) error {
		h(NewFiberCtx(c))
		return nil
	})
}

func (r *FiberRouter) PutUser(path string, h func(ctx rctx.Context)) {
	r.user.Put(path, func(c *fiber.Ctx) error {
		h(NewFiberCtx(c))
		return nil
	})
}

func (r *FiberRouter) PatchUser(path string, h func(ctx rctx.Context)) {
	r.user.Patch(path, func(c *fiber.Ctx) error {
		h(NewFiberCtx(c))
		return nil
	})
}

func (r *FiberRouter) DeleteUser(path string, h func(ctx rctx.Context)) {
	r.user.Delete(path, func(c *fiber.Ctx) error {
		h(NewFiberCtx(c))
		return nil
	})
}
