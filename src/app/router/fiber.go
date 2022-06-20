package router

import (
	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/google/uuid"
)

type FiberRouter struct {
	*fiber.App
	user fiber.Router
}

func NewFiberRouter() *FiberRouter {
	r := fiber.New(fiber.Config{
		StrictRouting: true,
		AppName:       "RNKM65 API",
	})

	r.Use(cors.New())
	r.Use(logger.New())

	r.Get("/docs/*", swagger.HandlerDefault)

	user := r.Group("/user")

	return &FiberRouter{r, user}
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
