package routing

import (
	"apigw/middleware"
	"apigw/proxy"
	"net/http"
)

func Routing() {
	// dashboard
	http.Handle("/apigw/users/health", middleware.AuthDashboard(http.HandlerFunc(proxy.HealthHandler)))
	http.Handle("/apigw/warehouse/unactive/", middleware.AuthDashboard(http.HandlerFunc(proxy.UnactiveWarehouseHandler)))
	http.Handle("/apigw/warehouse/stock-transfer", middleware.AuthDashboard(http.HandlerFunc(proxy.TransferStockHandler)))

	//customer
	http.Handle("/apigw/warehouse/stock", middleware.AuthCustomer(http.HandlerFunc(proxy.FindAllStockHandler)))
	http.Handle("/apigw/order", middleware.AuthCustomer(http.HandlerFunc(proxy.Order)))

	// dashboard & customer
	http.Handle("/apigw/users/login", http.HandlerFunc(proxy.LoginHandler))
}
