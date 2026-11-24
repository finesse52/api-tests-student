package advertisementCreate

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/advertisement"
	"api-tests-template/internal/managers/advertisement/models"
	"api-tests-template/internal/managers/response"
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
		{"пустое description", "description is required", func(r *models.CreateRequest) { r.Description = "" }},
		{"нет обязательных photos", "at least 1 photo is required", func(r *models.CreateRequest) { r.PhotoFiles = nil }},
		{"отрицательная цена", "price must be between 0 and 1000000 inclusive", func(r *models.CreateRequest) { r.Price = -1 }},
		{"строковая цена", "cannot convert string value invalid to float", func(r *models.CreateRequest) { r.PriceText = "invalid" }},
		{"отрицательное количество", "quantity must be between 1 and 100 inclusive", func(r *models.CreateRequest) { r.Quantity = -1 }},
		{"строковое количество", "cannot convert string value to int", func(r *models.CreateRequest) { r.QuantityText = "invalid" }},
		{"слишком длинный title", "title must not exceed 50 characters provided 256", func(r *models.CreateRequest) { r.Title = strings.Repeat("A", 256) }},
		{"неподдерживаемый формат фото", "photo 1: file type not allowed: extension \".txt\"", func(r *models.CreateRequest) { r.PhotoFiles = []string{unsupportedPhoto} }},
		{"более трёх фото", "maximum 3 photos allowed", func(r *models.CreateRequest) {
			r.PhotoFiles = []string{validPhoto1, validPhoto2, validPhoto3, validPhoto4}
		}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			request := fixtures.NewCreateRequest(t, validPhoto1)
			tc.change(&request)
			t.Logf("Подготовка: сформирован запрос с ошибкой: %s", tc.name)
			result := fixtures.RetryServiceUnavailable(t, "POST /advertisement с некорректными данными", func() response.Result {
				return advertisement.CreateExpectedError(t, owner.Token, request)
			})
			t.Log("Шаг: проверяем отказ сервера и тело ошибки")
			apiError := result.RequireError(t, http.StatusBadRequest, "invalid_request")
			require.Equal(t, tc.message, apiError.Message, "неверное сообщение об ошибке валидации")
			t.Logf("Результат: HTTP 400, error=%q, message=%q", apiError.Error, apiError.Message)
		})
	}
}
