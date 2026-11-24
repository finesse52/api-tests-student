package fixtures

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/advertisement"
	"api-tests-template/internal/managers/advertisement/models"
	"api-tests-template/internal/managers/response"
	"api-tests-template/internal/utils"
)

const (
	defaultPrice            = 1250.50
	defaultQuantity         = 2
	testAdvertisementPrefix = "QA advertisement "
)

type OwnedAdvertisement struct {
	models.Advertisement
	deleted bool
}

// MarkDeleted tells the cleanup hook that this test already removed its advertisement.
func (ad *OwnedAdvertisement) MarkDeleted() {
	ad.deleted = true
}

// NewCreateRequest builds a unique advertisement using photo files stored with the calling scenario.
func NewCreateRequest(t *testing.T, photoFiles ...string) models.CreateRequest {
	t.Helper()
	for _, file := range photoFiles {
		require.FileExists(t, file, "не найден файл фотографии %s", file)
	}
	request := models.CreateRequest{
		Title:       testAdvertisementPrefix + utils.RandomString(16),
		Description: "Created by the automated API test suite",
		Price:       defaultPrice,
		Quantity:    defaultQuantity,
		PhotoFiles:  append([]string(nil), photoFiles...),
	}
	return request
}

// CreateOwnedAdvertisement creates an ad and registers cleanup for this test's owner.
func CreateOwnedAdvertisement(t *testing.T, token string, request models.CreateRequest) *OwnedAdvertisement {
	t.Helper()
	t.Logf("Подготовка: создаём объявление %q", request.Title)
	result := advertisement.Create(t, token, request)
	result.RequireStatus(t, http.StatusCreated)
	var created models.Advertisement
	result.Decode(t, &created)
	require.NotEmpty(t, created.ID, "создание объявления не вернуло id")
	owned := &OwnedAdvertisement{Advertisement: created}
	t.Cleanup(func() {
		if owned.deleted {
			return
		}
		cleanupDeletedEntity(t, "объявление "+created.ID, func() response.Result {
			return advertisement.Delete(t, token, created.ID)
		})
	})
	t.Logf("Результат подготовки: объявление создано, id=%s", created.ID)
	return owned
}
