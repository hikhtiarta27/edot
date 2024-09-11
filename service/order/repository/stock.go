package repository

import (
	"context"
	"order/model"
	"proto_buffer/stock"
	"time"

	"github.com/oklog/ulid/v2"
)

type Stock interface {
	ReserveRelease(ctx context.Context, param *model.ReserveReleaseStock) (*model.Stock, error)
}

type stockRepo struct {
	stockGrpc stock.StockServiceClient
}

func NewStock(
	stockGrpc stock.StockServiceClient,
) Stock {
	return &stockRepo{
		stockGrpc: stockGrpc,
	}
}

func (r stockRepo) ReserveRelease(ctx context.Context, param *model.ReserveReleaseStock) (*model.Stock, error) {

	res, err := r.stockGrpc.ReserveRelease(ctx, &stock.ReserveReleaseRequest{
		ProductId: param.ProductID.String(),
		Action:    string(param.Action),
		Qty:       param.Qty,
	})

	if err != nil {
		return nil, err
	}

	return &model.Stock{
		ID:             ulid.MustParse(res.Id),
		ProductID:      ulid.MustParse(res.ProductId),
		AvailableStock: res.AvailableStock,
		ReservedStock:  res.ReservedStock,
		CreatedAt:      time.Unix(res.CreatedAt, 0),
	}, nil
}
