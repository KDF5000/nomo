package repository

import (
	"context"

	"github.com/KDF5000/nomo/domain/entity"
)

type BindInfoRepository interface {
	UpdateOrInsert(ctx context.Context, b *entity.BindInfo) error
	GetBindInfoByUnionUserID(ctx context.Context, id string) (*entity.BindInfo, error)
}
