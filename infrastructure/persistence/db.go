package persistence

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/KDF5000/nomo/domain/entity"
	"github.com/KDF5000/nomo/domain/repository"
)

type Repositories struct {
	BindInfoRepo        repository.BindInfoRepository
	LarkBotRegistarRepo repository.LarkBotRegistarRepository

	db *gorm.DB
}

func NewRepositories(dbUser, dbPassword, dbPort,
	dbHost, dbName string) (*Repositories, error) {
	dsn := fmt.Sprintf("%s:%s@%s(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPassword, "tcp", dbHost, dbPort, dbName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		return nil, err
	}

	return &Repositories{
		BindInfoRepo:        NewBindInfoRepo(db),
		LarkBotRegistarRepo: NewLarkBotRegistarRepo(db),
		db:                  db,
	}, nil
}

func (s *Repositories) AutoMigrate() error {
	return s.db.AutoMigrate(&entity.BindInfo{}, &entity.LarkBotRegistar{})
}
