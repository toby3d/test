package store

import "gitlab.com/toby3d/test/internal/model"

type (
	CartManager interface {
		Add(item *model.Item) error
		Delete(id uint64) error
		GetById(id uint64) *model.Item
		GetList() (int, []*model.Item)
		Update(item *model.Item) error
	}

	ProductReader interface {
		GetById(id uint64) *model.Product
	}
)
