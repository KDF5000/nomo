package repository

import (
	"context"

	"github.com/KDF5000/nomo/domain/entity"
)

type LarkBotRegistarRepository interface {
	UpdateOrInsert(ctx context.Context, b *entity.LarkBotRegistar) error
	GetLarkBotRegistarByUnionUserID(ctx context.Context, appID string) (*entity.LarkBotRegistar, error)
}
