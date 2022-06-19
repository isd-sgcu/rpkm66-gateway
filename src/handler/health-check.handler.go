package handler

import "net/http"

type HealthCheckHandler struct {
}

type HealthCheckContext interface {
	JSON(statusCode int, v interface{})
}

func NewHealthCheckHandler() *HealthCheckHandler {
	return &HealthCheckHandler{}
}

// HealthCheck is a function that use to check is service health is ok
// @Summary health check
// @Description Check is service heath is ok
// @Tags health check
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router / [get]
func (h *HealthCheckHandler) HealthCheck(c HealthCheckContext) {
	c.JSON(http.StatusOK, map[string]interface{}{
		"Health": "OK!",
	})
	return
}
