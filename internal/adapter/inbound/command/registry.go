package command_inbound_adapter

import (
	"go-template/internal/domain"
	inbound_port "go-template/internal/port/inbound"
)

type adapter struct {
	domain domain.Domain
}

func NewAdapter(
	domain domain.Domain,
) inbound_port.CommandPort {
	return &adapter{
		domain: domain,
	}
}

func (s *adapter) Client() inbound_port.ClientCommandPort {
	return NewClientAdapter(s.domain)
}
