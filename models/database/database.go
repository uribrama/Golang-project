package database

import (
	"fmt"

	"github.com/ansel1/merry"
	"github.com/uribrama/Golang-project/config"
	"github.com/uribrama/Golang-project/logger"
	"github.com/uribrama/Golang-project/models/person"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB interface {
	Start(logger.Logger, config.Config) (*gorm.DB, error)
}

type Database struct {
	*gorm.DB
}

func (d *Database) Start(ctx logger.Logger, cfg config.Config) (*gorm.DB, error) {
	if d.DB != nil {
		return d.DB, nil
	}

	host := cfg.Get().DBHost
	name := cfg.Get().DBName
	pwd := cfg.Get().DBPassword
	port := cfg.Get().DBPort
	user := cfg.Get().DBUser

	if host == "" || name == "" || pwd == "" || port == 0 || user == "" {
		ctx.Fatal(merry.New("Database is not configure properly"))
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", host, user, pwd, name, port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		ctx.Fatal(err)
	}

	err = db.AutoMigrate(&person.Person{})
	d.DB = db

	return db, err
}
