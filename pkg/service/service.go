package service

import (
	"context"
	"git.yandex-academy.ru/ooornament/code_architecture/pkg/types"
)

type Service interface {
	CheckCart(ctx context.Context, request *types.CartCheckRequest) types.CartCheckResponse
}
