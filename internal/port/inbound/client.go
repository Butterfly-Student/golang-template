package inbound_port

import "github.com/gin-gonic/gin"

type ClientHttpPort interface {
	Upsert(c *gin.Context)
	Find(c *gin.Context)
	Delete(c *gin.Context)
}

type ClientMessagePort interface {
	Upsert(a any) bool
}

type ClientCommandPort interface {
	PublishUpsert(name string)
	StartUpsert(name string)
}

type ClientWorkflowPort interface {
	Upsert()
}
