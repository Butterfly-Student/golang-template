package temporal_inbound_adapter

import (
	client_temporal_inbound_adapter "go-template/internal/adapter/inbound/temporal/client"
	"go-template/internal/domain"
	inbound_port "go-template/internal/port/inbound"
)

type adapter struct {
	domain domain.Domain
}

func NewAdapter(
	domain domain.Domain,
) inbound_port.WorkflowPort {
	return &adapter{
		domain: domain,
	}
}

func (a *adapter) Client() inbound_port.ClientWorkflowPort {
	return client_temporal_inbound_adapter.NewClientAdapter(a.domain)
}
