package repository

import (
	"time"

	"go.elastic.co/apm"
	"gorm.io/gorm"

	"warehousesvc/internal/core/domain"
	"warehousesvc/internal/core/port"
	"warehousesvc/pkg/common"
)

type repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) port.Repository {
	return &repository{db}
}

func (r *repository) FindAllStock(ctx common.ServiceContextManager, filter domain.FindAllStockRequest) ([]domain.FindAllStockResponse, error) {
	span, context := apm.StartSpan(ctx.GetContext(), "repository.warehouse.FindAllStock", "repository")
	defer span.End()

	ctx.SetContext(context)

	var warehouseStocks []WarehouseStockItem
	tx := r.db.WithContext(ctx.GetContext()).Model(&WarehouseStockItem{}).
		Preload("Warehouse").
		Where("warehouse_stock_item.shop_id = ?", filter.ShopID).
		Where("warehouse_stock_item.warehouse_id = ?", filter.WarehouseID)

	// Set pagination and sorting
	if err := tx.Scopes(
		common.Paginate(ctx, warehouseStocks, tx),
		common.Sorting(ctx, tx),
	).
		Find(&warehouseStocks).Error; err != nil {
		return nil, err
	}

	var result []domain.FindAllStockResponse
	for _, b := range warehouseStocks {
		warehousestock := b.ToStockList()
		result = append(result, warehousestock)
	}

	return result, nil
}

func (r *repository) FindStockByID(ctx common.ServiceContextManager, id int) (*domain.Warehouse, error) {
	span, context := apm.StartSpan(ctx.GetContext(), "repository.warehouse.FindStockItemByID", "repository")
	defer span.End()

	ctx.SetContext(context)

	query := "SELECT * FROM warehouse WHERE id = ? FOR UPDATE"
	tx := ctx.GetDB()

	if tx == nil {
		tx = r.db
		query = "SELECT * FROM warehouse WHERE id = ?"
	}

	var stock Warehouse
	if err := tx.Raw(query, id).Scan(&stock).Error; err != nil {
		return nil, err
	}

	if stock.ID == 0 {
		return nil, nil
	}

	return stock.ToDomain(), nil
}

func (r *repository) FindStockItemByID(ctx common.ServiceContextManager, itemID, warehouseID int) (*domain.WarehouseStockItem, error) {
	span, context := apm.StartSpan(ctx.GetContext(), "repository.warehouse.FindStockItemByID", "repository")
	defer span.End()

	ctx.SetContext(context)
	query := `
		SELECT 
			ws.id,
			ws.shop_id,
			ws.warehouse_id,
			ws.name,
			ws.price,
			ws.stock_total,
			ws.stocked_at,
			ws.created_at,
			ws.created_by,
			ws.updated_at,
			ws.updated_by,
			w.id as warehouse_id,
			w.shop_id as warehouse_shop_id,
			w.name as warehouse_name,
			w.address as warehouse_address,
			w.is_active as warehouse_is_active,
			w.stock_total as warehouse_stock_total,
			w.stocked_at as wareshouse_stocked_at,
			w.created_at,
			w.created_by,
			w.updated_at,
			w.updated_by,
			w.deleted_at,
			w.deleted_by
		FROM warehouse_stock_item ws
		JOIN warehouse w ON ws.warehouse_id = w.id 
		WHERE ws.id = ? AND ws.warehouse_id = ? FOR UPDATE
	`
	tx := ctx.GetDB()
	if tx == nil {
		tx = r.db
		query = `
		SELECT 
			ws.id,
			ws.shop_id,
			ws.warehouse_id,
			ws.name,
			ws.price,
			ws.stock_total,
			ws.stocked_at,
			ws.created_at,
			ws.created_by,
			ws.updated_at,
			ws.updated_by,
			w.id as warehouse_id,
			w.shop_id as warehouse_shop_id,
			w.name as warehouse_name,
			w.address as warehouse_address,
			w.is_active as warehouse_is_active,
			w.stock_total as warehouse_stock_total,
			w.stocked_at as wareshouse_stocked_at,
			w.created_at,
			w.created_by,
			w.updated_at,
			w.updated_by,
			w.deleted_at,
			w.deleted_by
		FROM warehouse_stock_item ws
		JOIN warehouse w ON ws.warehouse_id = w.id 
		WHERE ws.id = ? AND ws.warehouse_id = ? FOR UPDATE
	`
	}

	var stock WarehouseAndWarehouseStockItem
	if err := tx.Raw(query, itemID, warehouseID).Scan(&stock).Error; err != nil {
		return nil, err
	}

	if stock.ID == 0 {
		return nil, nil
	}

	return stock.ToDomain(), nil
}

