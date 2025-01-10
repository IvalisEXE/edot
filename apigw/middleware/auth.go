package middleware

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
)

var (
	authDasboardUrl = "http://usersvc:3001/users/dashboard-authenticate"
	authCustomerUrl = "http://usersvc:3001/users/customer-authenticate"
)

type User struct {
	Data struct {
		ID     int    `json:"id"`
		ShopID int    `json:"shop_id"`
		Name   string `json:"name,omitempty"`
		Phone  string `json:"phone"`
	} `json:"data"`
}

func AuthDashboard(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		token = strings.TrimPrefix(token, "Bearer ")
		header, ok := verifyAuthDashboard(token)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "authHeaders", header)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func verifyAuthDashboard(token string) (http.Header, bool) {
	if token == "" {
		return nil, false
	}

	request, err := http.NewRequest(http.MethodGet, authDasboardUrl, nil)
	if err != nil {
		return nil, false
	}

	request.Header.Add("Authorization", "Bearer "+token)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, false
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, false
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, false
	}

	var user User
	if err = json.Unmarshal(body, &user); err != nil {
		return nil, false
	}

	header := http.Header{}
	header.Add("X-User-Id", strconv.Itoa(user.Data.ID))
	header.Add("X-Shop-Id", strconv.Itoa(user.Data.ShopID))
	header.Add("X-Username", user.Data.Name)
	header.Add("X-Phone", user.Data.Phone)

	return header, true
}

func AuthCustomer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		token = strings.TrimPrefix(token, "Bearer ")
		header, ok := verifyAuthCustomer(token)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "authHeaders", header)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func verifyAuthCustomer(token string) (http.Header, bool) {
	if token == "" {
		return nil, false
	}

	request, err := http.NewRequest(http.MethodGet, authCustomerUrl, nil)
	if err != nil {
		return nil, false
	}

	request.Header.Add("Authorization", "Bearer "+token)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, false
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, false
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, false
	}

	var user User
	if err = json.Unmarshal(body, &user); err != nil {
		return nil, false
	}

	header := http.Header{}
	header.Add("X-User-Id", strconv.Itoa(user.Data.ID))
	header.Add("X-Shop-Id", strconv.Itoa(user.Data.ShopID))
	header.Add("X-Username", user.Data.Name)
	header.Add("X-Phone", user.Data.Phone)

	return header, true
}
