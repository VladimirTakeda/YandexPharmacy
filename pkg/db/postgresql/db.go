package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"git.yandex-academy.ru/ooornament/code_architecture/pkg/errors"
	"git.yandex-academy.ru/ooornament/code_architecture/pkg/types"
)

const UnknownSpeciality = 0

type Database struct {
	Client *sql.DB
}

func NewPostgresDb(config Config) (*sql.DB, error) {
	connInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.Database,
	)

	client, err := sql.Open("postgres", connInfo)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (db *Database) FindCommonItem(ctx context.Context, item *types.CartItem) {
	if item.Error != nil {
		return
	}

	row := db.Client.QueryRowContext(ctx, `SELECT id FROM common_item where id = $1 `, item.Id)
	var id int
	if err := row.Scan(&id); err == sql.ErrNoRows {
		item.Error = errors.ItemNotFound{}
	}
}

func (db *Database) FindReceiptItem(ctx context.Context, item *types.CartItem) {
	if item.Error != nil {
		return
	}
	row := db.Client.QueryRowContext(ctx, `SELECT id FROM receipt_item where id = $1 `, item.Id)
	var id int
	if err := row.Scan(&id); err == sql.ErrNoRows {
		item.Error = errors.ItemNotFound{}
	}
}

func (db *Database) FindSpecialItem(ctx context.Context, item *types.CartItem) {
	if item.Error != nil {
		return
	}

	row := db.Client.QueryRowContext(ctx, `SELECT id FROM special_item where id = $1 `, item.Id)
	var id int
	if err := row.Scan(&id); err == sql.ErrNoRows {
		item.Error = errors.ItemNotFound{}
	}
}

func (db *Database) CheckUserType(ctx context.Context, ID int) bool {
	row := db.Client.QueryRowContext(ctx, `SELECT id FROM user_account where id = $1 `, ID)
	var id int
	if err := row.Scan(&id); err != sql.ErrNoRows {
		return true
	}
	return false
}

func (db *Database) CheckDoctorType(ctx context.Context, ID int) bool {
	row := db.Client.QueryRowContext(ctx, `SELECT id FROM doctor_account where id = $1 `, ID)
	var id int
	if err := row.Scan(&id); err != sql.ErrNoRows {
		return true
	}
	return false
}

func (db *Database) GetDoctorSpeciality(ctx context.Context, doctorID int) int {
	row := db.Client.QueryRowContext(ctx, `SELECT specialty_id FROM doctor_account where id = $1 `, doctorID)
	var doctorSpecialtyID int
	if err := row.Scan(&doctorSpecialtyID); err == nil {
		return doctorSpecialtyID
	}
	return UnknownSpeciality
}

func (db *Database) GetItemSpeciality(ctx context.Context, ItemID int) int {
	row := db.Client.QueryRowContext(ctx, `SELECT specialty_id FROM special_item where id = $1 `, ItemID)
	var ItemSpecialtyID int
	if err := row.Scan(&ItemSpecialtyID); err == nil {
		return ItemSpecialtyID
	}
	return UnknownSpeciality
}

func (db *Database) CheckUserReceipt(ctx context.Context, UserID int, ItemID int) bool {
	row := db.Client.QueryRowContext(ctx, `SELECT id FROM receipt where user_id = $1 AND item_id = $2`, UserID, ItemID)
	var receiptID int
	if err := row.Scan(&receiptID); err != nil {
		return false
	}
	return true
}
