package infra

import (
	"context"
	"gorm.io/gorm"
)

var GormFactory func(ctx context.Context) *gorm.DB
