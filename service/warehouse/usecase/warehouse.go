package usecase

import (
	"context"
	"strings"
	"warehouse/model"
	"warehouse/repository"
	"warehouse/usecase/warehouse"

	"github.com/oklog/ulid/v2"
)

type Warehouse interface {
	Create(ctx context.Context, param *warehouse.CreateRequest) (*warehouse.Warehouse, error)
	Update(ctx context.Context, param *warehouse.UpdateRequest) (*warehouse.Warehouse, error)
	List(ctx context.Context, param *warehouse.ListRequest) ([]warehouse.Warehouse, error)
	TransferStock(ctx context.Context, param *warehouse.TransferRequest) (*warehouse.WarehouseTransfer, error)
}

type warehouseUsecase struct {
	warehouseRepo         repository.Warehouse
	productRepo           repository.Product
	warehouseTransferRepo repository.WarehouseTransfer
}

func NewWarehouse(
	warehouseRepo repository.Warehouse,
	productRepo repository.Product,
	warehouseTransferRepo repository.WarehouseTransfer,
) Warehouse {
	return &warehouseUsecase{
		warehouseRepo:         warehouseRepo,
		productRepo:           productRepo,
		warehouseTransferRepo: warehouseTransferRepo,
	}
}

func (s warehouseUsecase) Create(ctx context.Context, param *warehouse.CreateRequest) (*warehouse.Warehouse, error) {

	if err := param.Validate(); err != nil {
		return nil, err
	}

	warehouseModel, err := model.NewWarehouse(param.Name)
	if err != nil {
		return nil, err
	}

	err = s.warehouseRepo.Create(ctx, warehouseModel)
	if err != nil {
		return nil, err
	}

	return &warehouse.Warehouse{
		ID:        warehouseModel.ID,
		Name:      warehouseModel.Name,
		Status:    warehouseModel.Status,
		CreatedAt: warehouseModel.CreatedAt,
	}, nil
}

func (s warehouseUsecase) Update(ctx context.Context, param *warehouse.UpdateRequest) (*warehouse.Warehouse, error) {

	if err := param.Validate(); err != nil {
		return nil, err
	}

	id, err := ulid.Parse(param.ID)
	if err != nil {
		return nil, model.ErrInvalidUlid
	}

	existingWarehouse, err := s.warehouseRepo.Get(ctx, &model.GetWarehouse{
		ID: id,
	})
	if err != nil {
		return nil, err
	}

	if existingWarehouse == nil {
		return nil, model.ErrWarehouseNotFound
	}

	existingWarehouse.Name = strings.TrimSpace(param.Name)

	existingWarehouse.Status = param.Status

	err = s.warehouseRepo.Update(ctx, existingWarehouse)
	if err != nil {
		return nil, err
	}

	return &warehouse.Warehouse{
		ID:        existingWarehouse.ID,
		Name:      existingWarehouse.Name,
		Status:    existingWarehouse.Status,
		CreatedAt: existingWarehouse.CreatedAt,
		UpdatedAt: existingWarehouse.UpdatedAt,
	}, nil
}

func (s warehouseUsecase) List(ctx context.Context, param *warehouse.ListRequest) ([]warehouse.Warehouse, error) {

	warehouses, err := s.warehouseRepo.Select(ctx, &model.SelectWarehouse{})
	if err != nil {
		return nil, err
	}

	var res = make([]warehouse.Warehouse, 0)

	for _, whs := range warehouses {
		res = append(res, warehouse.Warehouse{
			ID:        whs.ID,
			Name:      whs.Name,
			Status:    whs.Status,
			CreatedAt: whs.CreatedAt,
			UpdatedAt: whs.UpdatedAt,
		})
	}

	return res, nil
}

func (s warehouseUsecase) TransferStock(ctx context.Context, param *warehouse.TransferRequest) (*warehouse.WarehouseTransfer, error) {

	if err := param.Validate(); err != nil {
		return nil, err
	}

	fromWarehouse, err := s.warehouseRepo.Get(ctx, &model.GetWarehouse{
		ID: param.FromWarehouseID,
	})
	if err != nil {
		return nil, err
	}

	if fromWarehouse == nil {
		return nil, model.ErrWarehouseNotFound
	}

	if fromWarehouse.Status != model.WarehouseActive {
		return nil, model.ErrWarehouseInactive
	}

	toWarehouse, err := s.warehouseRepo.Get(ctx, &model.GetWarehouse{
		ID: param.ToWarehouseID,
	})
	if err != nil {
		return nil, err
	}

	if toWarehouse == nil {
		return nil, model.ErrWarehouseNotFound
	}

	if toWarehouse.Status != model.WarehouseActive {
		return nil, model.ErrWarehouseInactive
	}

	// initialize stock
	// if fromWarehouse.ID == toWarehouse.ID {
	// 	return nil, &shared.Error{
	// 		HttpStatusCode: 422,
	// 		Message:        "failed to transfer stock. same warehouse id",
	// 	}
	// }

	product, err := s.productRepo.Get(ctx, &model.GetProduct{
		ID: param.ProductID,
	})
	if err != nil {
		return nil, err
	}

	warehouseTransfer, err := model.NewWarehouseTransfer(fromWarehouse.ID, toWarehouse.ID, product.ID, param.Stock)
	if err != nil {
		return nil, err
	}

	err = s.warehouseTransferRepo.Create(ctx, warehouseTransfer)
	if err != nil {
		return nil, err
	}

	return (*warehouse.WarehouseTransfer)(warehouseTransfer), nil
}
