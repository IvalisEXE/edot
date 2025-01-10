package proxy

type TransferStockRequest struct {
	WarehouseOriginID            int    `json:"warehouse_origin_id"`
	WarehouseDestinationID       int    `json:"warehouse_destination_id"`
	WarehouseOriginStockItemID   int    `json:"warehouse_origin_stock_item_id"`
	WarehouseOriginStockItemCode string `json:"warehouse_origin_stock_item_code"`
	Number                       int    `json:"number"`
}
