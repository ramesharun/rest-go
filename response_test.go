//+build unit

package rest_test

import (
	"encoding/json"
	"errors"
	"github.com/edermanoel94/rest-go"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestResponse_Json(t *testing.T) {

	recorder := httptest.NewRecorder()

	body := []byte("{\"name\": \"cale\"}")

	statusCode := http.StatusOK

	_, _ = rest.Json(recorder, body, statusCode)

	result := recorder.Result()

	defer result.Body.Close()

	bytes, err := ioutil.ReadAll(result.Body)

	if err != nil {
		t.Fatalf("cannot read recorder: %v", err)
	}

	if len(body) != len(bytes) {
		t.Fatalf("size of slice of bytes is different")
	}

	if statusCode != result.StatusCode {
		t.Fatalf("got status %d, but given %d", statusCode, result.StatusCode)
	}

	contentType := "Content-Type"

	if result.Header.Get(contentType) != "application/json" {
		t.Fatalf("should be application/json, got: %s", result.Header.Get(contentType))
	}
}

func TestJsonWithError(t *testing.T) {

	t.Run("should send a signal error", func(t *testing.T) {

		errorThrowed := errors.New("error")
		statusCode := http.StatusNotFound

		recorder := httptest.NewRecorder()

		_, _ = rest.Error(recorder, errorThrowed, statusCode)

		result := recorder.Result()

		defer result.Body.Close()

		bytes, err := ioutil.ReadAll(result.Body)

		if err != nil {
			t.Fatalf("cannot read recorder: %v", err)
		}

		content := map[string]string{}

		err = json.Unmarshal(bytes, &content)

		if err != nil {
			t.Fatalf("couldn't unmarshal: %v", err)
		}

		if errorThrowed.Error() != content["message"] {
			t.Fatalf("expected: %s, got: %s", errorThrowed.Error(), content["message"])
		}
	})
}

func TestJsonWithRedirect(t *testing.T) {

	recorder := httptest.NewRecorder()

	body := []byte("{\"name\": \"cale\"}")
	statusCode := http.StatusOK
	redirect := "http://localhost:8080/tenant/LASA"

	_, _ = rest.Redirect(recorder, body, redirect, statusCode)

	result := recorder.Result()

	defer result.Body.Close()

	_, err := ioutil.ReadAll(result.Body)

	if err != nil {
		t.Fatalf("cannot read recorder: %v", err)
	}

	location := "Location"

	headerLocation := result.Header.Get(location)

	if redirect != headerLocation {
		t.Fatalf("expected a redirect to %s, got: %s", headerLocation, headerLocation)
	}
}

func ExampleJson() {

	product := struct {
		Name  string  `json:"name"`
		Price float32 `json:"price"`
	}{
		Name:  "Smart TV",
		Price: 100.00,
	}

	bytes, _ := json.Marshal(&product)

	recorder := httptest.NewRecorder()

	_, _ = rest.Json(recorder, bytes, http.StatusOK)

	result := recorder.Result()

	defer result.Body.Close()

	_, _ = io.Copy(os.Stdout, result.Body)

	// Output: {"name":"Smart TV","price":100}

}
