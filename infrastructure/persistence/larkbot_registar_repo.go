package persistence

import (
	"context"

	"github.com/KDF5000/nomo/domain/entity"
	"github.com/KDF5000/nomo/domain/repository"
	"gorm.io/gorm"
)

type larkBotRegistarRepo struct {
	db *gorm.DB
}

func NewLarkBotRegistarRepo(db *gorm.DB) *larkBotRegistarRepo {
	return &larkBotRegistarRepo{db: db}
}

var _ repository.LarkBotRegistarRepository = &larkBotRegistarRepo{}

func (repo *larkBotRegistarRepo) UpdateOrInsert(ctx context.Context, b *entity.LarkBotRegistar) error {
	var reg entity.LarkBotRegistar
	err := repo.db.Where("app_id = ?", b.AppID).First(&reg).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	b.ID = reg.ID
	b.CreatedAt = reg.CreatedAt
	if err := repo.db.Save(b).Error; err != nil {
		return err
	}

	return nil
}

func (repo *larkBotRegistarRepo) GetLarkBotRegistarByUnionUserID(ctx context.Context, appID string) (*entity.LarkBotRegistar, error) {
	var reg entity.LarkBotRegistar
	err := repo.db.Where("app_id = ?", appID).First(&reg).Error
	if err != nil {
		return nil, err
	}

	return &reg, nil
}
