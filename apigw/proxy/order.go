package proxy

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

var (
	orderUrl = "http://ordersvc:3005/order"
)

type OrderRequest struct {
	ShopID                 int     `json:"shop_id"`
	WarehouseID            int     `json:"warehouse_id"`
	WarehouseStockItemID   int     `json:"warehouse_stock_item_id"`
	WarehouseStockItemCode string  `json:"warehouse_stock_item_code"`
	ProductName            string  `json:"product_name"`
	Price                  float32 `json:"price"`
	Quantity               int     `json:"quantity"`
	PaymentMethod          string  `json:"payment_method"`
	TotalPayment           float32 `json:"total_payment"`
}

func Order(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {

		var bodyReq OrderRequest
		if err := json.NewDecoder(r.Body).Decode(&bodyReq); err != nil {
			http.Error(w, "Invalid body request", http.StatusBadRequest)
			return
		}

		jsonData, err := json.Marshal(bodyReq)
		if err != nil {
			log.Fatalf("Error marshaling JSON: %v", err)
		}

		request, err := http.NewRequest(http.MethodPost, orderUrl, bytes.NewBuffer(jsonData))
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
