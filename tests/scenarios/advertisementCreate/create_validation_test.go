package advertisementCreate

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/advertisement"
	"api-tests-template/internal/managers/advertisement/models"
	"api-tests-template/tests/fixtures"
)

func TestCreateAdvertisementRejectsInvalidFields(t *testing.T) {
	owner := fixtures.Login(t, fixtures.PrimaryLoginEnv, fixtures.PrimaryPasswordEnv)
	cases := []struct {
		name    string
		message string
		change  func(*models.CreateRequest)
	}{
		{"пустой title", "title is required", func(r *models.CreateRequest) { r.Title = "" }},
		{"нет обязательных photos", "at least 1 photo is required", func(r *models.CreateRequest) { r.PhotoFiles = nil }},
		{"отрицательная цена", "price must be between 0 and 1000000 inclusive", func(r *models.CreateRequest) { r.Price = -1 }},
		{"слишком длинный title", "title must not exceed 50 characters provided 1001", func(r *models.CreateRequest) { r.Title = strings.Repeat("A", 1001) }},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			request := fixtures.NewCreateRequest(t, 1)
			tc.change(&request)
			t.Logf("Подготовка: сформирован запрос с ошибкой: %s", tc.name)
			result := advertisement.Create(t, owner.Token, request)
			t.Log("Шаг: проверяем отказ сервера и тело ошибки")
			apiError := result.RequireError(t, http.StatusBadRequest, "invalid_request")
			require.Equal(t, tc.message, apiError.Message, "неверное сообщение об ошибке валидации")
			t.Logf("Результат: HTTP 400, error=%q, message=%q", apiError.Error, apiError.Message)
		})
	}
}
