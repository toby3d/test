package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/toby3d/test/internal/model"
)

func TestInMemoryAdd(t *testing.T) {
	itemOne := model.Item{ProductId: 42, Quanity: 24}
	itemTwo := model.Item{ProductId: 24, Quanity: 42}
	t.Run("invalid", func(t *testing.T) {
		s := NewInMemoryCartStore()
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
	})
	t.Run("valid", func(t *testing.T) {
		t.Run("single", func(t *testing.T) {
			s := NewInMemoryCartStore()
			assert.NoError(t, s.Add(&itemOne))
			count, list := s.GetList()
			assert.Contains(t, list, &itemOne)
			assert.Equal(t, 1, count)
		})
		t.Run("add same", func(t *testing.T) {
			s := NewInMemoryCartStore()
			assert.NoError(t, s.Add(&itemOne))
			count, list := s.GetList()
			assert.Contains(t, list, &itemOne)
			assert.Equal(t, 1, count)

			assert.NoError(t, s.Add(&model.Item{ProductId: itemOne.ProductId, Quanity: 6}))

			count, list = s.GetList()
			assert.Len(t, list, 1)
			assert.Equal(t, 1, count)
			assert.Contains(t, list, &model.Item{
				ProductId: itemOne.ProductId,
				Quanity:   itemOne.Quanity + 6,
			})
		})
		t.Run("add different", func(t *testing.T) {
			s := NewInMemoryCartStore()
			assert.NoError(t, s.Add(&itemOne))
			count, list := s.GetList()
			assert.Contains(t, list, &itemOne)
			assert.Equal(t, 1, count)

			assert.NoError(t, s.Add(&itemTwo))

			count, list = s.GetList()
			assert.Len(t, list, 2)
			assert.Equal(t, 2, count)
			assert.Contains(t, list, &itemOne)
			assert.Contains(t, list, &itemTwo)
		})
	})
}

func TestInMemoryGetById(t *testing.T) {
	t.Run("product", func(t *testing.T) {
		t.Run("invalid", func(t *testing.T) {
			t.Run("empty", func(t *testing.T) {
				s := NewInMemoryProductStore()
				assert.Nil(t, s.GetById(0))
			})
			t.Run("not exist", func(t *testing.T) {
				s := NewInMemoryProductStore()
				s.Products = append(s.Products, &model.Product{Id: 42, Name: "Apple", Price: 4.99})
				assert.Nil(t, s.GetById(24))
			})
		})
		t.Run("valid", func(t *testing.T) {
			s := NewInMemoryProductStore()
			productOne := model.Product{Id: 42, Name: "Apple", Price: 4.99}
			productTwo := model.Product{Id: 24, Name: "Banana", Price: 2.49}
			s.Products = append(s.Products, &productOne)
			s.Products = append(s.Products, &productTwo)
			assert.Equal(t, &productOne, s.GetById(productOne.GetId()))
			assert.Equal(t, &productTwo, s.GetById(productTwo.GetId()))
		})
	})
	t.Run("cart", func(t *testing.T) {
		t.Run("invalid", func(t *testing.T) {
			t.Run("empty", func(t *testing.T) {
				s := NewInMemoryCartStore()
				assert.Nil(t, s.GetById(0))
			})
			t.Run("not exist", func(t *testing.T) {
				s := NewInMemoryCartStore()
				assert.NoError(t, s.Add(&model.Item{ProductId: 42, Quanity: 4}))
				assert.Nil(t, s.GetById(24))
			})
		})
		t.Run("valid", func(t *testing.T) {
			s := NewInMemoryCartStore()
			itemOne := model.Item{ProductId: 42, Quanity: 4}
			itemTwo := model.Item{ProductId: 24, Quanity: 6}
			assert.NoError(t, s.Add(&itemOne))
			assert.NoError(t, s.Add(&itemTwo))
			assert.Equal(t, &itemOne, s.GetById(itemOne.GetProductId()))
			assert.Equal(t, &itemTwo, s.GetById(itemTwo.GetProductId()))
		})
	})
}

func TestInMemoryGetList(t *testing.T) {
	s := NewInMemoryCartStore()
	count, list := s.GetList()
	assert.Empty(t, list)
	assert.Zero(t, count)

	itemOne := model.Item{ProductId: 42, Quanity: 4}
	itemTwo := model.Item{ProductId: 24, Quanity: 16}

	assert.NoError(t, s.Add(&itemOne))

	count, list = s.GetList()
	assert.Len(t, list, 1)
	assert.Equal(t, 1, count)
	assert.Contains(t, list, &itemOne)

	assert.NoError(t, s.Add(&itemTwo))

	count, list = s.GetList()
	assert.Len(t, list, 2)
	assert.Equal(t, 2, count)
	assert.Contains(t, list, &itemOne)
	assert.Contains(t, list, &itemTwo)
	assert.Equal(t, list[0], &itemTwo, "list must be sorted by product_id")
	assert.Equal(t, list[1], &itemOne, "list must be sorted by product_id")
}

