package registry

import (
	"order/infra"
	"order/repository"
	"proto_buffer/product"
	"proto_buffer/stock"
	"sync"
)

var (
	orderRepoOnce sync.Once
	orderRepo     repository.Order
)

func LoadOrderRepo() repository.Order {
	orderRepoOnce.Do(func() {
		orderRepo = repository.NewOrder(
			infra.LoadDB(),
			LoadStockRepo(),
		)
	})

	return orderRepo
}

var (
	orderDetailRepoOnce sync.Once
	orderDetailRepo     repository.OrderDetail
)

func LoadOrderDetailRepo() repository.OrderDetail {
	orderDetailRepoOnce.Do(func() {
		orderDetailRepo = repository.NewOrderDetail(
			infra.LoadDB(),
		)
	})

	return orderDetailRepo
}

var (
	productRepoOnce sync.Once
	productRepo     repository.Product
)

func LoadProductRepo() repository.Product {
	productRepoOnce.Do(func() {
		productRepo = repository.NewProduct(
			product.NewProductServiceClient(
				infra.LoadProductService(),
			),
		)
	})

	return productRepo
}

var (
	stockRepoOnce sync.Once
	stockRepo     repository.Stock
)

func LoadStockRepo() repository.Stock {
	stockRepoOnce.Do(func() {
		stockRepo = repository.NewStock(
			stock.NewStockServiceClient(
				infra.LoadWarehouseService(),
			),
		)
	})

	return stockRepo
}
