package proxy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

var (
	findAllStockUrl      = "http://warehousesvc:3004/warehouse/stock"
	transferStockUrl     = "http://warehousesvc:3004/warehouse/stock-transfer"
	unactiveWarehouseUrl = "http://warehousesvc:3004/warehouse/unactive"
)

func FindAllStockHandler(w http.ResponseWriter, r *http.Request) {
	shopID := r.URL.Query().Get("shop_id")
	warehouseID := r.URL.Query().Get("warehouse_id")
	perPage := r.URL.Query().Get("per_page")
	page := r.URL.Query().Get("page")
	sort := r.URL.Query().Get("sort")

	url := fmt.Sprintf("%s?shop_id=%s&warehouse_id=%s&per_page=%s&page=%s&sort=%s",
		findAllStockUrl,
		shopID,
		warehouseID,
		perPage,
		page,
		sort)
	request, err := http.NewRequest(http.MethodGet, url, nil)
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

	HandleResponse(w, request)
}

func UnactiveWarehouseHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) < 4 {
			http.Error(w, "Invalid URL", http.StatusBadRequest)
			return
		}

		warehouseID := parts[4]

		request, err := http.NewRequest(http.MethodPost, unactiveWarehouseUrl+"/"+warehouseID, nil)
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

		HandleResponse(w, request)
	}
}

func TransferStockHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {

		var bodyReq TransferStockRequest
		if err := json.NewDecoder(r.Body).Decode(&bodyReq); err != nil {
			http.Error(w, "Invalid body request", http.StatusBadRequest)
			return
		}

		jsonData, err := json.Marshal(bodyReq)
		if err != nil {
			log.Fatalf("Error marshaling JSON: %v", err)
		}

		request, err := http.NewRequest(http.MethodPost, transferStockUrl, bytes.NewBuffer(jsonData))
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		authHeaders, ok := r.Context().Value("authHeaders").(http.Header)
		if !ok {
			http.Error(w, "Unable to get Auth headers", http.StatusInternalServerError)
			return
		}

		request.Header.Add("Content-Type", "application/json")
		for key, values := range authHeaders {
			for _, value := range values {
				request.Header.Add(key, value)
			}
		}

		HandleResponse(w, request)
	}
}
