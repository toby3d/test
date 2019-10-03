package handler

import (
	"math"

	"gitlab.com/toby3d/test/internal/model"
	"gitlab.com/toby3d/test/internal/model/store"
	"golang.org/x/net/context"
)

// Handler представляет собой объект хендлеров с хранилищем данных
type Handler struct {
	cartManager   store.CartManager
	productReader store.ProductReader
}

// NewHandler создаёт хендлеры сервера с указанным хранилищем
func NewHandler(cartManager store.CartManager, productReader store.ProductReader) *Handler {
	return &Handler{
		cartManager:   cartManager,
		productReader: productReader,
	}
}

// Add добавляет объект в хранилище (если его не существует) или обновляет количество существующего объекта
func (h *Handler) Add(ctx context.Context, req *model.AddRequest) (*model.Response, error) {
	var resp model.Response

	if err := h.cartManager.Add(&model.Item{
		ProductId: req.GetProductId(),
		Quanity:   req.GetQuanity(),
	}); err != nil {
		resp.Description = err.Error()
		return &resp, err
	}

	resp.Ok = true
	resp.Result = &model.Response_Item{Item: h.cartManager.GetById(req.GetProductId())}
	return &resp, nil
}

func (h *Handler) Get(ctx context.Context, req *model.GetRequest) (*model.Response, error) {
	var resp model.Response

	var result model.Cart
	count, items := h.cartManager.GetList()
	result.Items = items
	result.ItemsCount = int32(count)
	for _, item := range items {
		result.QuanityCount += item.GetQuanity()
		product := h.productReader.GetById(item.GetProductId())
		result.TotalPrice += product.GetPrice() * float32(item.GetQuanity())
	}
	// NOTE(toby3d): округляем до двух знаков после запятой
	result.TotalPrice = float32(math.Round(float64(result.GetTotalPrice())*100) / 100)

	resp.Ok = true
	resp.Result = &model.Response_Cart{Cart: &result}
	return &resp, nil
}

// Update обновляет конкретные товары в корзине.
func (h *Handler) Update(ctx context.Context, req *model.UpdateRequest) (*model.Response, error) {
	var resp model.Response

	if err := h.cartManager.Update(&model.Item{
		ProductId: req.GetProductId(),
		Quanity:   req.GetQuanity(),
	}); err != nil {
		resp.Description = err.Error()
		return &resp, err
	}

	// NOTE(toby3d): Если количество отрицательно, то Result должен быть пустой как и в случае Remove
	if req.GetQuanity() > 0 {
		resp.Result = &model.Response_Item{Item: h.cartManager.GetById(req.GetProductId())}
	}

	resp.Ok = true
	return &resp, nil
}

// Remove удаляет конкретные товары в корзине вне зависимости от их количества.
func (h *Handler) Remove(ctx context.Context, req *model.RemoveRequest) (*model.Response, error) {
	var resp model.Response

	if err := h.cartManager.Delete(req.GetProductId()); err != nil {
		resp.Description = err.Error()
		return &resp, err
	}

	resp.Ok = true
	return &resp, nil
}
