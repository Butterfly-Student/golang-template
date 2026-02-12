package inbound_port

import "github.com/gin-gonic/gin"

type PingHttpPort interface {
	GetResource(c *gin.Context)
}
