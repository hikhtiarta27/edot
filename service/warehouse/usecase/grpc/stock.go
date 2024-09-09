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
	stockRepo repository.Stock
}

func NewStock(
	stockRepo repository.Stock,
) *StockGrpc {
	return &StockGrpc{
		stockRepo: stockRepo,
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
