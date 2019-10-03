package store

import (
	"github.com/jmoiron/sqlx"
	"gitlab.com/toby3d/test/internal/model"
)

type (
	CartStore    struct{ conn *sqlx.DB }
	ProductStore struct{ conn *sqlx.DB }

	itemResult struct {
		ProductID uint64 `db:"product_id"`
		Quanity   int32  `db:"quanity"`
	}

	productResult struct {
		ID    uint64  `db:"id"`
		Name  string  `db:"name"`
		Price float32 `db:"price"`
	}
)

func NewCartStore(conn *sqlx.DB) *CartStore { return &CartStore{conn: conn} }

func NewProductStore(conn *sqlx.DB) *ProductStore { return &ProductStore{conn: conn} }

func (s *CartStore) Add(item *model.Item) (err error) {
	switch {
	case item == nil, item.GetProductId() <= 0:
		return ErrNoProductId
	case item.GetQuanity() <= 0:
		return ErrZeroQuanity
	}

	// NOTE(toby3d): Сначала проверяем существование продукта, который мы хотим добавить в корзину
	var p productResult
	if err = s.conn.Get(&p, `SELECT * FROM products WHERE id = $1`, item.GetProductId()); err != nil {
		return err
	}

	// NOTE(toby3d): возможно продукт уже в корзине и мы просто хотим увеличить его количество
	if i := s.GetById(item.ProductId); i != nil { // NOTE(toby3d): продукт уже в корзине, увеличиваем
		i.Quanity += item.Quanity
		_, err = s.conn.Exec(`UPDATE cart SET quanity = $2 WHERE product_id = $1`, p.ID, i.GetQuanity())
	} else { // NOTE(toby3d): объекта в корзине ещё нет, добавляем
		_, err = s.conn.Exec(
			"INSERT INTO cart (product_id, quanity) VALUES ($1, $2)", p.ID, item.GetQuanity(),
		)
	}
	return err
}

func (s *CartStore) Delete(id uint64) (err error) {
	_, err = s.conn.Exec(`DELETE FROM cart WHERE product_id = $1`, id)
	return
}

func (s *CartStore) GetById(id uint64) *model.Item {
	var item itemResult
	if err := s.conn.Get(&item, "SELECT * FROM cart WHERE product_id=$1", id); err != nil {
		return nil
	}

	return &model.Item{
		ProductId: item.ProductID,
		Quanity:   item.Quanity,
	}
}

func (s *CartStore) GetList() (int, []*model.Item) {
	rows, err := s.conn.Queryx(`SELECT * FROM cart ORDER BY product_id ASC`)
	if err != nil {
		return 0, nil
	}
	defer rows.Close()

	var items []*model.Item
	for rows.Next() {
		var item itemResult
		if err = rows.StructScan(&item); err != nil {
			continue
		}
		items = append(items, &model.Item{
			ProductId: item.ProductID,
			Quanity:   item.Quanity,
		})
	}
	if rows.Err() != nil {
		return 0, nil
	}

	return len(items), items
}

func (s *CartStore) Update(item *model.Item) (err error) {
	if i := s.GetById(item.GetProductId()); i != nil {
		if item.GetQuanity() <= 0 {
			return s.Delete(i.GetProductId())
		}

		_, err = s.conn.Exec(`UPDATE cart SET quanity = $2 WHERE product_id = $1`, i.GetProductId(), item.GetQuanity())
		return err
	}
	return s.Add(item)
}

func (s *ProductStore) GetById(id uint64) *model.Product {
	var product productResult
	if err := s.conn.Get(&product, "SELECT * FROM products WHERE id=$1", id); err != nil {
		return nil
	}

	return &model.Product{
		Id:    product.ID,
		Name:  product.Name,
		Price: product.Price,
	}
}
