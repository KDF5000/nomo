package persistence

import (
	"context"

	"github.com/KDF5000/nomo/domain/entity"
	"github.com/KDF5000/nomo/domain/repository"
	"gorm.io/gorm"
)

type bindInfoRepo struct {
	db *gorm.DB
}

func NewBindInfoRepo(db *gorm.DB) *bindInfoRepo {
	return &bindInfoRepo{db: db}
}

var _ repository.BindInfoRepository = &bindInfoRepo{}

func (repo *bindInfoRepo) UpdateOrInsert(ctx context.Context, b *entity.BindInfo) error {
	var bind entity.BindInfo
	err := repo.db.Where("union_user_id = ?", b.UnionUserID).First(&bind).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	b.ID = bind.ID
	b.CreatedAt = bind.CreatedAt
	if err := repo.db.Save(b).Error; err != nil {
		return err
	}

	return nil
}

func (repo *bindInfoRepo) GetBindInfoByUnionUserID(ctx context.Context, id string) (*entity.BindInfo, error) {
	var bind entity.BindInfo
	err := repo.db.Where("union_user_id = ?", id).First(&bind).Error
	if err != nil {
		return nil, err
	}

	return &bind, nil
}
