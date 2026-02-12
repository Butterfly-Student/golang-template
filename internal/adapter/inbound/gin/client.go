package gin_inbound_adapter

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/palantir/stacktrace"

	"go-template/internal/domain"
	"go-template/internal/model"
	inbound_port "go-template/internal/port/inbound"
	"go-template/utils/activity"
)

type clientAdapter struct {
	domain domain.Domain
}

func NewClientAdapter(
	domain domain.Domain,
) inbound_port.ClientHttpPort {
	return &clientAdapter{
		domain: domain,
	}
}

func (h *clientAdapter) Upsert(c *gin.Context) {
	ctx := activity.NewContext("http_client_upsert")
	var payload []model.ClientInput

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, model.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	ctx = activity.WithPayload(ctx, payload)

	results, err := h.domain.Client().Upsert(ctx, payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{
			Success: false,
			Error:   stacktrace.RootCause(err).Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Success: true,
		Data:    results,
	})
}

func (h *clientAdapter) Find(c *gin.Context) {
	ctx := activity.NewContext("http_client_find_by_filter")
	var payload model.ClientFilter

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, model.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	ctx = activity.WithPayload(ctx, payload)

	results, err := h.domain.Client().FindByFilter(ctx, payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{
			Success: false,
			Error:   stacktrace.RootCause(err).Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Success: true,
		Data:    results,
	})
}

func (h *clientAdapter) Delete(c *gin.Context) {
	ctx := activity.NewContext("http_client_delete_by_filter")
	var payload model.ClientFilter

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, model.Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	ctx = activity.WithPayload(ctx, payload)

	err := h.domain.Client().DeleteByFilter(ctx, payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Response{
			Success: false,
			Error:   stacktrace.RootCause(err).Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Success: true,
	})
}
