package users

import (
	"context"
	"git.yandex-academy.ru/ooornament/code_architecture/pkg/types"
)

type UserType interface {
	ValidateCommon(ctx context.Context, item *types.CartItem)
	ValidateReceipt(ctx context.Context, item *types.CartItem)
	ValidateSpecial(ctx context.Context, item *types.CartItem)
}

type Validator struct {
	UserType
}

func (validator *Validator) Validate(ctx context.Context, request *types.CartCheckRequest) {
	for _, item := range request.Items {
		if item.Error != nil {
			continue
		}

		if item.ItemType == types.Common {
			validator.ValidateCommon(ctx, item)
		}

		if item.ItemType == types.Receipt {
			validator.ValidateReceipt(ctx, item)
		}

		if item.ItemType == types.Special {
			validator.ValidateSpecial(ctx, item)
		}
	}
}
