package handler

import (
	"git.yandex-academy.ru/ooornament/code_architecture/pkg/errors"
	"git.yandex-academy.ru/ooornament/code_architecture/pkg/types"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

const NulID = 0

func parseType(itemType string) (string, errors.Error) {
	if _, ok := types.ItemTypes[itemType]; ok {
		return itemType, nil
	}
	return itemType, errors.WrongCategory{}
}

func ItemsTranslateFromJson(items []string) []*types.CartItem {
	arr := make([]*types.CartItem, 0)
	for _, str := range items {
		str = strings.ToLower(str)

		parts := strings.Split(str, "_")

		category, parseErr := parseType(parts[0])
		if parseErr != nil {
			arr = append(arr, &types.CartItem{ItemID: str, ItemType: category, Id: NulID, Error: parseErr})
			continue
		}

		itemId, err := strconv.Atoi(parts[1])

		if err != nil {
			arr = append(arr, &types.CartItem{ItemID: str, ItemType: category, Id: NulID, Error: errors.IncorrectItemId{}})
			continue
		}

		arr = append(arr, &types.CartItem{ItemID: str, ItemType: category, Id: itemId, Error: nil})
	}
	return arr
}

func (h *Handler) checkCart(c *gin.Context) {
	userId, err := strconv.Atoi(c.Query("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	itemIds := c.QueryArray("item_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	items := ItemsTranslateFromJson(itemIds)

	request := &types.CartCheckRequest{Items: items, UserId: userId}

	c.JSON(http.StatusOK, h.services.CheckCart(c, request).Items)
}
