package server

import (
	"net"

	"gitlab.com/toby3d/test/internal/model"
	"golang.org/x/xerrors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Server представляет собой объект gRPC сервера
type Server struct {
	listener net.Listener
	server   *grpc.Server
}

// ErrServerNotInitialized описывает ошибку инициализации сервера
var ErrServerNotInitialized = xerrors.New("server is not initialized")

// NewServer создаёт новое TCP соединение по указанному адресу с указанным набором хендлеров
func NewServer(addr string, handlers model.ShopCartServer) (*Server, error) {
	s := Server{server: grpc.NewServer()}

	var err error
	if s.listener, err = net.Listen("tcp", addr); err != nil {
		return nil, err
	}

	model.RegisterShopCartServer(s.server, handlers)
	reflection.Register(s.server)

	return &s, nil
}

// Start запускает gRPC сервер.
// Возвращает ошибку если была предпринята попытка запуска без предварительной инициализации сервера.
func (s *Server) Start() error {
	if s == nil || s.server == nil || s.listener == nil {
		return ErrServerNotInitialized
	}
	return s.server.Serve(s.listener)
}

// Stop останавливает gRPC сервер.
// Возвращает ошибку если была предпринята попытка остановки без предварительной инициализации сервера.
func (s *Server) Stop() error {
	if s == nil || s.server == nil {
		return ErrServerNotInitialized
	}

	s.server.Stop()
	return nil
}
