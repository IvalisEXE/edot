package proxy

import (
	"io"
	"net/http"
)

func HandleResponse(w http.ResponseWriter, request *http.Request) {
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	w.WriteHeader(response.StatusCode)

	for key, values := range response.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	_, err = io.Copy(w, response.Body)
	if err != nil {
		http.Error(w, "Error copying response body", http.StatusInternalServerError)
	}
}
