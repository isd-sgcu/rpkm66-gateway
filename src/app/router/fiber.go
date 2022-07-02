package router

import (
	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/google/uuid"
	guard "github.com/isd-sgcu/rnkm65-gateway/src/app/middleware/auth"
)

type FiberRouter struct {
	*fiber.App
	user fiber.Router
	auth fiber.Router
}

type IGuard interface {
	Use(guard.IContext)
}

func NewFiberRouter(authGuard IGuard, isDebug bool) *FiberRouter {
	r := fiber.New(fiber.Config{
		StrictRouting: true,
		AppName:       "RNKM65 API",
	})

	r.Use(cors.New())
	if isDebug {
		r.Use(logger.New())
	}

	r.Get("/docs/*", swagger.HandlerDefault)

	user := NewGroupRouteWithAuthMiddleware(r, "/user", authGuard.Use)
	auth := NewGroupRouteWithAuthMiddleware(r, "/auth", authGuard.Use)

	return &FiberRouter{r, user, auth}
}

func NewGroupRouteWithAuthMiddleware(r *fiber.App, path string, middleware func(ctx guard.IContext)) fiber.Router {
	return r.Group(path, func(c *fiber.Ctx) error {
		middleware(NewFiberCtx(c))
		return nil
	})
}

type FiberCtx struct {
	*fiber.Ctx
}

func (c *FiberCtx) UserID() string {
	return c.Ctx.Locals("UserId").(string)
}

func NewFiberCtx(c *fiber.Ctx) *FiberCtx {
	return &FiberCtx{c}
}

func (c *FiberCtx) Bind(v interface{}) error {
	return c.Ctx.BodyParser(v)
}

func (c *FiberCtx) JSON(statusCode int, v interface{}) {
	c.Ctx.Status(statusCode).JSON(v)
}

func (c *FiberCtx) ID() (id string, err error) {
	id = c.Params("id")

	_, err = uuid.Parse(id)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (c *FiberCtx) Token() string {
	return c.Ctx.Get(fiber.HeaderAuthorization, "")
}

func (c *FiberCtx) Method() string {
	return c.Ctx.Method()
}

func (c *FiberCtx) Path() string {
	return c.Ctx.Path()
}

func (c *FiberCtx) StoreValue(k string, v string) {
	c.Locals(k, v)
}

func (c *FiberCtx) Next() {
	c.Ctx.Next()
}
