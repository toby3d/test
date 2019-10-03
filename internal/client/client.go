package client

import (
	"gitlab.com/toby3d/test/internal/model"
	"golang.org/x/xerrors"
	"google.golang.org/grpc"
)

// Client представляет собой простой gRPC клиент
type Client struct {
	model.ShopCartClient
	listener *grpc.ClientConn
}

// ErrClientNotInitialized описывает ошибку инициализации сервера
var ErrClientNotInitialized = xerrors.New("client is not initialized")

// NewClient создаёт новый клиент
func NewClient(addr string) (*Client, error) {
	var c Client

	var err error
	if c.listener, err = grpc.Dial(addr, grpc.WithInsecure()); err != nil {
		return nil, err
	}

	c.ShopCartClient = model.NewShopCartClient(c.listener)
	return &c, nil
}

// Close закрывает все активные соединения с клиентом
func (c *Client) Close() error {
	if c == nil || c.listener == nil {
		return ErrClientNotInitialized
	}
	return c.listener.Close()
}
