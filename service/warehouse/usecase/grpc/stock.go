package grpc

import (
	"context"
	"proto_buffer/stock"
	"warehouse/model"
	"warehouse/repository"

	"github.com/oklog/ulid/v2"
)

type StockGrpc struct {
	stock.UnimplementedStockServiceServer
	stockRepo     repository.Stock
	productRepo   repository.Product
	warehouseRepo repository.Warehouse
}

func NewStock(
	stockRepo repository.Stock,
	productRepo repository.Product,
	warehouseRepo repository.Warehouse,
) *StockGrpc {
	return &StockGrpc{
		stockRepo:     stockRepo,
		productRepo:   productRepo,
		warehouseRepo: warehouseRepo,
	}
}

func (s StockGrpc) Get(ctx context.Context, param *stock.GetRequest) (*stock.GetResponse, error) {

	var productIDs []ulid.ULID

	for _, productID := range param.ProductId {
		id, err := ulid.Parse(productID)
		if err != nil {
			return nil, model.ErrInvalidUlid
		}

		productIDs = append(productIDs, id)
	}

	stocks, err := s.stockRepo.Select(ctx, &model.SelectStock{
		ProductIDs: productIDs,
	})

	if err != nil {
		return nil, err
	}

	var stockGrpcs []*stock.Stock

	for _, stk := range stocks {
		stockGrpcs = append(stockGrpcs, &stock.Stock{
			Id:             stk.ID.String(),
			ProductId:      stk.ProductID.String(),
			AvailableStock: stk.AvailableStock,
			ReservedStock:  stk.ReservedStock,
			CreatedAt:      stk.CreatedAt.Unix(),
		})
	}

	return &stock.GetResponse{Stock: stockGrpcs}, nil
}

func (s StockGrpc) Create(ctx context.Context, param *stock.CreateRequest) (*stock.Stock, error) {

	productID, err := ulid.Parse(param.ProductId)
	if err != nil {
		return nil, model.ErrInvalidUlid
	}

	_, err = s.productRepo.Get(ctx, &model.GetProduct{
		ID: productID,
	})
	if err != nil {
		return nil, err
	}

	warehouseID, err := ulid.Parse(param.WarehouseId)
	if err != nil {
		return nil, model.ErrInvalidUlid
	}

	warehouse, err := s.warehouseRepo.Get(ctx, &model.GetWarehouse{
		ID:     warehouseID,
		Status: model.WarehouseActive,
	})
	if err != nil {
		return nil, err
	}

	if warehouse == nil {
		return nil, model.ErrWarehouseNotFound
	}

	warehouseTransfer, err := model.NewWarehouseTransfer(warehouseID, warehouseID, productID, param.Stock)
	if err != nil {
		return nil, err
	}

	stockModel, err := s.stockRepo.Get(ctx, &model.GetStock{
		ProductID: productID,
	})
	if err != nil {
		return nil, err
	}

	if stockModel != nil {
		return nil, model.ErrDuplicateStock
	}

	stockModel, err = model.NewStock(productID, param.Stock)
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
		Id:             stockModel.ID.String(),
		ProductId:      stockModel.ProductID.String(),
		AvailableStock: stockModel.AvailableStock,
		ReservedStock:  stockModel.ReservedStock,
		CreatedAt:      stockModel.CreatedAt.Unix(),
	}, nil
}
