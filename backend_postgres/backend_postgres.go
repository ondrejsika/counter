package backend_postgres

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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
	postgresSslmode string,
	hostname string,
) (int, error) {
	// Connection string for the PostgreSQL database
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		postgresHost,
		postgresPort,
		postgresUser,
		postgresPassword,
		postgresDatabase,
		postgresSslmode,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Error().
			Str("hostname", hostname).
			Msg(fmt.Sprintf("error=%s", err))
		return -1, err
	}

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
	postgresSslmode string,
	hostname string,
) (int, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		postgresHost,
		postgresPort,
		postgresUser,
		postgresPassword,
		postgresDatabase,
		postgresSslmode,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Error().
			Str("hostname", hostname).
			Msg(fmt.Sprintf("error=%s", err))
		return -1, err
	}

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

func RunMigrations(
	postgresHost string,
	postgresPort int,
	postgresUser string,
	postgresPassword string,
	postgresDatabase string,
	postgresSslmode string,
	hostname string,
) error {
	// Connection string for the PostgreSQL database
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		postgresHost,
		postgresPort,
		postgresUser,
		postgresPassword,
		postgresDatabase,
		postgresSslmode,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Error().
			Str("hostname", hostname).
			Msg(fmt.Sprintf("error=%s", err))
		return err
	}

	// Migrate the schema
	db.AutoMigrate(&Counter{})

	return nil
}
