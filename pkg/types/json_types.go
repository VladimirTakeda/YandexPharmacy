package types

import (
	"git.yandex-academy.ru/ooornament/code_architecture/pkg/errors"
)

const Special = "special"
const Common = "common"
const Receipt = "receipt"

var ItemTypes = map[string]bool{
	"special": true,
	"common":  true,
	"receipt": true,
}

type CartItem struct {
	ItemID   string
	ItemType string
	Id       int
	Error    errors.Error
}

type CartCheckRequest struct {
	Items  []*CartItem
	UserId int
}

type CartItemStatus struct {
	Id    string `json:"item_id"`
	Error string `json:"problem"`
}

type CartCheckResponse struct {
	Items []CartItemStatus
}
