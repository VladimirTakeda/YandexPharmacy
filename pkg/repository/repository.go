package repository

import (
	"context"
	"database/sql"
	"git.yandex-academy.ru/ooornament/code_architecture/pkg/db/postgresql"
	"git.yandex-academy.ru/ooornament/code_architecture/pkg/types"
)

type Cart interface {
	FindCommonItem(ctx context.Context, item *types.CartItem)
	FindReceiptItem(ctx context.Context, item *types.CartItem)
	FindSpecialItem(ctx context.Context, item *types.CartItem)
	CheckUserType(ctx context.Context, ID int) bool
	CheckDoctorType(ctx context.Context, ID int) bool
	GetDoctorSpeciality(ctx context.Context, doctorID int) int
	GetItemSpeciality(ctx context.Context, ItemID int) int
	CheckUserReceipt(ctx context.Context, UserId int, ItemID int) bool
}

type Repository struct {
	Cart
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Cart: &postgresql.Database{Client: db},
	}
}
