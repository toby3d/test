package store

import (
	"sort"
	"sync"

	"gitlab.com/toby3d/test/internal/model"
	"golang.org/x/xerrors"
)

type (
	// InMemoryProductStore представляет собой объект хранилища продуктов доступный только для чтения
	InMemoryProductStore struct {
		mutex    sync.RWMutex
		Products []*model.Product
	}

	// InMemoryCartStore представляет собой простой менеджер объектов корзины.
	InMemoryCartStore struct {
		mutex sync.RWMutex
		items []*model.Item
	}
)

var (
	ErrNoProductId = xerrors.New("product_id not provided")
	ErrZeroQuanity = xerrors.New("item quanity is zero")
	ErrNotExist    = xerrors.New("item not exists or already removed")
)

// NewInMemoryProductStore создаёт новый менеджер продуктов
func NewInMemoryProductStore() *InMemoryProductStore {
	return &InMemoryProductStore{mutex: sync.RWMutex{}}
}

// GetById возвращает информацию о продукте по его ID, если он существует
func (imps *InMemoryProductStore) GetById(id uint64) *model.Product {
	for _, product := range imps.Products {
		if product.GetId() != id {
			continue
		}
		return product
	}
	return nil
}

// NewInMemoryCartStore создаёт новый менеджер объектов корзины
func NewInMemoryCartStore() *InMemoryCartStore { return &InMemoryCartStore{mutex: sync.RWMutex{}} }

// Add добавляет новый объект в корзину.
//
// * Если объекта не существует, то он будет добавлен в список.
// * Если объект уже существует в корзине, то его количество в корзине будет увеличено.
// * Ошибка будет возвращена если не был указан ID продукта или его количество равно или меньше ноля.
//
// BUG(toby3d): InMemoryStore не проверяет наличие продукта по его ID, подразумевая, что он всегда существует в базе.
func (ims *InMemoryCartStore) Add(i *model.Item) error {
	switch {
	case i == nil, i.GetProductId() <= 0:
		return ErrNoProductId
	case i.GetQuanity() <= 0:
		return ErrZeroQuanity
	}

	if item := ims.GetById(i.GetProductId()); item != nil {
		ims.mutex.Lock()
		item.Quanity += i.GetQuanity()
		ims.mutex.Unlock()
		return nil
	}

	ims.mutex.Lock()
	ims.items = append(ims.items, &model.Item{
		ProductId: i.GetProductId(),
		Quanity:   i.GetQuanity(),
	})
	sort.Slice(ims.items, func(i, j int) bool {
		return ims.items[i].GetProductId() < ims.items[j].GetProductId()
	})
	ims.mutex.Unlock()
	return nil
}

// GetById возвращает объект корзины по ID его продукта (если существует).
func (ims *InMemoryCartStore) GetById(id uint64) *model.Item {
	if id == 0 {
		return nil
	}
	ims.mutex.RLock()
	defer ims.mutex.RUnlock()
	for _, item := range ims.items {
		if item.GetProductId() != id {
			continue
		}
		return item
	}
	return nil
}

// GetList возвращает массив всех доступных в корзине объектов.
func (ims *InMemoryCartStore) GetList() (int, []*model.Item) {
	ims.mutex.RLock()
	defer ims.mutex.RUnlock()
	return len(ims.items), ims.items
}

// Update обновляет конкретный объект в корзине.
//
// * Если объекта с указанным ProductId ещё не существует в корзине, то он будет добавлен.
// * Если объект уже существует в корзине, то его количество будет перезаписано.
// * Если желаемое количество объекта равна нулю или отрицательно, то объект будет удалён из корзины.
func (ims *InMemoryCartStore) Update(i *model.Item) error {
	if item := ims.GetById(i.GetProductId()); item != nil {
		if i.GetQuanity() <= 0 {
			return ims.Delete(i.GetProductId())
		}

		ims.mutex.Lock()
		item.Quanity = i.GetQuanity()
		ims.mutex.Unlock()
		return nil
	}
	return ims.Add(i)
}

// Delete удаляет объект из корзины по его ProductId (если он существует).
func (ims *InMemoryCartStore) Delete(id uint64) error {
	if item := ims.GetById(id); item == nil {
		return ErrNotExist
	}
	ims.mutex.Lock()
	defer ims.mutex.Unlock()
	for i := range ims.items {
		if ims.items[i].GetProductId() != id {
			continue
		}
		// NOTE(toby3d): см. https://github.com/golang/go/wiki/SliceTricks
		ims.items[i] = ims.items[len(ims.items)-1]
		ims.items[len(ims.items)-1] = nil
		ims.items = ims.items[:len(ims.items)-1]
		break
	}
	sort.Slice(ims.items, func(i, j int) bool {
		return ims.items[i].GetProductId() < ims.items[j].GetProductId()
	})
	return nil
}
