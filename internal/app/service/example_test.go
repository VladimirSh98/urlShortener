package service

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"net/http/httptest"
)

func Example() {
	// Здесь должна быть инициализация зависимостей:
	// - Настройка логгера (zap)
	// - Подключение к БД
	// - Создание репозитория и сервиса

	// Для примера создаем тестовый роутер
	router := chi.NewMux()
	router.Get("/ping", func(w http.ResponseWriter, r *http.Request) {})

	// Запускаем тестовый сервер
	testServer := httptest.NewServer(router)
	defer testServer.Close()

	// Output:
}
