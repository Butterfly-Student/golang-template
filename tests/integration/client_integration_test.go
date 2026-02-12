//go:build integration
// +build integration

package integration_test

import (
	"context"
	"testing"
	"time"

	_ "github.com/lib/pq"
	. "github.com/smartystreets/goconvey/convey"
	"gorm.io/gorm"

	postgres_outbound_adapter "go-template/internal/adapter/outbound/postgres"
	"go-template/internal/model"
	"go-template/tests/helpers"
)

func TestClientIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := context.Background()

	// Use shared helper for container setup
	pgContainer, err := helpers.SetupPostgresContainer(ctx)
	if err != nil {
		t.Fatalf("Failed to setup postgres container: %v", err)
	}
	defer pgContainer.Terminate(ctx)

	// Use GORM AutoMigrate instead of manual SQL
	err = pgContainer.DB.AutoMigrate(&model.Client{})
	if err != nil {
		t.Fatalf("Failed to migrate table: %v", err)
	}

	Convey("Test Client Integration with PostgreSQL", t, func() {
		adapter := postgres_outbound_adapter.NewClientAdapter(pgContainer.DB)

		// Cleanup before test
		pgContainer.DB.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&model.Client{})

		Convey("Full CRUD cycle", func() {
			bearerKey := "integration-key-" + time.Now().Format("20060102150405.000")
			input := model.ClientInput{
				Name:      "Integration Test Client",
				BearerKey: bearerKey,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			Convey("Upsert creates a new client", func() {
				err := adapter.Upsert([]model.ClientInput{input})
				So(err, ShouldBeNil)

				Convey("FindByFilter retrieves the client", func() {
					filter := model.ClientFilter{BearerKeys: []string{bearerKey}}
					results, err := adapter.FindByFilter(filter, false)
					So(err, ShouldBeNil)
					So(len(results), ShouldEqual, 1)
					So(results[0].Name, ShouldEqual, "Integration Test Client")
					So(results[0].BearerKey, ShouldEqual, bearerKey)
				})

				Convey("IsExists returns true for existing client", func() {
					exists, err := adapter.IsExists(bearerKey)
					So(err, ShouldBeNil)
					So(exists, ShouldBeTrue)
				})

				Convey("Upsert updates existing client", func() {
					updatedInput := model.ClientInput{
						Name:      "Updated Client Name",
						BearerKey: bearerKey,
						CreatedAt: input.CreatedAt,
						UpdatedAt: time.Now(),
					}
					err := adapter.Upsert([]model.ClientInput{updatedInput})
					So(err, ShouldBeNil)

					filter := model.ClientFilter{BearerKeys: []string{bearerKey}}
					results, err := adapter.FindByFilter(filter, false)
					So(err, ShouldBeNil)
					So(len(results), ShouldEqual, 1)
					So(results[0].Name, ShouldEqual, "Updated Client Name")
				})

				Convey("DeleteByFilter removes the client", func() {
					filter := model.ClientFilter{BearerKeys: []string{bearerKey}}
					err := adapter.DeleteByFilter(filter)
					So(err, ShouldBeNil)

					exists, err := adapter.IsExists(bearerKey)
					So(err, ShouldBeNil)
					So(exists, ShouldBeFalse)
				})
			})
		})

		Convey("IsExists returns false for non-existent client", func() {
			exists, err := adapter.IsExists("nonexistent-key-" + time.Now().Format("20060102150405"))
			So(err, ShouldBeNil)
			So(exists, ShouldBeFalse)
		})

		Convey("FindByFilter with multiple filters", func() {
			// Cleanup
			pgContainer.DB.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&model.Client{})

			now := time.Now()
			clients := []model.ClientInput{
				{Name: "Client A", BearerKey: "key-a-" + now.Format("150405.000"), CreatedAt: now, UpdatedAt: now},
				{Name: "Client B", BearerKey: "key-b-" + now.Format("150405.000"), CreatedAt: now, UpdatedAt: now},
			}

			err := adapter.Upsert(clients)
			So(err, ShouldBeNil)

			filter := model.ClientFilter{Names: []string{"Client A", "Client B"}}
			results, err := adapter.FindByFilter(filter, false)
			So(err, ShouldBeNil)
			So(len(results), ShouldBeGreaterThanOrEqualTo, 2)
		})
	})
}
