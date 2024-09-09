package repository

import (
	"context"
	"product/model"
	"proto_buffer/stock"
	"time"

	"github.com/oklog/ulid/v2"
)

type Stock interface {
	Select(ctx context.Context, param *model.SelectStock) (model.Stocks, error)
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

func (r stockRepo) Select(ctx context.Context, param *model.SelectStock) (model.Stocks, error) {

	var productIDs []string

	for _, productId := range param.ProductIDs {
		productIDs = append(productIDs, productId.String())
	}

	res, err := r.stockGrpc.Get(ctx, &stock.GetRequest{
		ProductId: productIDs,
	})

	if err != nil {
		return nil, err
	}

	var stocks model.Stocks

	for _, stk := range res.Stock {
		stocks = append(stocks, model.Stock{
			ID:             ulid.MustParse(stk.Id),
			ProductID:      ulid.MustParse(stk.ProductId),
			AvailableStock: stk.AvailableStock,
			ReservedStock:  stk.ReservedStock,
			CreatedAt:      time.Unix(stk.CreatedAt, 0),
		})
	}

	return stocks, nil
}
