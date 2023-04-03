package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/igorrnk/ypdiploma.git/internal/auth"
	"github.com/igorrnk/ypdiploma.git/internal/configs"
	"github.com/igorrnk/ypdiploma.git/internal/model"
	"net/http"
)

type Server struct {
	http.Server
	config   *configs.ServerConfigType
	servicer model.Servicer
}

func NewServer(config *configs.ServerConfigType, servicer model.Servicer) *Server {
	server := &Server{
		Server: http.Server{
			Addr: config.ServerAddress,
		},
		config:   config,
		servicer: servicer,
	}
	server.Handler = newChiRouter(server)
	return server
}

func newChiRouter(s *Server) http.Handler {
	r := chi.NewRouter()

	// Public routes
	r.Group(func(r chi.Router) {
		r.Use(middleware.Logger)
		r.Post("/api/user/register", s.PostRegister)
		r.Post("/api/user/login", s.PostLogin)
	})

	r.Group(func(r chi.Router) {
		r.Use(middleware.Logger)
		r.Use(auth.AuthenticatorJWT)
		r.Post("/api/user/orders", s.PostOrders)
		r.Get("/api/user/orders", s.GetOrders)
		r.Get("/api/user/balance", s.GetBalance)
		r.Post("/api/user/balance/withdraw", s.PostWithdraw)
		r.Get("/api/user/withdrawals", s.GetWithDrawals)
	})
	return r
	/*
		POST /api/user/register — регистрация пользователя;
		POST /api/user/login — аутентификация пользователя;
		POST /api/user/orders — загрузка пользователем номера заказа для расчёта;
		GET /api/user/orders — получение списка загруженных пользователем номеров заказов, статусов их обработки и информации о начислениях;
		GET /api/user/balance — получение текущего баланса счёта баллов лояльности пользователя;
		POST /api/user/balance/withdraw — запрос на списание баллов с накопительного счёта в счёт оплаты нового заказа;
		GET /api/user/withdrawals — получение информации о выводе средств с накопительного счёта пользователем.
	*/
}
