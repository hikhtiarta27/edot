package registry

import (
	"shop/usecase"
	"sync"
)

var (
	shopUsecaseOnce sync.Once
	shopUsecase     usecase.Shop
)

func LoadShopUsecase() usecase.Shop {
	shopUsecaseOnce.Do(func() {
		shopUsecase = usecase.NewShop(
			LoadShopRepo(),
			LoadWarehouseRepo(),
		)
	})

	return shopUsecase
}
