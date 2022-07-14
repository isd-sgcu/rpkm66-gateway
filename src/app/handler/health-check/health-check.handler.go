package health_check

import "net/http"

type Handler struct {
}

type IContext interface {
	JSON(statusCode int, v interface{})
}

func NewHandler() *Handler {
	return &Handler{}
}

// HealthCheck is a function that use to check is service health is ok
// @Summary health check
// @Description Check is service heath is ok
// @Tags health check
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router / [get]
func (h *Handler) HealthCheck(c IContext) {
	c.JSON(http.StatusOK, map[string]interface{}{
		"Health": "OK!",
	})
	return
}
