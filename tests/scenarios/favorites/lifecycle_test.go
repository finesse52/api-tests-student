package favorites

import (
	"testing"

	"api-tests-template/internal/managers/advertisement"
	favoriteManager "api-tests-template/internal/managers/favorites"
	"api-tests-template/internal/managers/response"
	"api-tests-template/tests/fixtures"
)

func TestFavoriteLifecycleAndAdvertisementDeletion(t *testing.T) {
	viewer := fixtures.Login(t, fixtures.PrimaryLoginEnv, fixtures.PrimaryPasswordEnv)
	owner := fixtures.Login(t, fixtures.SecondLoginEnv, fixtures.SecondPasswordEnv)
	request := fixtures.NewCreateRequest(t, 1)
	created := fixtures.CreateOwnedAdvertisement(t, owner.Token, request)
	fixtures.RequireFavorite(t, viewer.Token, request.Title, created.ID, false)

	t.Log("Шаг: добавляем чужое объявление в избранное")
	fixtures.AddFavorite(t, viewer.Token, request.Title, created.ID)
	fixtures.RequireFavorite(t, viewer.Token, request.Title, created.ID, true)
	t.Log("Результат: объявление появилось в избранном")

	t.Log("Шаг: удаляем объявление из избранного")
	fixtures.DeleteUntilCompleted(t, "из избранного "+created.ID, func() response.Result {
		return favoriteManager.Delete(t, viewer.Token, created.ID)
	})
	fixtures.RequireFavorite(t, viewer.Token, request.Title, created.ID, false)
	t.Log("Результат: объявление исчезло из избранного")

	t.Log("Шаг: добавляем объявление снова и удаляем его владельцем")
	fixtures.AddFavorite(t, viewer.Token, request.Title, created.ID)
	fixtures.RequireFavorite(t, viewer.Token, request.Title, created.ID, true)
	fixtures.DeleteUntilCompleted(t, "объявления "+created.ID, func() response.Result {
		return advertisement.Delete(t, owner.Token, created.ID)
	})
	fixtures.RequireAdvertisementDeleted(t, created.ID)
	created.MarkDeleted()
	fixtures.RequireFavorite(t, viewer.Token, request.Title, created.ID, false)
	t.Log("Результат: удалённое объявление исчезло из избранного пользователя")
}
