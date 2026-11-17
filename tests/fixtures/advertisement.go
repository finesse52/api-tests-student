package fixtures

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/advertisement"
	"api-tests-template/internal/managers/advertisement/models"
	"api-tests-template/internal/managers/response"
	"api-tests-template/internal/utils"
)

const (
	DefaultPrice    = 1250.50
	DefaultQuantity = 2
)

type OwnedAdvertisement struct {
	models.Advertisement
	deleted bool
}

// MarkDeleted tells the cleanup hook that this test already removed its advertisement.
func (ad *OwnedAdvertisement) MarkDeleted() {
	ad.deleted = true
}

// NewCreateRequest builds a unique advertisement and temporary JPEG files for one test.
func NewCreateRequest(t *testing.T, photoCount int) models.CreateRequest {
	t.Helper()
	require.GreaterOrEqual(t, photoCount, 0)
	require.LessOrEqual(t, photoCount, 3)
	request := models.CreateRequest{
		Title:       "QA advertisement " + utils.RandomString(16),
		Description: "Created by the automated API test suite",
		Price:       DefaultPrice,
		Quantity:    DefaultQuantity,
	}
	for i := range photoCount {
		request.PhotoFiles = append(request.PhotoFiles, createPhotoFile(t, i))
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

// createPhotoFile writes a small valid JPEG into the test's temporary directory.
func createPhotoFile(t *testing.T, index int) string {
	t.Helper()
	filePath := filepath.Join(t.TempDir(), fmt.Sprintf("photo-%d.jpg", index))
	file, err := os.Create(filePath)
	require.NoError(t, err, "не удалось создать временную фотографию")
	defer file.Close()

	photo := image.NewRGBA(image.Rect(0, 0, 256, 256))
	shade := uint8(50 + index*50) //nolint:gosec // index is limited to three photos.
	draw.Draw(photo, photo.Bounds(), &image.Uniform{C: color.RGBA{R: shade, G: 80, B: 180, A: 255}}, image.Point{}, draw.Src)
	require.NoError(t, jpeg.Encode(file, photo, &jpeg.Options{Quality: 85}), "не удалось записать JPEG")
	return filePath
}
