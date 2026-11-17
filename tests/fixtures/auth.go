package fixtures

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/auth"
	"api-tests-template/internal/managers/auth/models"
	"api-tests-template/internal/utils"
)

const (
	PrimaryLoginEnv    = "TEST_LOGIN"
	PrimaryPasswordEnv = "TEST_PASSWORD"
	SecondLoginEnv     = "TEST_SECOND_LOGIN"
	SecondPasswordEnv  = "TEST_SECOND_PASSWORD"
)

// Login loads the local environment and authorizes the selected test account.
func Login(t *testing.T, loginEnv, passwordEnv string) models.LoginOkResponse {
	t.Helper()
	utils.LoadEnv()
	login := os.Getenv(loginEnv)
	password := os.Getenv(passwordEnv)
	require.NotEmptyf(t, login, "задайте %s в .env.override", loginEnv)
	require.NotEmptyf(t, password, "задайте %s в .env.override", passwordEnv)
	t.Logf("Подготовка: авторизуем пользователя из %s", loginEnv)
	response := auth.Login(t, login, password)
	require.NotEmpty(t, response.Token, "авторизация не вернула токен")
	require.NotEmpty(t, response.User.Id, "авторизация не вернула id пользователя")
	return response
}
