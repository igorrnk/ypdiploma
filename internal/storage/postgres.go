package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/igorrnk/ypdiploma.git/internal/configs"
	"github.com/igorrnk/ypdiploma.git/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type PostgresStorage struct {
	context context.Context
	dbPool  *pgxpool.Pool
}

const (
	insertUser        string = `INSERT INTO users(login, hash) VALUES($1, $2) RETURNING user_id;`
	selectUserByLogin string = `SELECT user_id,hash FROM users WHERE login = $1`
	selectUserID      string = `SELECT user_id FROM users WHERE login = $1`
	selectIsUser      string = `SELECT COUNT(*) FROM users WHERE login = $1`

	insertOrder string = `INSERT INTO orders(user_id, number, status, accrual, uploaded_at) 
		VALUES ($1, $2, $3, $4, $5)`
	selectOrder     string = `SELECT user_id FROM orders WHERE number = $1`
	selectUserOrder string = `SELECT user_id, order_id FROM orders WHERE number = $1`
	selectOrders    string = `SELECT number, status, accrual, uploaded_at FROM orders 
    	WHERE user_id = $1 ORDER BY uploaded_at DESC`
	selectAllOrders string = `SELECT number FROM orders 
    	WHERE status IN ('NEW', 'PROCESSING') ORDER BY uploaded_at`
	updateOrder string = `UPDATE orders SET status = $2, accrual = $3 WHERE number = $1;`

	selectSumAccrual   string = `SELECT sum(accrual) FROM orders WHERE user_id = $1`
	selectSumWithdrawn string = `SELECT sum(sum) FROM withdrawals WHERE user_id = $1`
	insertWithdraw     string = `INSERT INTO withdrawals(user_id, order_num, sum, processed_at)
		VALUES ($1, $2, $3, $4)`
	selectWithdraws string = `SELECT order_num, sum, processed_at FROM withdrawals
		WHERE user_id = $1 ORDER BY processed_at DESC`
)

func NewPostgresStorage(ctx context.Context, config *configs.DBConfigType) (*PostgresStorage, error) {
	storage := &PostgresStorage{
		context: ctx,
	}
	var err error
	storage.dbPool, err = pgxpool.New(context.Background(), config.DBUri)
	if err != nil {
		log.Error().Msgf("New database pool: %v", err)
		return nil, err
	}

	_, err = storage.dbPool.Exec(ctx, createDatabase)
	if err != nil {
		log.Error().Msgf("NewDBStorage: failed create tables: %v", err)
		return nil, err
	}
	return storage, nil
}

// AddUser adds a user into storage.
// Return nil, ErrDB, ErrLoginOccupied.
func (storage *PostgresStorage) AddUser(ctx context.Context, user *model.User) error {
	var count int64
	err := storage.dbPool.QueryRow(ctx, selectIsUser, user.Login).Scan(&count)
	if err != nil {
		return model.ErrDB
	}
	if count != 0 {
		return model.ErrLoginOccupied
	}
	var id int64
	err = storage.dbPool.QueryRow(ctx, insertUser, user.Login, user.Hash).Scan(&id)
	if err != nil {
		return model.ErrDB
	}
	user.ID = fmt.Sprint(id)
	return nil
}

func (storage *PostgresStorage) GetUser(ctx context.Context, user *model.User) error {
	var id int32
	err := storage.dbPool.QueryRow(ctx, selectUserByLogin, user.Login).Scan(&id, &user.Hash)
	if errors.Is(err, pgx.ErrNoRows) {
		return model.ErrNoUser
	}
	user.ID = fmt.Sprint(id)
	return nil
}

func (storage *PostgresStorage) AddOrder(ctx context.Context, user *model.User, order *model.Order) error {
	var userID int32
	err := storage.dbPool.QueryRow(ctx, selectOrder, order.Number).Scan(&userID)
	if errors.Is(err, pgx.ErrNoRows) {
		_, err = storage.dbPool.Exec(ctx, insertOrder, user.ID, order.Number, order.Status, order.Accrual, order.UploadedAt)
		if err != nil {
			return model.ErrDB
		}
		return nil
	}
	if err != nil {
		return model.ErrDB
	}
	if fmt.Sprint(userID) == user.ID {
		return model.ErrOrderUpload
	} else {
		return model.ErrOrderOccupied
	}
}

func (storage *PostgresStorage) GetOrders(ctx context.Context, user *model.User) (orders []*model.Order, err error) {
	rows, err := storage.dbPool.Query(ctx, selectOrders, user.ID)
	orders, err = pgx.CollectRows(rows, func(row pgx.CollectableRow) (*model.Order, error) {
		order := &model.Order{}
		e := row.Scan(&order.Number, &order.Status, &order.Accrual, &order.UploadedAt)
		return order, e
	})
	return orders, err
}

func (storage *PostgresStorage) GetAllOrders(ctx context.Context) (orders []*model.Order, err error) {
	rows, err := storage.dbPool.Query(ctx, selectAllOrders)
	orders, err = pgx.CollectRows(rows, func(row pgx.CollectableRow) (*model.Order, error) {
		order := &model.Order{}
		e := row.Scan(&order.Number)
		return order, e
	})
	return orders, err
}

func (storage *PostgresStorage) UpdateOrder(ctx context.Context, order *model.Order) error {
	_, err := storage.dbPool.Exec(ctx, updateOrder, order.Number, order.Status, order.Accrual)
	if err != nil {
		return model.ErrDB
	}
	return nil
}

func (storage *PostgresStorage) GetSumAccrual(ctx context.Context, user *model.User) (int32, error) {
	var sum *int32
	err := storage.dbPool.QueryRow(ctx, selectSumAccrual, user.ID).Scan(&sum)
	if err != nil {
		return 0, model.ErrDB
	}
	if sum == nil {
		return 0, err
	}
	return *sum, nil
}
func (storage *PostgresStorage) GetSumWithdrawn(ctx context.Context, user *model.User) (int32, error) {
	var sum *int32
	err := storage.dbPool.QueryRow(ctx, selectSumWithdrawn, user.ID).Scan(&sum)
	if err != nil {
		return 0, model.ErrDB
	}
	if sum == nil {
		return 0, err
	}
	return *sum, nil
}

func (storage *PostgresStorage) AddWithdraw(ctx context.Context, user *model.User, withdraw *model.Withdraw) error {

	_, err := storage.dbPool.Exec(ctx, insertWithdraw, user.ID, withdraw.Order, withdraw.Sum, withdraw.ProcessedAt)
	if err != nil {
		log.Error().Err(err)
		return model.ErrDB
	}
	return nil
}

func (storage *PostgresStorage) GetWithdraws(ctx context.Context, user *model.User) (withdraws []*model.Withdraw, err error) {
	rows, err := storage.dbPool.Query(ctx, selectWithdraws, user.ID)
	if err != nil {
		return nil, model.ErrDB
	}
	withdraws, err = pgx.CollectRows(rows, func(row pgx.CollectableRow) (*model.Withdraw, error) {
		withdraw := &model.Withdraw{}
		e := row.Scan(&withdraw.Order, &withdraw.Sum, &withdraw.ProcessedAt)
		return withdraw, e
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, model.ErrNoWithdraws
	}
	if err != nil {
		return nil, model.ErrDB
	}
	return withdraws, err
}
