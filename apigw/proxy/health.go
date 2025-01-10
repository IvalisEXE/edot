package proxy

import (
	"io"
	"net/http"
)

var (
	healthUrl = "http://usersvc:3001/health"
)

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	request, err := http.NewRequest(http.MethodGet, healthUrl, nil)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	authHeaders, ok := r.Context().Value("authHeaders").(http.Header)
	if !ok {
		http.Error(w, "Unable to get Auth headers", http.StatusInternalServerError)
		return
	}

	for key, values := range authHeaders {
		for _, value := range values {
			request.Header.Add(key, value)
		}
	}

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
