package store

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"gitlab.com/toby3d/test/internal/model"
)

func newDataBase(t *testing.T) (*sqlx.DB, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	if !assert.NoError(t, err) {
		assert.FailNow(t, err.Error())
	}
	return sqlx.NewDb(db, "sqlmock"), mock, func() { assert.NoError(t, mock.ExpectationsWereMet()) }
}

func TestAdd(t *testing.T) {
	db, mock, release := newDataBase(t)
	defer release()

	s := NewCartStore(db)
	itemOne := model.Item{ProductId: 42, Quanity: 24}
	itemTwo := model.Item{ProductId: 24, Quanity: 42}

	t.Run("invalid", func(t *testing.T) {
		t.Run("empty", func(t *testing.T) {
			assert.Error(t, s.Add(&model.Item{}))
			count, list := s.GetList()
			assert.Empty(t, list, 0)
			assert.Zero(t, count)
		})
		t.Run("zero quanity", func(t *testing.T) {
			assert.Error(t, s.Add(&model.Item{ProductId: 42}))
			count, list := s.GetList()
			assert.Empty(t, list, 0)
			assert.Zero(t, count)
		})
		t.Run("not exists product", func(t *testing.T) {
			mock.ExpectQuery(`SELECT \* FROM products WHERE id`).
				WillReturnError(sql.ErrNoRows)

			assert.Error(t, s.Add(&model.Item{ProductId: 420, Quanity: 7}))
		})
	})
	t.Run("valid", func(t *testing.T) {
		t.Run("single", func(t *testing.T) {
			mock.ExpectQuery(`SELECT \* FROM products WHERE id`).WillReturnRows(
				mock.NewRows([]string{"id", "name", "price"}).
					AddRow(itemOne.GetProductId(), "Apple", 4.99),
			)
			mock.ExpectQuery(`SELECT \* FROM cart WHERE product_id`)
			mock.ExpectExec(`INSERT INTO cart`).
				WithArgs(itemOne.GetProductId(), itemOne.GetQuanity()).
				WillReturnResult(sqlmock.NewResult(int64(itemOne.GetProductId()), 1))

			assert.NoError(t, s.Add(&itemOne))

			mock.ExpectQuery(`SELECT \* FROM cart`).
				WillReturnRows(sqlmock.NewRows([]string{"product_id", "quanity"}).
					AddRow(itemOne.GetProductId(), itemOne.GetQuanity()),
				)

			count, list := s.GetList()
			assert.Contains(t, list, &itemOne)
			assert.Equal(t, 1, count)
		})
		t.Run("append", func(t *testing.T) {
			mock.ExpectQuery(`SELECT \* FROM products WHERE id`).WillReturnRows(
				mock.NewRows([]string{"id", "name", "price"}).
					AddRow(itemOne.GetProductId(), "Apple", 4.99),
			)
			mock.ExpectQuery(`SELECT \* FROM cart WHERE product_id`).
				WillReturnRows(mock.NewRows([]string{"product_id", "quanity"}).
					AddRow(itemOne.GetProductId(), itemOne.GetQuanity()),
				)
			mock.ExpectExec(`UPDATE cart SET quanity`).
				WithArgs(itemOne.GetProductId(), itemOne.GetQuanity()+6).
				WillReturnResult(sqlmock.NewResult(int64(itemOne.GetProductId()), 1))

			assert.NoError(t, s.Add(&model.Item{
				ProductId: itemOne.GetProductId(),
				Quanity:   6,
			}))

			mock.ExpectQuery(`SELECT \* FROM cart`).
				WillReturnRows(sqlmock.NewRows([]string{"product_id", "quanity"}).
					AddRow(itemOne.GetProductId(), itemOne.GetQuanity()+6),
				)

			count, list := s.GetList()
			assert.Equal(t, 1, count)
			assert.Contains(t, list, &model.Item{
				ProductId: itemOne.ProductId,
				Quanity:   itemOne.GetQuanity() + 6,
			})
		})
		t.Run("add different", func(t *testing.T) {
			mock.ExpectQuery(`SELECT \* FROM products WHERE id`).WillReturnRows(
				mock.NewRows([]string{"id", "name", "price"}).
					AddRow(itemTwo.GetProductId(), "Banana", 2.49),
			)
			mock.ExpectQuery(`SELECT \* FROM cart WHERE product_id`)
			mock.ExpectExec(`INSERT INTO cart`).
				WithArgs(itemTwo.GetProductId(), itemTwo.GetQuanity()).
				WillReturnResult(sqlmock.NewResult(int64(itemTwo.GetProductId()), 1))

			assert.NoError(t, s.Add(&itemTwo))

			mock.ExpectQuery(`SELECT \* FROM cart`).
				WillReturnRows(sqlmock.NewRows([]string{"product_id", "quanity"}).
					AddRow(itemOne.GetProductId(), itemOne.GetQuanity()).
					AddRow(itemTwo.GetProductId(), itemTwo.GetQuanity()),
				)

			count, list := s.GetList()
			assert.Equal(t, 2, count)
			assert.Contains(t, list, &itemOne)
			assert.Contains(t, list, &itemTwo)
		})
	})
}

