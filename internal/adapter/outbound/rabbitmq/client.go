package rabbitmq_outbound_adapter

import (
	"context"

	"go-template/internal/model"
	outbound_port "go-template/internal/port/outbound"
	"go-template/utils/rabbitmq"
)

type clientAdapter struct{}

func NewClientAdapter() outbound_port.ClientMessagePort {
	return &clientAdapter{}
}

func (adapter *clientAdapter) PublishUpsert(datas []model.ClientInput) error {
	err := rabbitmq.Publish(context.Background(), model.UpsertClientMessage, rabbitmq.KindFanOut, "", datas)
	if err != nil {
		return err
	}

	return nil
}
