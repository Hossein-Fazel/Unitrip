package cmd

import (
	"database/sql"
	"log"
	"sync"

	"unitrip/internal/config"
	"unitrip/internal/infrastructure/database"
)

type Container struct {
	db *sql.DB

	mutex sync.RWMutex
}

func NewContainer() *Container {
	return &Container{}
}

func (c *Container) SetUp(db *sql.DB) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.db = db
}

func Run() error {

	conf, err := config.New()
	if err != nil {
		return err
	}

	db := database.NewDB(conf)
	if db == nil {
		log.Println("db is nil ...")
	}
	container := NewContainer()

	container.SetUp(db)

	return nil
}
