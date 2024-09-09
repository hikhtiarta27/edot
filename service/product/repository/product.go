package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"product/model"
	"shared"

	"github.com/meilisearch/meilisearch-go"
	"gorm.io/gorm"
)

type Product interface {
	Get(ctx context.Context, param *model.GetProduct) (*model.Product, error)
	Select(ctx context.Context, param *model.SelectProduct) ([]model.Product, error)
	Create(ctx context.Context, product *model.Product) error
	UpdateAvailableStock(ctx context.Context, param *model.UpdateStockProduct) error
}

type productRepo struct {
	db                *gorm.DB
	meilisearchClient meilisearch.ServiceManager
}

func NewProduct(
	db *gorm.DB,
	meilisearchClient meilisearch.ServiceManager,
) Product {
	return &productRepo{
		db:                db,
		meilisearchClient: meilisearchClient,
	}
}

func (r productRepo) Get(ctx context.Context, param *model.GetProduct) (*model.Product, error) {

	product := &model.Product{}

	q := r.db.
		WithContext(ctx)

	if !shared.IsZero(param.ID) {
		q = q.Where("id = ?", param.ID)
	}

	err := q.
		First(&product).
		Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return product, nil
}

func (r productRepo) Select(ctx context.Context, param *model.SelectProduct) ([]model.Product, error) {

	var products []model.Product

	if param.UseSearchEngine {

		res, err := r.meilisearchClient.
			Index("products").
			SearchWithContext(ctx, param.Keyword, &meilisearch.SearchRequest{})
		if err != nil {
			return nil, err
		}

		for _, hit := range res.Hits {
			b, err := json.Marshal(hit)
			if err != nil {
				return nil, err
			}

			var product model.Product

			err = json.Unmarshal(b, &product)
			if err != nil {
				return nil, err
			}

			products = append(products, product)
		}

		return products, nil
	}

	q := r.db.
		WithContext(ctx)

	err := q.
		Find(&products).
		Error

	if err != nil {
		return nil, err
	}

	return products, nil
}

func (r productRepo) Create(ctx context.Context, product *model.Product) error {

	tx := r.db.Begin(&sql.TxOptions{})

	err := tx.
		WithContext(ctx).
		Create(product).
		Error

	if err != nil {
		tx.Rollback()
		return err
	}

	// add to search engine
	_, err = r.meilisearchClient.Index("products").
		AddDocumentsWithContext(ctx, product)
	if err != nil {
		return err
	}

	// save redis stock

	tx.Commit()
	return nil
}

func (r productRepo) UpdateAvailableStock(ctx context.Context, param *model.UpdateStockProduct) error {
	tx := r.db.Begin(&sql.TxOptions{})

	err := tx.
		WithContext(ctx).
		Model(&model.Product{}).
		Where("id = ?", param.ID).
		Updates(map[string]interface{}{
			"available_stock": param.Stock,
		}).
		Error

	if err != nil {
		tx.Rollback()
		return err
	}

	// update to search engine
	_, err = r.meilisearchClient.Index("products").
		UpdateDocumentsWithContext(ctx, map[string]interface{}{
			"id":              param.ID,
			"available_stock": param.Stock,
		})
	if err != nil {
		return err
	}

	// save redis

	tx.Commit()
	return nil
}