func TestInMemoryUpdate(t *testing.T) {
	itemOne := model.Item{ProductId: 42, Quanity: 16}
	t.Run("invalid", func(t *testing.T) {
		t.Run("empty", func(t *testing.T) {
			s := NewInMemoryCartStore()
			assert.Error(t, s.Update(&model.Item{}))
			count, list := s.GetList()
			assert.Empty(t, list)
			assert.Zero(t, count)
		})
		t.Run("zero quanity", func(t *testing.T) {
			s := NewInMemoryCartStore()
			assert.Error(t, s.Update(&model.Item{ProductId: itemOne.ProductId}))
			count, list := s.GetList()
			assert.Empty(t, list)
			assert.Zero(t, count)
		})
	})
	t.Run("valid", func(t *testing.T) {
		t.Run("create", func(t *testing.T) {
			s := NewInMemoryCartStore()
			count, list := s.GetList()
			assert.Empty(t, list)
			assert.Zero(t, count)

			assert.NoError(t, s.Update(&itemOne))

			count, list = s.GetList()
			assert.Contains(t, list, &itemOne)
			assert.Equal(t, 1, count)
		})
		t.Run("change", func(t *testing.T) {
			s := NewInMemoryCartStore()
			count, list := s.GetList()
			assert.Empty(t, list)
			assert.Zero(t, count)

			assert.NoError(t, s.Add(&itemOne))

			count, list = s.GetList()
			assert.Contains(t, list, &itemOne)
			assert.Equal(t, 1, count)

			assert.NoError(t, s.Update(&model.Item{
				ProductId: itemOne.ProductId,
				Quanity:   7,
			}))

			count, list = s.GetList()
			assert.Len(t, list, 1, "length of items store must not be changed")
			assert.Equal(t, 1, count)
			assert.Contains(t, list, &model.Item{
				ProductId: itemOne.ProductId,
				Quanity:   7,
			})
		})
		t.Run("remove", func(t *testing.T) {
			t.Run("zero", func(t *testing.T) {
				s := NewInMemoryCartStore()
				assert.NoError(t, s.Add(&itemOne))

				count, list := s.GetList()
				assert.Contains(t, list, &itemOne)
				assert.Equal(t, 1, count)

				assert.NoError(t, s.Update(&model.Item{
					ProductId: itemOne.ProductId,
					Quanity:   itemOne.Quanity * 0,
				}))

				count, list = s.GetList()
				assert.Empty(t, list)
				assert.Zero(t, count)
			})
			t.Run("less than quanity", func(t *testing.T) {
				s := NewInMemoryCartStore()
				assert.NoError(t, s.Add(&itemOne))

				count, list := s.GetList()
				assert.Contains(t, list, &itemOne)
				assert.Equal(t, 1, count)

				assert.NoError(t, s.Update(&model.Item{
					ProductId: itemOne.ProductId,
					Quanity:   itemOne.Quanity * -1,
				}))
				count, list = s.GetList()
				assert.Empty(t, list)
				assert.Zero(t, count)
			})
		})
	})
}

func TestInMemoryDelete(t *testing.T) {
	itemOne := model.Item{ProductId: 42, Quanity: 4}
	itemTwo := model.Item{ProductId: 24, Quanity: 15}
	itemThree := model.Item{ProductId: 420, Quanity: 7}
	t.Run("invalid", func(t *testing.T) {
		s := NewInMemoryCartStore()
		assert.NoError(t, s.Add(&itemOne))
		count, list := s.GetList()
		assert.Contains(t, list, &itemOne)
		assert.Equal(t, 1, count)

		assert.Error(t, s.Delete(itemTwo.ProductId))

		count, list = s.GetList()
		assert.Contains(t, list, &itemOne)
		assert.Equal(t, 1, count)
	})
	t.Run("valid", func(t *testing.T) {
		s := NewInMemoryCartStore()
		assert.NoError(t, s.Add(&itemOne))
		assert.NoError(t, s.Add(&itemTwo))
		assert.NoError(t, s.Add(&itemThree))

		count, list := s.GetList()
		assert.Contains(t, list, &itemOne)
		assert.Contains(t, list, &itemTwo)
		assert.Contains(t, list, &itemThree)
		assert.Equal(t, count, 3)

		assert.NoError(t, s.Delete(itemThree.ProductId))
		count, list = s.GetList()
		assert.Contains(t, list, &itemOne)
		assert.Contains(t, list, &itemTwo)
		assert.NotContains(t, list, &itemThree)
		assert.Equal(t, 2, count)
	})
}
