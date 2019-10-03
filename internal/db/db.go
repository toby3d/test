package db

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"gitlab.com/toby3d/test/internal/model"
	"golang.org/x/xerrors"
)

// demoProducts представляют собой демо-набор продуктов, добавляемые через AutoMigrate для дальнейшего чтения сторами
var demoProducts = []*model.Product{
	&model.Product{Id: 24, Name: "Banana", Price: 2.49},
	&model.Product{Id: 42, Name: "Apple", Price: 4.99},
	&model.Product{Id: 420, Name: "Bottle of Soda", Price: 10},
}

var ErrDataBaseNotInitialized = xerrors.New("database is not initialized")

// Open открывает соединение с PostgreSQL по указанному адресу с параметрами
func Open(addr string) (*sqlx.DB, error) {
	client, err := sqlx.Connect("postgres", addr)
	if err != nil {
		return nil, err
	}
	if err = client.Ping(); err != nil {
		_ = client.Close()
		return nil, err
	}
	return client, nil
}

// AutoMigrate создаёт таблицы, если они не существуют, для адекватной работы сторов
func AutoMigrate(db *sqlx.DB) (err error) {
	if db == nil || db.DB == nil {
		return ErrDataBaseNotInitialized
	}

	if _, err = db.Exec("CREATE TABLE IF NOT EXISTS products (id SERIAL PRIMARY KEY, name TEXT, price FLOAT)"); err != nil {
		return
	}
	for _, p := range demoProducts {
		if _, err = db.Exec(
			"INSERT INTO products (id, name, price) VALUES ($1, $2, $3)",
			p.GetId(), p.GetName(), p.GetPrice(),
		); err != nil {
			return
		}
	}
	if _, err = db.Exec("CREATE TABLE IF NOT EXISTS cart (product_id SERIAL PRIMARY KEY, quanity INT)"); err != nil {
		return
	}
	return nil
}

// AutoClean удаляет таблицы созданные AutoMigrate
func AutoClean(db *sqlx.DB) (err error) {
	if db == nil || db.DB == nil {
		return ErrDataBaseNotInitialized
	}

	if _, err = db.Exec("DROP TABLE IF EXISTS cart"); err != nil {
		return
	}
	_, err = db.Exec("DROP TABLE IF EXISTS products")
	return
}
