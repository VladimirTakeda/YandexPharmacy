package errors

// interface

type Error interface {
	ToString() string
}

// base errors

type WrongCategory struct{}

func (e WrongCategory) ToString() string {
	return "WRONG_CATEGORY"
}

type IncorrectItemId struct{}

func (e IncorrectItemId) ToString() string {
	return "INCORRECT_ITEM_ID"
}

type ItemNotFound struct{}

func (e ItemNotFound) ToString() string {
	return "ITEM_NOT_FOUND"
}

type NoUser struct{}

func (e NoUser) ToString() string {
	return "NO_USER"
}

type NoUserNoReceipt struct{}

func (e NoUserNoReceipt) ToString() string {
	return "NO_USER_NO_RECEIPT"
}

type NoUserSpecialItem struct{}

func (e NoUserSpecialItem) ToString() string {
	return "NO_USER_SPECIAL_ITEM"
}

type NoReceipt struct{}

func (e NoReceipt) ToString() string {
	return "NO_RECEIPT"
}

type ItemIsSpecial struct{}

func (e ItemIsSpecial) ToString() string {
	return "ITEM_IS_SPECIAL"
}

type ItemSpecialWrongSpecific struct{}

func (e ItemSpecialWrongSpecific) ToString() string {
	return "ITEM_SPECIAL_WRONG_SPECIFIC"
}
