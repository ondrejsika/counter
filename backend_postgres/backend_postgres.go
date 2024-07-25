package backend_postgres

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Counter struct {
	gorm.Model
	ID      int `gorm:"primaryKey"`
	Counter int
}

func DoCountPostgres(
	postgresHost string,
	postgresPort int,
	postgresUser string,
	postgresPassword string,
	postgresDatabase string,
	hostname string,
) (int, error) {
	// Connection string for the PostgreSQL database
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		postgresHost,
		postgresPort,
		postgresUser,
		postgresPassword,
		postgresDatabase,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return -1, err
	}

	// Migrate the schema
	db.AutoMigrate(&Counter{})

	// Ensure the counter exists
	db.Where(Counter{ID: 1}).FirstOrCreate(&Counter{ID: 1, Counter: 0})

	// Increment the counter using Exec to execute SQL directly
	db.Exec("UPDATE counters SET counter = counter + 1 WHERE id = 1")

	// Get the counter value
	var counter Counter
	db.First(&counter, 1)

	// Close DB connection
	sqlDB, _ := db.DB()
	sqlDB.Close()

	return counter.Counter, nil
}

func GetCountPostgres(
	postgresHost string,
	postgresPort int,
	postgresUser string,
	postgresPassword string,
	postgresDatabase string,
	hostname string,
) (int, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		postgresHost,
		postgresPort,
		postgresUser,
		postgresPassword,
		postgresDatabase,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return -1, err
	}

	// Migrate the schema
	db.AutoMigrate(&Counter{})

	// Ensure the counter exists
	db.Where(Counter{ID: 1}).FirstOrCreate(&Counter{ID: 1, Counter: 0})

	// Get the counter value
	var counter Counter
	db.First(&counter, 1)

	// Close DB connection
	sqlDB, _ := db.DB()
	sqlDB.Close()

	return counter.Counter, nil
}
