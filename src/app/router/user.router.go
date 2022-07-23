package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/isd-sgcu/rnkm65-gateway/src/app/handler/user"
)

func (r *FiberRouter) GetUser(path string, h func(ctx user.IContext)) {
	r.user.Get(path, func(c *fiber.Ctx) error {
		h(NewFiberCtx(c))
		return nil
	})
}

func (r *FiberRouter) PostUser(path string, h func(ctx user.IContext)) {
	r.user.Post(path, func(c *fiber.Ctx) error {
		h(NewFiberCtx(c))
		return nil
	})
}

func (r *FiberRouter) PutUser(path string, h func(ctx user.IContext)) {
	r.user.Put(path, func(c *fiber.Ctx) error {
		h(NewFiberCtx(c))
		return nil
	})
}

func (r *FiberRouter) PatchUser(path string, h func(ctx user.IContext)) {
	r.user.Patch(path, func(c *fiber.Ctx) error {
		h(NewFiberCtx(c))
		return nil
	})
}

func (r *FiberRouter) DeleteUser(path string, h func(ctx user.IContext)) {
	r.user.Delete(path, func(c *fiber.Ctx) error {
		h(NewFiberCtx(c))
		return nil
	})
}
