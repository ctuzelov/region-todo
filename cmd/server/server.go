package server

import (
	"context"
	"net/http"
	"time"

	"github.com/spf13/viper"
)

type Server struct {
	httpServer *http.Server
}

// * Запуск HTTP-сервера на указанном порту с заданным обработчиком.
// * Возвращает ошибку, если не удается запустить сервер.
func (s *Server) Run(port string, handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:           ":" + port,
		Handler:        handler,
		MaxHeaderBytes: 1 << 20, // 1MB
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
	}
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func InitConfig() error {
	viper.AddConfigPath("pkg/configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
