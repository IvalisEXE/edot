package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"ordersvc/internal/core/domain"
	"ordersvc/internal/core/port"
	"ordersvc/pkg/cache"
	"ordersvc/pkg/common"
	"ordersvc/pkg/database"
	"time"

	"go.elastic.co/apm"
)

var (
	stockBookingUrl = "http://localhost:3004/warehouse/stock-booking"
	stockReleaseUrl = "http://localhost:3004/warehouse/stock-booking-release"
)

type service struct {
	repository port.Repository
	cache      cache.CacheManager
}

// New creates a new user service
func New(
	repository port.Repository,
	cache cache.CacheManager,
) port.Service {
	return &service{
		repository: repository,
		cache:      cache,
	}
}

func (s *service) Order(ctx common.ServiceContextManager, request domain.OrderRequest) error {
	span, context := apm.StartSpan(ctx.GetContext(), "repository.order.Order", "repository")
	defer span.End()

	ctx.SetContext(context)

	tx := ctx.BeginTransaction(database.DB)
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	// validasi warehouse_stock_item / product ada atau tidak -> warehousesvc call

	// validasi jumlah stok ada atau tidak

	// buat order
	if err := s.repository.Order(ctx, request); err != nil {
		return err
	}

	tx.Commit()

	return nil
}

func (s *service) LockStock(ctx common.ServiceContextManager) error {

	var (
		now = time.Now()
	)
	orders, err := s.repository.CheckOrder(ctx, "OPEN", now, now.Add(-time.Hour))
	if err != nil {
		log.Printf("[LockStock] Error %s\n", err.Error())
	}

	if orders == nil {
		log.Printf("[LockStock] Empty orders at  %s\n", now.String())
	}

	for _, order := range orders {
		s.stockAPI(stockBookingUrl, order)
	}

	return nil
}

func (s *service) ReleaseStock(ctx common.ServiceContextManager) error {

	var (
		now = time.Now()
	)
	orders, err := s.repository.CheckOrder(ctx, "CLOSED", now, now.Add(-time.Hour))
	if err != nil {
		return fmt.Errorf("[ReleaseStock] Error %s", err.Error())
	}

	if orders == nil {
		return fmt.Errorf("[ReleaseStock] Empty orders at  %s", now.String())
	}

	for _, order := range orders {
		s.stockAPI(stockReleaseUrl, order)
	}

	return nil
}

func (s *service) stockAPI(url string, order domain.Order) error {
	orderjson, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %s", err.Error())
	}

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(orderjson))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %s", err.Error())
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("failed to send request to warehouse service: %s", err.Error())
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %s", err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error response from warehouse service: %s", string(respBody))
	}

	return nil
}
