package service

import (
	"time"
	"warehousesvc/internal/core/domain"
	"warehousesvc/internal/core/port"
	"warehousesvc/pkg/cache"
	"warehousesvc/pkg/common"
	"warehousesvc/pkg/database"

	"go.elastic.co/apm"
)

type service struct {
	repository port.Repository
	cache      cache.CacheManager
}

func New(
	repository port.Repository,
	cache cache.CacheManager,
) port.Service {
	return &service{
		repository: repository,
		cache:      cache,
	}
}

func (s *service) FindAllStock(ctx common.ServiceContextManager, filter domain.FindAllStockRequest) ([]domain.FindAllStockResponse, error) {
	span, context := apm.StartSpan(ctx.GetContext(), "repository.warehouse.FindAllStock", "repository")
	defer span.End()

	ctx.SetContext(context)

	return s.repository.FindAllStock(ctx, filter)
}

func (s *service) TransferStock(ctx common.ServiceContextManager, request domain.TransferStockRequest) error {
	span, context := apm.StartSpan(ctx.GetContext(), "repository.warehouse.TransferStock", "repository")
	defer span.End()

	ctx.SetContext(context)

	var (
		now = time.Now()
	)

	tx := ctx.BeginTransaction(database.DB)
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	originStockItem, err := s.repository.FindStockItemByID(ctx, request.WarehouseOriginStockItemID, request.WarehouseOriginID)
	if err != nil {
		tx.Rollback()
		return err
	}

	if originStockItem == nil {
		tx.Rollback()
		return domain.ErrItemNotFound
	}

	if originStockItem.StockTotal == 0 || request.Number > originStockItem.StockTotal {
		tx.Rollback()
		return domain.ErrInsufficientItem
	}

	if err := s.repository.UpdateStock(ctx,
		request.WarehouseOriginID,
		request.FindAndReducePreviousStockTotal(originStockItem.Warehouse.StockTotal),
	); err != nil {
		tx.Rollback()
		return err
	}

	if err := s.repository.UpdateStockItem(ctx,
		request.WarehouseOriginStockItemID,
		request.FindAndReducePreviousStockTotal(originStockItem.StockTotal)); err != nil {
		tx.Rollback()
		return err
	}

	destinaitonStock, err := s.repository.FindStockItemByID(ctx, request.WarehouseOriginStockItemID, request.WarehouseDestinationID)
	if err != nil {
		tx.Rollback()
		return err
	}

	var SID int
	if destinaitonStock != nil {
		if err := s.repository.UpdateStock(ctx,
			request.WarehouseDestinationID,
			request.FindAndAddPreviousStockTotal(destinaitonStock.Warehouse.StockTotal),
		); err != nil {
			tx.Rollback()
			return err
		}

		if err := s.repository.UpdateStockItem(ctx,
			request.WarehouseOriginStockItemID,
			request.FindAndAddPreviousStockTotal(destinaitonStock.StockTotal)); err != nil {
			tx.Rollback()
			return err
		}
		SID = destinaitonStock.ID
	} else {
		stock, err := s.repository.FindStockByID(ctx, request.WarehouseDestinationID)
		if err != nil {
			tx.Rollback()
			return err
		}

		if err := s.repository.UpdateStock(ctx, stock.ID, stock.StockTotal+request.Number); err != nil {
			tx.Rollback()
			return err
		}

		stockItemID, err := s.repository.CreateStockItem(ctx, domain.WarehouseStockItem{
			ShopID:      originStockItem.ShopID,
			WarehouseID: request.WarehouseDestinationID,
			Code:        request.WarehouseOriginStockItemCode,
			Name:        originStockItem.Name,
			Price:       originStockItem.Price,
			StockTotal:  request.Number,
			StockedAt:   &now,
		})
		if err != nil {
			tx.Rollback()
			return err
		}
		SID = *stockItemID
	}

	_, err = s.repository.CreateStockItemTransfer(ctx, domain.WarehouseStockTransfer{
		ShopID:                 originStockItem.ShopID,
		WarehouseOriginID:      request.WarehouseOriginID,
		WarehouseDestinationID: request.WarehouseDestinationID,
		WarehouseStockItemID:   SID,
		TransferTotal:          request.Number,
		TransferAt:             &now,
	})
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (s *service) UnactiveWarehouse(ctx common.ServiceContextManager, id int) error {
	span, context := apm.StartSpan(ctx.GetContext(), "repository.warehouse.UnactiveWarehouse", "repository")
	defer span.End()

	ctx.SetContext(context)

	warehouse, err := s.repository.FindStockByID(ctx, id)
	if err != nil {
		return err
	}

	if warehouse == nil {
		return domain.ErrWarehouseNotFound
	}

	if !warehouse.IsActive {
		return domain.ErrWarehouseIsUnactive
	}

	return s.repository.UnactiveWarehouse(ctx, id)
}
