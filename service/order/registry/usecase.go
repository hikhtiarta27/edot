package registry

import (
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
		)
	})

	return orderUsecase
}