func TestDelete(t *testing.T) {
	db, mock, release := newDataBase(t)
	defer release()
	t.Run("invalid", func(t *testing.T) {
		mock.ExpectExec(`DELETE FROM cart WHERE product_id`).WithArgs(240)
		assert.Error(t, NewCartStore(db).Delete(240))
	})
	t.Run("valid", func(t *testing.T) {
		mock.ExpectExec(`DELETE FROM cart WHERE product_id`).WithArgs(42).
			WillReturnResult(sqlmock.NewResult(42, 1))
		assert.NoError(t, NewCartStore(db).Delete(42))
	})
}

func TestGetById(t *testing.T) {
	t.Run("cart", func(t *testing.T) {
		t.Run("invalid", func(t *testing.T) {
			db, mock, release := newDataBase(t)
			defer release()
			mock.ExpectQuery(`SELECT \* FROM cart WHERE product_id`).WithArgs(24).
				WillReturnError(sql.ErrNoRows)
			assert.Nil(t, NewCartStore(db).GetById(24))
		})
		t.Run("valid", func(t *testing.T) {
			db, mock, release := newDataBase(t)
			defer release()
			mock.ExpectQuery(`SELECT \* FROM cart WHERE product_id`).WithArgs(42).
				WillReturnRows(mock.NewRows([]string{"product_id", "quanity"}).AddRow(42, 24))

			assert.Equal(t, &model.Item{
				ProductId: 42,
				Quanity:   24,
			}, NewCartStore(db).GetById(42))
		})
	})
	t.Run("product", func(t *testing.T) {
		t.Run("invalid", func(t *testing.T) {
			db, mock, release := newDataBase(t)
			defer release()
			mock.ExpectQuery(`SELECT \* FROM products WHERE id`).WithArgs(24).
				WillReturnError(sql.ErrNoRows)
			assert.Nil(t, NewProductStore(db).GetById(24))
		})
		t.Run("valid", func(t *testing.T) {
			db, mock, release := newDataBase(t)
			defer release()
			mock.ExpectQuery(`SELECT \* FROM products WHERE id`).WithArgs(42).
				WillReturnRows(
					mock.NewRows([]string{"id", "name", "price"}).AddRow(42, "Apple", 4.99),
				)

			assert.Equal(t, &model.Product{
				Id:    42,
				Name:  "Apple",
				Price: 4.99,
			}, NewProductStore(db).GetById(42))
		})
	})
}

func TestGetList(t *testing.T) {
	db, mock, release := newDataBase(t)
	defer release()
	mock.ExpectQuery(`SELECT \* FROM cart ORDER BY product_id ASC`).
		WillReturnRows(mock.NewRows([]string{"product_id", "quanity"}).
			AddRow(24, 7).
			AddRow(420, 11).
			AddRow(42, 24),
		)

	count, list := NewCartStore(db).GetList()
	assert.Equal(t, 3, count)
	assert.Contains(t, list, &model.Item{ProductId: 24, Quanity: 7})
	assert.Contains(t, list, &model.Item{ProductId: 42, Quanity: 24})
	assert.Contains(t, list, &model.Item{ProductId: 420, Quanity: 11})
}

