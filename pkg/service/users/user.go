package users

import (
	"context"
	"git.yandex-academy.ru/ooornament/code_architecture/pkg/errors"
	"git.yandex-academy.ru/ooornament/code_architecture/pkg/repository"
	"git.yandex-academy.ru/ooornament/code_architecture/pkg/types"
)

type User struct {
	repo repository.Cart
	id   int
}

func NewUser(repo repository.Cart, ID int) *User {
	return &User{repo: repo, id: ID}
}

func (user *User) ValidateCommon(_ context.Context, _ *types.CartItem) {

}

func (user *User) ValidateReceipt(ctx context.Context, item *types.CartItem) {
	if !user.repo.CheckUserReceipt(ctx, user.id, item.Id) {
		item.Error = errors.NoReceipt{}
	}
}
func (user *User) ValidateSpecial(ctx context.Context, item *types.CartItem) {
	item.Error = errors.ItemIsSpecial{}
}
