package handler

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/toby3d/test/internal/model"
	"gitlab.com/toby3d/test/internal/store"
	"golang.org/x/net/context"
)

var productsManager = store.InMemoryProductStore{Products: []*model.Product{
	&model.Product{Id: 24, Name: "Apple", Price: 4.99},
	&model.Product{Id: 42, Name: "Banana", Price: 2.49},
}}

func TestAdd(t *testing.T) {
	t.Run("invalid", func(t *testing.T) {
		handlerAdd := NewHandler(store.NewInMemoryCartStore(), &productsManager).Add
		t.Run("empty", func(t *testing.T) {
			resp, err := handlerAdd(context.TODO(), &model.AddRequest{})
			assert.Error(t, err)
			assert.False(t, resp.GetOk())
			assert.NotEmpty(t, resp.GetDescription())
			assert.Nil(t, resp.GetResult())
		})
		t.Run("zero quanity", func(t *testing.T) {
			resp, err := handlerAdd(context.TODO(), &model.AddRequest{ProductId: 42})
			assert.Error(t, err)
			assert.False(t, resp.GetOk())
			assert.NotEmpty(t, resp.GetDescription())
			assert.Nil(t, resp.GetResult())
		})
	})
	t.Run("valid", func(t *testing.T) {
		handlerAdd := NewHandler(store.NewInMemoryCartStore(), &productsManager).Add
		itemOne := model.Item{ProductId: 42, Quanity: 5}

		resp, err := handlerAdd(context.TODO(), &model.AddRequest{
			ProductId: itemOne.GetProductId(),
			Quanity:   itemOne.GetQuanity(),
		})
		assert.NoError(t, err)
		assert.True(t, resp.GetOk())
		assert.Empty(t, resp.GetDescription())

		assert.Equal(t, &itemOne, resp.GetItem())

		t.Run("append to exist item", func(t *testing.T) {
			resp, err = handlerAdd(context.TODO(), &model.AddRequest{
				ProductId: itemOne.GetProductId(),
				Quanity:   24,
			})
			assert.NoError(t, err)
			assert.True(t, resp.GetOk())
			assert.Empty(t, resp.GetDescription())

			assert.Equal(t, &model.Item{
				ProductId: itemOne.GetProductId(),
				Quanity:   itemOne.GetQuanity() + 24,
			}, resp.GetItem())
		})
	})
}

func TestGet(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		s := store.NewInMemoryCartStore()
		handlerGet := NewHandler(s, &productsManager).Get

		resp, err := handlerGet(context.TODO(), &model.GetRequest{})
		assert.NoError(t, err)
		assert.True(t, resp.GetOk())
		assert.Empty(t, resp.GetDescription())

		assert.Empty(t, resp.GetCart())
	})
	t.Run("have items", func(t *testing.T) {
		for _, tc := range []struct {
			name       string
			items      []*model.Item
			expCount   int32
			expQuanity int32
			expPrice   float32
		}{{
			name: "2 apples",
			items: []*model.Item{
				&model.Item{ProductId: 24, Quanity: 2},
			},
			expCount:   1,
			expQuanity: 2,
			expPrice:   4.99 * 2,
		}, {
			name: "2 bananas",
			items: []*model.Item{
				&model.Item{ProductId: 42, Quanity: 2},
			},
			expCount:   1,
			expQuanity: 2,
			expPrice:   2.49 * 2,
		}, {
			name: "2 bananas and apples",
			items: []*model.Item{
				&model.Item{ProductId: 42, Quanity: 2},
				&model.Item{ProductId: 24, Quanity: 2},
			},
			expCount:   2,
			expQuanity: 4,
			expPrice:   2.49*2 + 4.99*2,
		}, {
			name: "5 bananas and 3 apples",
			items: []*model.Item{
				&model.Item{ProductId: 42, Quanity: 5},
				&model.Item{ProductId: 24, Quanity: 3},
			},
			expCount:   2,
			expQuanity: 8,
			expPrice:   2.49*5 + 4.99*3,
		}, {
			name: "1+4 bananas and 3 apples",
			items: []*model.Item{
				&model.Item{ProductId: 42, Quanity: 1},
				&model.Item{ProductId: 42, Quanity: 4},
				&model.Item{ProductId: 24, Quanity: 3},
			},
			expCount:   2,
			expQuanity: 8,
			expPrice:   2.49*5 + 4.99*3,
		}, {
			name: "2+3 bananas and 2+4 apples",
			items: []*model.Item{
				&model.Item{ProductId: 42, Quanity: 2},
				&model.Item{ProductId: 24, Quanity: 2},
				&model.Item{ProductId: 42, Quanity: 3},
				&model.Item{ProductId: 24, Quanity: 4},
			},
			expCount:   2,
			expQuanity: 11,
			expPrice:   2.49*5 + 4.99*6,
		}} {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				s := store.NewInMemoryCartStore()
				for _, item := range tc.items {
					assert.NoError(t, s.Add(item))
				}
				handlerGet := NewHandler(s, &productsManager).Get

				resp, err := handlerGet(context.TODO(), &model.GetRequest{})
				assert.NoError(t, err)
				assert.True(t, resp.GetOk())
				assert.Empty(t, resp.GetDescription())

				result := resp.GetCart()
				assert.Equal(t, tc.expCount, result.GetItemsCount())
				assert.Equal(t, tc.expQuanity, result.GetQuanityCount())
				assert.Equal(t, tc.expPrice, result.GetTotalPrice())
			})
		}
	})
}