func TestUpdate(t *testing.T) {
	itemOne := model.Item{ProductId: 42, Quanity: 16}
	t.Run("invalid", func(t *testing.T) {
		t.Run("empty", func(t *testing.T) {
			db, _, release := newDataBase(t)
			defer release()
			s := NewCartStore(db)

			assert.Error(t, s.Update(&model.Item{}))
			count, list := s.GetList()
			assert.Empty(t, list)
			assert.Zero(t, count)
		})
		t.Run("zero quanity", func(t *testing.T) {
			db, _, release := newDataBase(t)
			defer release()
			s := NewCartStore(db)

			assert.Error(t, s.Update(&model.Item{ProductId: itemOne.ProductId}))
			count, list := s.GetList()
			assert.Empty(t, list)
			assert.Zero(t, count)
		})
		t.Run("non exists product", func(t *testing.T) {
			db, mock, release := newDataBase(t)
			defer release()
			s := NewCartStore(db)

			mock.ExpectQuery(`SELECT \* FROM products WHERE id`).
				WillReturnError(sql.ErrNoRows)

			assert.Error(t, s.Update(&itemOne))
			count, list := s.GetList()
			assert.Empty(t, list)
			assert.Zero(t, count)
		})
	})
	t.Run("valid", func(t *testing.T) {
		t.Run("create", func(t *testing.T) {
			db, mock, release := newDataBase(t)
			defer release()
			s := NewCartStore(db)

			mock.ExpectQuery(`SELECT \* FROM cart WHERE product_id`)
			mock.ExpectQuery(`SELECT \* FROM products WHERE id`).WillReturnRows(
				mock.NewRows([]string{"id", "name", "price"}).
					AddRow(itemOne.GetProductId(), "Apple", 4.99),
			)
			mock.ExpectQuery(`SELECT \* FROM cart WHERE product_id`)
			mock.ExpectExec(`INSERT INTO cart`).
				WithArgs(itemOne.GetProductId(), itemOne.GetQuanity()).
				WillReturnResult(sqlmock.NewResult(int64(itemOne.GetProductId()), 1))

			assert.NoError(t, s.Update(&itemOne))

			mock.ExpectQuery(`SELECT \* FROM cart`).
				WillReturnRows(sqlmock.NewRows([]string{"product_id", "quanity"}).
					AddRow(itemOne.GetProductId(), itemOne.GetQuanity()),
				)

			count, list := s.GetList()
			assert.Contains(t, list, &itemOne)
			assert.Equal(t, 1, count)
		})
		t.Run("change", func(t *testing.T) {
			db, mock, release := newDataBase(t)
			defer release()
			s := NewCartStore(db)

			mock.ExpectQuery(`SELECT \* FROM cart WHERE product_id`).
				WillReturnRows(mock.NewRows([]string{"product_id", "quanity"}).
					AddRow(itemOne.GetProductId(), itemOne.GetQuanity()),
				)
			mock.ExpectExec(`UPDATE cart SET quanity`).
				WithArgs(itemOne.GetProductId(), 7).
				WillReturnResult(sqlmock.NewResult(int64(itemOne.GetProductId()), 1))

			assert.NoError(t, s.Update(&model.Item{
				ProductId: itemOne.GetProductId(),
				Quanity:   7,
			}))

			mock.ExpectQuery(`SELECT \* FROM cart`).
				WillReturnRows(sqlmock.NewRows([]string{"product_id", "quanity"}).
					AddRow(itemOne.GetProductId(), 7),
				)

			count, list := s.GetList()
			assert.Len(t, list, 1, "length of items store must not be changed")
			assert.Equal(t, 1, count)
			assert.Contains(t, list, &model.Item{
				ProductId: itemOne.GetProductId(),
				Quanity:   7,
			})
		})
		t.Run("remove", func(t *testing.T) {
			t.Run("zero", func(t *testing.T) {
				db, mock, release := newDataBase(t)
				defer release()
				s := NewCartStore(db)

				mock.ExpectQuery(`SELECT \* FROM cart WHERE product_id`).
					WillReturnRows(mock.NewRows([]string{"product_id", "quanity"}).
						AddRow(itemOne.GetProductId(), itemOne.GetQuanity()),
					)
				mock.ExpectExec(`DELETE FROM cart WHERE product_id`).WithArgs(itemOne.GetProductId()).
					WillReturnResult(sqlmock.NewResult(int64(itemOne.GetProductId()), 1))

				assert.NoError(t, s.Update(&model.Item{
					ProductId: itemOne.ProductId,
					Quanity:   itemOne.Quanity * 0,
				}))

				count, list := s.GetList()
				assert.Empty(t, list)
				assert.Zero(t, count)
			})
			t.Run("less than quanity", func(t *testing.T) {
				db, mock, release := newDataBase(t)
				defer release()
				s := NewCartStore(db)

				mock.ExpectQuery(`SELECT \* FROM cart WHERE product_id`).
					WillReturnRows(mock.NewRows([]string{"product_id", "quanity"}).
						AddRow(itemOne.GetProductId(), itemOne.GetQuanity()),
					)
				mock.ExpectExec(`DELETE FROM cart WHERE product_id`).WithArgs(itemOne.GetProductId()).
					WillReturnResult(sqlmock.NewResult(int64(itemOne.GetProductId()), 1))

				assert.NoError(t, s.Update(&model.Item{
					ProductId: itemOne.ProductId,
					Quanity:   itemOne.Quanity * -1,
				}))

				count, list := s.GetList()
				assert.Empty(t, list)
				assert.Zero(t, count)
			})
		})
	})
}
