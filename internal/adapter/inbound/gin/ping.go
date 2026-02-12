package gin_inbound_adapter

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"go-template/internal/domain"
	inbound_port "go-template/internal/port/inbound"
)

type pingAdapter struct {
	domain domain.Domain
}

func NewPingAdapter(
	domain domain.Domain,
) inbound_port.PingHttpPort {
	return &pingAdapter{
		domain: domain,
	}
}

func (h *pingAdapter) GetResource(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
