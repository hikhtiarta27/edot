package registry

import (
	"order/infra"
	"order/usecase"
	"sync"
)

var (
	orderUsecaseOnce sync.Once
	orderUsecase     usecase.Order
)

func LoadOrderUsecase() usecase.Order {
	orderUsecaseOnce.Do(func() {
		orderUsecase = usecase.NewOrder(
			LoadOrderRepo(),
			LoadProductRepo(),
			infra.LoadRedis(),
			LoadStockRepo(),
			LoadOrderDetailRepo(),
			LoadTxRepo(),
		)
	})

	return orderUsecase
}
