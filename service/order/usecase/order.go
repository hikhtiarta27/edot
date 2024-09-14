package usecase

import (
	"context"
	"fmt"
	"order/model"
	"order/repository"
	"order/usecase/order"
	"time"

	"log"

	"github.com/oklog/ulid/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Order interface {
	Create(ctx context.Context, param *order.CreateRequest) (*order.Order, error)
	Release(ctx context.Context) error
}

type orderUsecase struct {
	orderRepo       repository.Order
	productRepo     repository.Product
	redisClient     *redis.Client
	stockRepo       repository.Stock
	orderDetailRepo repository.OrderDetail
	txRepo          repository.Tx
}

func NewOrder(
	orderRepo repository.Order,
	productRepo repository.Product,
	redisClient *redis.Client,
	stockRepo repository.Stock,
	orderDetailRepo repository.OrderDetail,
	txRepo repository.Tx,
) Order {
	return &orderUsecase{
		orderRepo:       orderRepo,
		productRepo:     productRepo,
		redisClient:     redisClient,
		stockRepo:       stockRepo,
		orderDetailRepo: orderDetailRepo,
		txRepo:          txRepo,
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

	err = s.redisClient.ZAdd(ctx, "order_release", redis.Z{
		Score:  float64(newOrder.ExpiredAt.Unix()),
		Member: newOrder.ID.String(),
	}).Err()

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

func (s *orderUsecase) Release(ctx context.Context) error {
	member, err := s.redisClient.ZRangeByScore(ctx, "order_release", &redis.ZRangeBy{
		Min: "-inf",
		Max: fmt.Sprintf("%d", time.Now().Unix()),
	}).Result()

	if err != nil {
		log.Println(err)
		return err
	}

	for _, k := range member {

		go func() {

			orderID, err := ulid.Parse(k)
			if err != nil {
				return
			}

			err = s.txRepo.DoInTransaction(context.Background(), func(ctx context.Context, tx *gorm.DB) error {

				order, err := s.orderRepo.Get(ctx, &model.GetOrder{
					ID: orderID,
				})
				if err != nil {
					return err
				}

				orderDetails, err := s.orderDetailRepo.Select(ctx, &model.SelectOrderDetail{
					OrderID: orderID,
				})

				for _, orderDetail := range orderDetails {
					_, err = s.stockRepo.ReserveRelease(ctx, &model.ReserveReleaseStock{
						Action:    model.StockRelease,
						ProductID: orderDetail.ProductID,
						Qty:       orderDetail.Qty,
					})

					if err != nil {
						log.Printf("failed to release expired order detail %s. err: %v", orderDetail.ID, err)
					}
				}

				order.Status = model.OrderStatusExpired

				err = s.orderRepo.Update(ctx, tx, order)
				if err != nil {
					return err
				}

				return nil
			})

			if err != nil {
				log.Printf("failed to release expired order %s. err: %v", orderID, err)
			}
		}()
	}

	return nil
}
