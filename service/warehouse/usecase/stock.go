package usecase

import (
	"context"
	"warehouse/model"
	"warehouse/repository"
	"warehouse/usecase/stock"
)

type Stock interface {
	Create(ctx context.Context, param *stock.CreateRequest) (*stock.Stock, error)
}

type stockUsecase struct {
	stockRepo     repository.Stock
	productRepo   repository.Product
	warehouseRepo repository.Warehouse
}

func NewStock(
	stockRepo repository.Stock,
	productRepo repository.Product,
	warehouseRepo repository.Warehouse,
) Stock {
	return &stockUsecase{
		stockRepo:     stockRepo,
		productRepo:   productRepo,
		warehouseRepo: warehouseRepo,
	}
}

func (s stockUsecase) Create(ctx context.Context, param *stock.CreateRequest) (*stock.Stock, error) {

	if err := param.Validate(); err != nil {
		return nil, err
	}

	_, err := s.productRepo.Get(ctx, &model.GetProduct{
		ID: param.ProductID,
	})
	if err != nil {
		return nil, err
	}

	warehouse, err := s.warehouseRepo.Get(ctx, &model.GetWarehouse{
		ID:     param.WarehouseID,
		Status: model.WarehouseActive,
	})
	if err != nil {
		return nil, err
	}

	if warehouse == nil {
		return nil, model.ErrWarehouseNotFound
	}

	stockModel, err := model.NewStock(param.ProductID, param.Stock)
	if err != nil {
		return nil, err
	}

	warehouseTransfer, err := model.NewWarehouseTransfer(param.WarehouseID, param.WarehouseID, param.ProductID, param.Stock)
	if err != nil {
		return nil, err
	}

	err = s.stockRepo.Create(ctx, &model.CreateStock{
		Stock:             stockModel,
		WarehouseTransfer: warehouseTransfer,
	})
	if err != nil {
		return nil, err
	}

	return &stock.Stock{
		ID:             stockModel.ID,
		ProductID:      stockModel.ProductID,
		AvailableStock: stockModel.AvailableStock,
		ReservedStock:  stockModel.ReservedStock,
	}, nil
}
