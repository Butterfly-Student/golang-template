package postgres_outbound_adapter_test

import (
	"context"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"gorm.io/gorm"

	postgres_outbound_adapter "go-template/internal/adapter/outbound/postgres"
	"go-template/internal/model"
	"go-template/tests/helpers"
)

func TestClientAdapter(t *testing.T) {
	// Integration tests usually take longer, skip in short mode
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := context.Background()

	// Start Postgres Container
	pgContainer, err := helpers.SetupPostgresContainer(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer pgContainer.Terminate(ctx)

	// AutoMigrate the schema
	err = pgContainer.DB.AutoMigrate(&model.Client{})
	if err != nil {
		t.Fatal(err)
	}

	adapter := postgres_outbound_adapter.NewClientAdapter(pgContainer.DB)

	Convey("Test Postgres Client Adapter (Integration)", t, func() {
		// Cleanup before each test to ensure clean state
		pgContainer.DB.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&model.Client{})

		now := time.Now().Truncate(time.Microsecond)
		input := model.ClientInput{
			Name:      "Test Client",
			BearerKey: "test-key-integration",
			CreatedAt: now,
			UpdatedAt: now,
		}

		Convey("Upsert", func() {
			Convey("Insert new record", func() {
				err := adapter.Upsert([]model.ClientInput{input})
				So(err, ShouldBeNil)

				var count int64
				pgContainer.DB.Model(&model.Client{}).Count(&count)
				So(count, ShouldEqual, 1)

				var stored model.Client
				pgContainer.DB.First(&stored)
				So(stored.Name, ShouldEqual, input.Name)
				So(stored.BearerKey, ShouldEqual, input.BearerKey)
			})

			Convey("Update existing record (Conflict on BearerKey)", func() {
				// First insert
				adapter.Upsert([]model.ClientInput{input})

				// Update data
				updatedInput := input
				updatedInput.Name = "Updated Name"

				// Same BearerKey -> Should Update
				err := adapter.Upsert([]model.ClientInput{updatedInput})
				So(err, ShouldBeNil)

				var stored model.Client
				pgContainer.DB.First(&stored, "bearer_key = ?", input.BearerKey)
				So(stored.Name, ShouldEqual, "Updated Name")

				var count int64
				pgContainer.DB.Model(&model.Client{}).Count(&count)
				So(count, ShouldEqual, 1)
			})
		})

		Convey("FindByFilter", func() {
			// Seed data
			adapter.Upsert([]model.ClientInput{input})

			// Get actual ID
			var stored model.Client
			pgContainer.DB.First(&stored, "bearer_key = ?", input.BearerKey)

			Convey("Find by ID", func() {
				filter := model.ClientFilter{IDs: []int{stored.ID}}
				results, err := adapter.FindByFilter(filter, false)
				So(err, ShouldBeNil)
				So(len(results), ShouldEqual, 1)
				So(results[0].ID, ShouldEqual, stored.ID)
			})

			Convey("Find by Name", func() {
				filter := model.ClientFilter{Names: []string{input.Name}}
				results, err := adapter.FindByFilter(filter, false)
				So(err, ShouldBeNil)
				So(len(results), ShouldEqual, 1)
				So(results[0].Name, ShouldEqual, input.Name)
			})

			Convey("With Lock", func() {
				filter := model.ClientFilter{BearerKeys: []string{input.BearerKey}}
				results, err := adapter.FindByFilter(filter, true)
				So(err, ShouldBeNil)
				So(len(results), ShouldEqual, 1)
			})

			Convey("Empty Result", func() {
				filter := model.ClientFilter{Names: []string{"Non Existent"}}
				results, err := adapter.FindByFilter(filter, false)
				So(err, ShouldBeNil)
				So(len(results), ShouldEqual, 0)
			})
		})

		Convey("IsExists", func() {
			adapter.Upsert([]model.ClientInput{input})

			Convey("Exists", func() {
				exists, err := adapter.IsExists(input.BearerKey)
				So(err, ShouldBeNil)
				So(exists, ShouldBeTrue)
			})

			Convey("Not Exists", func() {
				exists, err := adapter.IsExists("non-existent-key")
				So(err, ShouldBeNil)
				So(exists, ShouldBeFalse)
			})
		})

		Convey("DeleteByFilter", func() {
			adapter.Upsert([]model.ClientInput{input})

			Convey("Delete by BearerKey", func() {
				filter := model.ClientFilter{BearerKeys: []string{input.BearerKey}}
				err := adapter.DeleteByFilter(filter)
				So(err, ShouldBeNil)

				var count int64
				pgContainer.DB.Model(&model.Client{}).Count(&count)
				So(count, ShouldEqual, 0)
			})
		})
	})
}
