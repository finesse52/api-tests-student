package advertisementCreate

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/advertisement"
	"api-tests-template/internal/managers/advertisement/models"
	"api-tests-template/internal/managers/advertisements"
	"api-tests-template/tests/fixtures"
)

func TestCreateAdvertisementWithAllFields(t *testing.T) {
	owner := fixtures.Login(t, fixtures.PrimaryLoginEnv, fixtures.PrimaryPasswordEnv)
	request := fixtures.NewCreateRequest(t, 3)
	created := fixtures.CreateOwnedAdvertisement(t, owner.Token, request)

	t.Log("Шаг: проверяем данные ответа на создание")
	require.Equal(t, request.Title, created.Title, "название созданного объявления")
	require.Equal(t, request.Description, created.Description, "описание созданного объявления")
	require.Equal(t, request.Price, created.Price, "цена созданного объявления")
	require.Equal(t, request.Quantity, created.Quantity, "количество созданного объявления")
	require.Equal(t, owner.User.Id, created.UserID, "владелец созданного объявления")
	require.Len(t, created.Photos, len(request.PhotoFiles), "количество загруженных фото")
	t.Log("Результат: сервер сохранил все параметры и три фотографии")

	t.Log("Шаг: читаем объявление через GET /advertisement?id=...")
	getResult := advertisement.Get(t, created.ID)
	getResult.RequireStatus(t, http.StatusOK)
	var fetched models.Advertisement
	getResult.Decode(t, &fetched)
	require.Equal(t, created.ID, fetched.ID, "id объявления после чтения")
	require.Equal(t, request.Title, fetched.Title, "название объявления после чтения")
	require.Equal(t, request.Description, fetched.Description, "описание объявления после чтения")
	require.Equal(t, request.Price, fetched.Price, "цена объявления после чтения")
	require.Equal(t, request.Quantity, fetched.Quantity, "количество объявления после чтения")
	require.Equal(t, owner.User.Id, fetched.UserID, "владелец объявления после чтения")
	require.Len(t, fetched.Photos, len(request.PhotoFiles), "фото объявления после чтения")
	t.Log("Результат: GET возвращает переданные параметры")

	t.Log("Шаг: ищем объявление по полному названию")
	found := false
	deadline := time.Now().Add(10 * time.Second)
	for !found && time.Now().Before(deadline) {
		list := advertisements.GetAdvertisementsList(t, owner.Token, request.Title)
		for _, item := range list.Items {
			if item.Title == request.Title {
				found = true
				break
			}
		}
		if !found {
			time.Sleep(250 * time.Millisecond)
		}
	}
	require.Truef(t, found, "объявление %s не появилось в поиске по полному названию %q", created.ID, request.Title)
	t.Log("Результат: объявление найдено в каталоге")

	t.Log("Шаг: читаем фотографии через GET /advertisements/{id}/photos")
	photosResult := advertisement.GetPhotos(t, owner.Token, created.ID)
	photosResult.RequireStatus(t, http.StatusOK)
	var storedPhotos []models.Photo
	photosResult.Decode(t, &storedPhotos)
	require.Len(t, storedPhotos, len(request.PhotoFiles), "количество фото на сервере")
	createdPhotoIDs := make([]string, 0, len(created.Photos))
	storedPhotoIDs := make([]string, 0, len(storedPhotos))
	for _, photo := range created.Photos {
		createdPhotoIDs = append(createdPhotoIDs, photo.ID)
	}
	for _, photo := range storedPhotos {
		require.NotEmpty(t, photo.ID, "id загруженной фотографии")
		require.NotEmpty(t, photo.URL, "URL загруженной фотографии")
		storedPhotoIDs = append(storedPhotoIDs, photo.ID)
	}
	require.ElementsMatch(t, createdPhotoIDs, storedPhotoIDs, "фото после создания и через GET не совпадают")
	t.Log("Результат: все загруженные фотографии доступны через API")
}
