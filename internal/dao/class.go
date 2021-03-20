package dao

import (
	"context"
	"github.com/go-pg/pg/v10/orm"
	"github.com/xuxusheng/time-frequency-be/internal/model"
	"time"
)

type IClass interface {
	Create(ctx context.Context, createdById int, name, description string) (*model.Class, error)
	Get(ctx context.Context, id int) (*model.Class, error)
	Update(ctx context.Context, id int, name, description string) (*model.Class, error)
	Delete(ctx context.Context, id int) error
	IsNameExist(ctx context.Context, name string, excludeId int) (bool, error)
}

func NewClass(db orm.DB) *Class {
	return &Class{db: db}
}

type Class struct {
	db orm.DB
}

func (c Class) Create(ctx context.Context, createdById int, name, description string) (*model.Class, error) {
	class := model.Class{
		Name:        name,
		Description: description,
		CreatedById: createdById,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	_, err := c.db.ModelContext(ctx, &class).Returning("*").Insert()
	if err != nil {
		return nil, err
	}
	return &class, err
}

func (c Class) Get(ctx context.Context, id int) (*model.Class, error) {
	class := model.Class{Id: id}
	err := c.db.ModelContext(ctx, &class).WherePK().Select()
	if err != nil {
		return nil, err
	}
	return &class, nil
}

func (c Class) Update(ctx context.Context, id int, name, description string) (*model.Class, error) {
	class := model.Class{Id: id, Name: name, Description: description, UpdatedAt: time.Now()}
	_, err := c.db.
		ModelContext(ctx, &class).
		Column("name", "description", "updated_at").
		WherePK().
		Returning("*").
		Update()
	if err != nil {
		return nil, err
	}
	return &class, nil
}

func (c Class) Delete(ctx context.Context, id int) error {
	_, err := c.db.ModelContext(ctx, &model.Class{Id: id}).WherePK().Delete()
	return err
}

func (c Class) IsNameExist(ctx context.Context, name string, excludeId int) (bool, error) {
	db := c.db.ModelContext(ctx, &model.Class{})
	if excludeId != 0 {
		db = db.Where("id != ?", excludeId)
	}
	return db.Where("name = ?", name).Exists()
}