func (r *repository) UpdateStock(ctx common.ServiceContextManager, id, totalStock int) error {
	span, context := apm.StartSpan(ctx.GetContext(), "repository.warehouse.UpdateStock", "repository")
	defer span.End()

	ctx.SetContext(context)

	var (
		now = time.Now()
	)

	// Get db transaction if exist
	tx := ctx.GetDB()
	if tx == nil {
		tx = r.db
	}

	tx = tx.WithContext(ctx.GetContext()).Where("id = ?", id)
	if err := tx.WithContext(ctx.GetContext()).Model(Warehouse{}).Updates(map[string]interface{}{
		"stock_total": totalStock,
		"stocked_at":  now,
		"updated_at":  now,
		"updated_by":  ctx.GetUserHeader().ID,
	}).Error; err != nil {
		return err
	}

	return nil
}

func (r *repository) UpdateStockItem(ctx common.ServiceContextManager, itemID, totalStockItem int) error {
	span, context := apm.StartSpan(ctx.GetContext(), "repository.warehouse.UpdateStockItem", "repository")
	defer span.End()

	ctx.SetContext(context)

	var (
		now = time.Now()
	)

	// Get db transaction if exist
	tx := ctx.GetDB()
	if tx == nil {
		tx = r.db
	}

	tx = tx.WithContext(ctx.GetContext()).Where("id = ?", itemID)
	if err := tx.WithContext(ctx.GetContext()).Model(WarehouseStockItem{}).Updates(map[string]interface{}{
		"stock_total": totalStockItem,
		"stocked_at":  now,
		"updated_at":  now,
		"updated_by":  ctx.GetUserHeader().ID,
	}).Error; err != nil {
		return err
	}

	return nil
}

func (r *repository) CreateStockItem(ctx common.ServiceContextManager, request domain.WarehouseStockItem) (*int, error) {
	span, context := apm.StartSpan(ctx.GetContext(), "repository.warehouse.CreateStockItem", "repository")
	defer span.End()

	ctx.SetContext(context)

	tx := ctx.GetDB()
	if tx == nil {
		tx = r.db
	}

	stockItem := WarehouseStockItem{
		ShopID:      request.ShopID,
		WarehouseID: request.WarehouseID,
		Code:        request.Code,
		Name:        request.Name,
		Price:       request.Price,
		StockTotal:  request.StockTotal,
		StockedAt:   request.StockedAt,
	}

	if err := tx.WithContext(ctx.GetContext()).Create(&stockItem).Error; err != nil {
		return nil, err
	}
	return &stockItem.ID, nil
}

func (r *repository) CreateStockItemTransfer(ctx common.ServiceContextManager, request domain.WarehouseStockTransfer) (*int, error) {
	span, context := apm.StartSpan(ctx.GetContext(), "repository.warehouse.CreateStockItem", "repository")
	defer span.End()

	ctx.SetContext(context)

	tx := ctx.GetDB()
	if tx == nil {
		tx = r.db
	}

	stockItemTransfer := WarehouseStockTransfer{
		ID:                     request.ID,
		ShopID:                 request.ShopID,
		WarehouseOriginID:      request.WarehouseOriginID,
		WarehouseDestinationID: request.WarehouseDestinationID,
		WarehouseStockItemID:   request.WarehouseStockItemID,
		TransferTotal:          request.TransferTotal,
		TransferAt:             request.TransferAt,
	}

	if err := tx.WithContext(ctx.GetContext()).Create(&stockItemTransfer).Error; err != nil {
		return nil, err
	}

	return &stockItemTransfer.ID, nil
}

func (r *repository) UnactiveWarehouse(ctx common.ServiceContextManager, id int) error {
	span, context := apm.StartSpan(ctx.GetContext(), "repository.warehouse.UnactiveWarehouse", "repository")
	defer span.End()

	ctx.SetContext(context)

	var (
		now = time.Now()
	)

	return r.db.WithContext(ctx.GetContext()).Where("id = ?", id).Model(Warehouse{}).Updates(map[string]interface{}{
		"is_active":  false,
		"updated_at": now,
		"updated_by": ctx.GetUserHeader().ID,
	}).Error
}
