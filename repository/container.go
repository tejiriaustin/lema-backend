package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/tejiriaustin/lema/logger"
	"gorm.io/gorm"
	"log"
	"time"

	"github.com/tejiriaustin/lema/database"
	"github.com/tejiriaustin/lema/models"
)

type (
	Container struct {
		UserRepo    *Repository[models.User]
		PostRepo    *Repository[models.Post]
		AddressRepo *Repository[models.Address]
	}
	Repository[T models.Models] struct {
		db *gorm.DB
	}
)

func NewRepositoryContainer(lemaLogger logger.Logger, dbConn *database.Client) *Container {
	log.Println("building repository container...")

	return &Container{
		UserRepo:    NewRepository[models.User](dbConn.GetModel("users")),
		PostRepo:    NewRepository[models.Post](dbConn.GetModel("posts")),
		AddressRepo: NewRepository[models.Address](dbConn.GetModel("addresses")),
	}
}

func NewRepository[T models.Models](client database.Client) *Repository[T] {
	return &Repository[T]{db: client.DB}
}

var _ RepoInterface[models.Shared] = (*Repository[models.Shared])(nil)

func (r *Repository[T]) Create(ctx context.Context, data T) (*T, error) {
	data.SetVersion(1) // Initialize version to 1 for new records
	data.Initialize(uuid.New(), time.Now())
	fmt.Println(data)
	result := r.db.WithContext(ctx).Create(&data)
	if result.Error != nil {
		return &data, result.Error
	}
	return &data, nil
}

func (r *Repository[T]) FindOne(ctx context.Context, queryFilter *Query) (*T, error) {
	var result *T
	db := r.db.WithContext(ctx)

	if queryFilter != nil {
		db = db.Where(queryFilter.query, queryFilter.args...)
	}

	if err := db.First(&result).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return result, ErrNotFound
		}
		return result, err
	}
	return result, nil
}

func (r *Repository[T]) FindManyPaginated(ctx context.Context, queryFilter *Query, page, perPage int64) ([]*T, *Paginator, error) {
	paginator := newPaginator(page, perPage)
	paginator.setOffset()

	var total int64
	db := r.db.WithContext(ctx).Model(new(T))

	if queryFilter != nil {
		db = db.Where(queryFilter.query, queryFilter.args...)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, nil, fmt.Errorf("count failed: %w", err)
	}

	paginator.TotalRows = total
	paginator.setTotalPages()
	paginator.setPrevPage()
	paginator.setNextPage()

	var results []*T
	if err := db.Offset(int(paginator.Offset)).Limit(int(paginator.PerPage)).Find(&results).Error; err != nil {
		return nil, nil, fmt.Errorf("find failed: %w", err)
	}

	return results, paginator, nil
}

func (r *Repository[T]) DeleteMany(ctx context.Context, queryFilter *Query) error {
	db := r.db.WithContext(ctx)

	if queryFilter != nil {
		db = db.Where(queryFilter.query, queryFilter.args...)
	}

	var model *T
	if err := db.Delete(&model).Error; err != nil {
		return err
	}
	return nil
}

func (r *Repository[T]) Update(ctx context.Context, dataObject T) (*T, error) {
	currentVersion := dataObject.GetVersion()
	newVersion := currentVersion + 1
	dataObject.SetVersion(newVersion)

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&dataObject).
			Where("_version = ?", currentVersion).
			Updates(dataObject)

		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected == 0 {
			return ErrConcurrentModification
		}

		return nil
	})

	if err != nil {
		// Revert version change on error
		dataObject.SetVersion(currentVersion)
		return &dataObject, err
	}

	return &dataObject, nil
}

func (r *Repository[T]) Select(ctx context.Context, target interface{}, query string, args ...interface{}) error {
	return r.db.WithContext(ctx).Select(query, args...).Scan(target).Error
}

func (r *Repository[T]) Count(ctx context.Context, queryFilter *Query) (int64, error) {
	var count int64
	db := r.db.WithContext(ctx).Model(new(T))

	if queryFilter != nil {
		db = db.Where(queryFilter.query, queryFilter.args...)
	}

	if err := db.Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count records: %w", err)
	}

	return count, nil
}
