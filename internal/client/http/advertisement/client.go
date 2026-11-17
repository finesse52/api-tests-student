package advertisement

import (
	"net/http"
	"strconv"
	"testing"

	"api-tests-template/internal/constants/path"
	apiRunner "api-tests-template/internal/helpers/api-runner"
	"api-tests-template/internal/managers/advertisement/models"
)

// HttpCreate sends a multipart advertisement request, including each photo under the same field name.
func HttpCreate(t *testing.T, token string, request models.CreateRequest, retryTransport bool) *http.Response {
	t.Helper()
	runner := apiRunner.GetRunner().Auth(token)
	call := runner.CreateWithTimeout()
	if retryTransport {
		call = runner.CreateWithTransientTransportRetry()
	}
	createRequest := call.Post(path.AdvertisementPath).
		MultipartFormData("title", request.Title).
		MultipartFormData("description", request.Description).
		MultipartFormData("price", strconv.FormatFloat(request.Price, 'f', 2, 64)).
		MultipartFormData("quantity", strconv.Itoa(request.Quantity))
	for _, photo := range request.PhotoFiles {
		createRequest = createRequest.MultipartFile("photos", photo)
	}
	return createRequest.Expect(t).End().Response
}

// HttpGet retrieves an advertisement by its query-string ID.
func HttpGet(t *testing.T, id string) *http.Response {
	t.Helper()
	return apiRunner.GetRunner().CreateWithTimeout().Get(path.AdvertisementPath).
		Query("id", id).Expect(t).End().Response
}

// HttpDelete removes an advertisement owned by the authorized user.
func HttpDelete(t *testing.T, token, id string) *http.Response {
	t.Helper()
	return apiRunner.GetRunner().Auth(token).CreateWithTimeout().Delete(path.AdvertisementPath).
		Query("id", id).Expect(t).End().Response
}

// HttpGetPhotos retrieves the photos stored for an advertisement.
func HttpGetPhotos(t *testing.T, token, id string) *http.Response {
	t.Helper()
	return apiRunner.GetRunner().Auth(token).CreateWithTimeout().Get(path.AdvertisementsPath + "/" + id + "/photos").
		Expect(t).End().Response
}
