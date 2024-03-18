package users

import (
	"context"
	"git.yandex-academy.ru/ooornament/code_architecture/pkg/errors"
	"git.yandex-academy.ru/ooornament/code_architecture/pkg/repository"
	"git.yandex-academy.ru/ooornament/code_architecture/pkg/types"
)

type Doctor struct {
	repo repository.Cart
	id   int
}

func NewDoctor(repo repository.Cart, ID int) *Doctor {
	return &Doctor{repo: repo, id: ID}
}

func (doctor *Doctor) ValidateCommon(_ context.Context, _ *types.CartItem) {

}

func (doctor *Doctor) ValidateReceipt(_ context.Context, _ *types.CartItem) {

}
func (doctor *Doctor) ValidateSpecial(ctx context.Context, item *types.CartItem) {
	doctorSpeciality := doctor.repo.GetDoctorSpeciality(ctx, doctor.id)
	itemSpeciality := doctor.repo.GetItemSpeciality(ctx, item.Id)
	if doctorSpeciality != 0 && itemSpeciality != 0 && doctorSpeciality != itemSpeciality {
		item.Error = errors.ItemSpecialWrongSpecific{}
	}
}
