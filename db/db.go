package db

import (
	"fmt"
	//"gorm.io/driver/postgres"
	"gorm.io/driver/mysql"
	"log"
	"os"
	"time"

	"github.com/decagonhq/meddle-api/config"
	"github.com/decagonhq/meddle-api/models"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type GormDB struct {
	DB *gorm.DB
}

func GetDB(c *config.Config) *GormDB {
	gormDB := &GormDB{}
	gormDB.Init(c)
	return gormDB
}

func (g *GormDB) Init(c *config.Config) {
	g.DB = getPostgresDB(c)

	if err := migrate(g.DB); err != nil {
		log.Fatalf("unable to run migrations: %v", err)
	}
}

func getPostgresDB(c *config.Config) *gorm.DB {
	log.Printf("Connecting to postgres: %+v", c)
	//postgresDSN := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d TimeZone=Africa/Lagos",
	//	c.PostgresHost, c.PostgresUser, c.PostgresPassword, c.PostgresDB, c.PostgresPort)
	dsn := "root:toluwase@tcp(127.0.0.1:3306)/meddle?charset=utf8mb4&parseTime=True&loc=Local"
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level Info, Silent, Warn, Error
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,       // Disable color
		},
	)
	postgresDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Fatal(err)
	}
	return postgresDB
}

func migrate(db *gorm.DB) error {
	err := db.AutoMigrate(&models.User{}, &models.BlackList{}, &models.Medication{}, &models.FCMNotificationToken{})
	if err != nil {
		return fmt.Errorf("migrations error: %v", err)
	}

	return nil
}
