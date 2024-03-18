package service

import (
	"context"
	"git.yandex-academy.ru/ooornament/code_architecture/pkg/repository"
	"git.yandex-academy.ru/ooornament/code_architecture/pkg/service/users"
	"git.yandex-academy.ru/ooornament/code_architecture/pkg/types"
)

type CartCheckService struct {
	repo repository.Cart
}

func NewCartCheckService(repo repository.Cart) *CartCheckService {
	return &CartCheckService{repo: repo}
}

func (s *CartCheckService) GetUser(ctx context.Context, UserId int) users.UserType {
	flag := s.repo.CheckUserType(ctx, UserId)
	if flag {
		return users.NewUser(s.repo, UserId)
	}
	flag = s.repo.CheckDoctorType(ctx, UserId)
	if flag {
		return users.NewDoctor(s.repo, UserId)
	}
	return users.NewNoUser()
}

func (s *CartCheckService) FindItems(ctx context.Context, request *types.CartCheckRequest) {
	for _, item := range request.Items {
		if item.ItemType == types.Common {
			s.repo.FindCommonItem(ctx, item)
		}
		if item.ItemType == types.Receipt {
			s.repo.FindReceiptItem(ctx, item)
		}
		if item.ItemType == types.Special {
			s.repo.FindSpecialItem(ctx, item)
		}
	}
}

func (s *CartCheckService) CheckCart(ctx context.Context, request *types.CartCheckRequest) types.CartCheckResponse {
	s.FindItems(ctx, request)

	validator := &users.Validator{UserType: s.GetUser(ctx, request.UserId)}

	validator.Validate(ctx, request)

	response := types.CartCheckResponse{Items: make([]types.CartItemStatus, 0)}

	for _, item := range request.Items {
		if item.Error != nil {
			response.Items = append(response.Items, types.CartItemStatus{Id: item.ItemID, Error: item.Error.ToString()})
		}
	}
	return response
}
