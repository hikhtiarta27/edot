package usecase

import (
	"context"
	"order/model"
	"order/repository"
	"order/usecase/order"
	"time"
)

type Order interface {
	Create(ctx context.Context, param *order.CreateRequest) (*order.Order, error)
}

type orderUsecase struct {
	orderRepo   repository.Order
	productRepo repository.Product
}

func NewOrder(
	orderRepo repository.Order,
	productRepo repository.Product,
) Order {
	return &orderUsecase{
		orderRepo:   orderRepo,
		productRepo: productRepo,
	}
}

func (s *orderUsecase) Create(ctx context.Context, param *order.CreateRequest) (*order.Order, error) {

	if err := param.Validate(); err != nil {
		return nil, err
	}

	newOrder, err := model.NewOrder(uint64(len(param.Product)))
	if err != nil {
		return nil, err
	}

	for _, prd := range param.Product {
		productModel, err := s.productRepo.Get(ctx, &model.GetProduct{
			ID: prd.ProductID,
		})

		if err != nil {
			return nil, err
		}

		newOrder.AddDetail(productModel, prd.Qty)
	}

	newOrder.ExpiredAt = time.Now().Add(time.Second * 10)

	err = s.orderRepo.Create(ctx, newOrder)
	if err != nil {
		return nil, err
	}

	res := &order.Order{
		ID:         newOrder.ID,
		TotalItem:  newOrder.TotalItem,
		TotalPrice: newOrder.TotalPrice,
		Status:     newOrder.Status,
		ExpiredAt:  newOrder.ExpiredAt,
		CreatedAt:  newOrder.CreatedAt,
	}

	for _, prd := range newOrder.Detail {
		res.Detail = append(res.Detail, order.OrderDetail{
			ID:          prd.ID,
			ProductID:   prd.ProductID,
			ProductName: prd.ProductName,
			Qty:         prd.Qty,
			Price:       prd.Price,
			TotalPrice:  prd.GetTotalPrice(),
		})
	}

	return res, nil
}
