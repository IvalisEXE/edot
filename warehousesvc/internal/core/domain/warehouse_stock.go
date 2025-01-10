package domain

type FindAllStockRequest struct {
	ShopID      int `query:"shop_id" validate:"required"`
	WarehouseID int `query:"warehouse_id" validate:"required"`
	QueryParam
}

type FindAllStockResponse struct {
	ID         int     `json:"id"`
	Name       string  `json:"name"`
	Price      float32 `json:"price"`
	StockTotal int     `json:"stock_total"`
}

type FindStockItemByIDRequest struct {
	ShopID               int `query:"shop_id"`
	WarehouseID          int `query:"warehouse_id"`
	WarehouseStockItemID int `query:"warehouse_stock_item_id"`
}

type UpsertStockRequest struct {
	WarehouseID  int     `json:"warehouse_id"`
	ProductID    int     `json:"product_id"`
	ProductName  string  `json:"product_name"`
	ProductPrice float32 `json:"product_price"`
	Number       int     `json:"number"`
}

type TransferStockRequest struct {
	WarehouseOriginID            int    `json:"warehouse_origin_id" validate:"required"`
	WarehouseDestinationID       int    `json:"warehouse_destination_id" validate:"required"`
	WarehouseOriginStockItemID   int    `json:"warehouse_origin_stock_item_id" validate:"required"`
	WarehouseOriginStockItemCode string `json:"warehouse_origin_stock_item_code" validate:"required"`
	Number                       int    `json:"number" validate:"required"`
}

func (t *TransferStockRequest) FindAndReducePreviousStockTotal(previousStockTotal int) int {
	return previousStockTotal - t.Number
}

func (t *TransferStockRequest) FindAndAddPreviousStockTotal(previousStockTotal int) int {
	return previousStockTotal + t.Number
}