func TestUpdate(t *testing.T) {
	t.Run("invalid", func(t *testing.T) {
		handlerUpdate := NewHandler(store.NewInMemoryCartStore(), &productsManager).Update
		t.Run("empty", func(t *testing.T) {
			resp, err := handlerUpdate(context.TODO(), &model.UpdateRequest{})
			assert.Error(t, err)
			assert.False(t, resp.GetOk())
			assert.NotEmpty(t, resp.GetDescription())
			assert.Nil(t, resp.GetResult())
		})
		t.Run("not exists", func(t *testing.T) {
			resp, err := handlerUpdate(context.TODO(), &model.UpdateRequest{ProductId: 42})
			assert.Error(t, err)
			assert.False(t, resp.GetOk())
			assert.NotEmpty(t, resp.GetDescription())
			assert.Nil(t, resp.GetResult())
		})
	})
	t.Run("valid", func(t *testing.T) {
		t.Run("create", func(t *testing.T) {
			s := store.NewInMemoryCartStore()
			itemOne := model.Item{ProductId: 42, Quanity: 5}
			handlerUpdate := NewHandler(s, &productsManager).Update

			resp, err := handlerUpdate(
				context.TODO(),
				&model.UpdateRequest{ProductId: itemOne.GetProductId(), Quanity: itemOne.GetQuanity()},
			)
			assert.NoError(t, err)
			assert.True(t, resp.GetOk())
			assert.Empty(t, resp.GetDescription())

			assert.Equal(t, &itemOne, resp.GetItem())
		})
		t.Run("update", func(t *testing.T) {
			s := store.NewInMemoryCartStore()
			itemOne := model.Item{ProductId: 42, Quanity: 5}
			assert.NoError(t, s.Add(&itemOne))
			handlerUpdate := NewHandler(s, &productsManager).Update

			resp, err := handlerUpdate(
				context.TODO(),
				&model.UpdateRequest{ProductId: itemOne.GetProductId(), Quanity: 2},
			)
			assert.NoError(t, err)
			assert.True(t, resp.GetOk())
			assert.Empty(t, resp.GetDescription())

			assert.Equal(t, &model.Item{
				ProductId: itemOne.GetProductId(),
				Quanity:   2,
			}, resp.GetItem())
		})
		t.Run("delete", func(t *testing.T) {
			s := store.NewInMemoryCartStore()
			itemOne := model.Item{ProductId: 42, Quanity: 5}
			assert.NoError(t, s.Add(&itemOne))
			handlerUpdate := NewHandler(s, &productsManager).Update

			resp, err := handlerUpdate(
				context.TODO(),
				&model.UpdateRequest{ProductId: itemOne.GetProductId(), Quanity: 0},
			)
			assert.NoError(t, err)
			assert.True(t, resp.GetOk())
			assert.Empty(t, resp.GetDescription())

			assert.Empty(t, resp.GetResult())
		})
	})
}
func TestRemove(t *testing.T) {
	t.Run("invalid", func(t *testing.T) {
		itemOne := model.Item{ProductId: 42, Quanity: 5}

		s := store.NewInMemoryCartStore()
		assert.NoError(t, s.Add(&itemOne))
		handlerRemove := NewHandler(s, &productsManager).Remove

		t.Run("empty", func(t *testing.T) {
			resp, err := handlerRemove(context.TODO(), &model.RemoveRequest{})
			assert.Error(t, err)
			assert.False(t, resp.GetOk())
			assert.NotEmpty(t, resp.GetDescription())
			assert.Nil(t, resp.GetResult())
			count, list := s.GetList()
			assert.Contains(t, list, &itemOne)
			assert.Equal(t, 1, count)
		})
		t.Run("not exist", func(t *testing.T) {
			resp, err := handlerRemove(context.TODO(), &model.RemoveRequest{ProductId: 24})
			assert.Error(t, err)
			assert.False(t, resp.GetOk())
			assert.NotEmpty(t, resp.GetDescription())
			assert.Nil(t, resp.GetResult())
			count, list := s.GetList()
			assert.Contains(t, list, &itemOne)
			assert.Equal(t, 1, count)
		})
	})
	t.Run("valid", func(t *testing.T) {
		itemOne := model.Item{ProductId: 42, Quanity: 5}

		s := store.NewInMemoryCartStore()
		assert.NoError(t, s.Add(&itemOne))
		handlerRemove := NewHandler(s, &productsManager).Remove

		resp, err := handlerRemove(
			context.TODO(), &model.RemoveRequest{ProductId: itemOne.GetProductId()},
		)
		assert.NoError(t, err)
		assert.True(t, resp.GetOk())
		assert.Empty(t, resp.GetDescription())
		assert.Empty(t, resp.GetResult())
		count, list := s.GetList()
		assert.NotContains(t, list, &itemOne)
		assert.Zero(t, count)
	})
}
