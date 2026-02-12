package postgres_outbound_adapter

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"go-template/internal/model"
	outbound_port "go-template/internal/port/outbound"
)

const tableClient = "clients"

type clientAdapter struct {
	db *gorm.DB
}

func NewClientAdapter(
	db *gorm.DB,
) outbound_port.ClientDatabasePort {
	return &clientAdapter{
		db: db,
	}
}

// Upsert inserts or updates client records
func (adapter *clientAdapter) Upsert(datas []model.ClientInput) error {
	// Build the data structures for GORM
	clients := make([]map[string]interface{}, len(datas))
	for i, data := range datas {
		clients[i] = map[string]interface{}{
			"name":       data.Name,
			"bearer_key": data.BearerKey,
			"created_at": data.CreatedAt,
			"updated_at": data.UpdatedAt,
		}
	}

	// Use GORM's Clauses for ON CONFLICT handling
	return adapter.db.Table(tableClient).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "bearer_key"}},
			DoUpdates: clause.AssignmentColumns([]string{"name", "updated_at"}),
		}).
		Create(clients).Error
}

// FindByFilter retrieves clients based on filter criteria
func (adapter *clientAdapter) FindByFilter(filter model.ClientFilter, lock bool) ([]model.Client, error) {
	var clients []model.Client

	query := adapter.db.Table(tableClient)

	// Apply filters
	if len(filter.IDs) > 0 {
		query = query.Where("id IN ?", filter.IDs)
	}

	if len(filter.Names) > 0 {
		query = query.Where("name IN ?", filter.Names)
	}

	if len(filter.BearerKeys) > 0 {
		query = query.Where("bearer_key IN ?", filter.BearerKeys)
	}

	// Add row locking if requested
	if lock {
		query = query.Clauses(clause.Locking{Strength: "UPDATE"})
	}

	// Execute query
	err := query.Find(&clients).Error
	if err != nil {
		return nil, err
	}

	return clients, nil
}

// DeleteByFilter deletes clients based on filter criteria
func (adapter *clientAdapter) DeleteByFilter(filter model.ClientFilter) error {
	query := adapter.db.Table(tableClient)

	// Apply filters
	if len(filter.IDs) > 0 {
		query = query.Where("id IN ?", filter.IDs)
	}

	if len(filter.Names) > 0 {
		query = query.Where("name IN ?", filter.Names)
	}

	if len(filter.BearerKeys) > 0 {
		query = query.Where("bearer_key IN ?", filter.BearerKeys)
	}

	// Execute delete
	return query.Delete(&model.Client{}).Error
}

// IsExists checks if a client exists by bearer key
func (adapter *clientAdapter) IsExists(bearerKey string) (bool, error) {
	var count int64

	err := adapter.db.Table(tableClient).
		Where("bearer_key = ?", bearerKey).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
