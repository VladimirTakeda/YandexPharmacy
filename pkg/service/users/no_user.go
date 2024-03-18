package users

import (
	"context"
	"git.yandex-academy.ru/ooornament/code_architecture/pkg/errors"
	"git.yandex-academy.ru/ooornament/code_architecture/pkg/types"
)

type NoUser struct {
}

func NewNoUser() *NoUser {
	return &NoUser{}
}

func (user *NoUser) ValidateCommon(_ context.Context, item *types.CartItem) {
	item.Error = errors.NoUser{}
}

func (user *NoUser) ValidateReceipt(_ context.Context, item *types.CartItem) {
	item.Error = errors.NoUserNoReceipt{}
}
func (user *NoUser) ValidateSpecial(_ context.Context, item *types.CartItem) {
	item.Error = errors.NoUserSpecialItem{}
}
