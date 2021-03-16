package service

import (
	"context"
	"github.com/pkg/errors"
	"github.com/xuxusheng/time-frequency-be/internal/dao"
	"github.com/xuxusheng/time-frequency-be/internal/model"
)

type IClassSvc interface {
	Create(ctx context.Context, createdById int, name, description string) (*model.Class, error)
	Get(ctx context.Context, id int) (*model.User, error)
	Update(ctx context.Context, id, updatedBy int, name, description string) (*model.Class, error)
	Delete(ctx context.Context, id int) error
	IsNameExist(ctx context.Context, name string, excludeId int) (bool, error)
}

func NewClassSvc(dao dao.IClassDao) *ClassSvc {
	return &ClassSvc{Dao: dao}
}

type ClassSvc struct {
	Dao dao.IClassDao
}

func (c ClassSvc) Create(ctx context.Context, createdById int, name, description string) (*model.Class, error) {
	d := c.Dao

	// 判断班级名称是否重复
	is, err := d.IsNameExist(ctx, name, 0)
	if err != nil {
		return nil, err
	}
	if is {
		return nil, errors.New("班级名称已存在")
	}

	class, err := d.Create(ctx, createdById, name, description)
	if err != nil {
		return nil, err
	}
	return class, nil
}

func (c ClassSvc) Get(ctx context.Context, id int) (*model.User, error) {
	return c.Get(ctx, id)
}

func (c ClassSvc) Update(ctx context.Context, id, updatedBy int, name, description string) (*model.Class, error) {
	d := c.Dao
	// 判断班级名称是否已被占用
	is, err := d.IsNameExist(ctx, name, id)
	if err != nil {
		return nil, err
	}
	if is {
		return nil, errors.New("班级名称已存在")
	}
	class, err := d.Update(ctx, id, name, description)
	if err != nil {
		return nil, err
	}
	return class, nil
}

func (c ClassSvc) Delete(ctx context.Context, id int) error {
	return c.Dao.Delete(ctx, id)
}

func (c ClassSvc) IsNameExist(ctx context.Context, name string, excludeId int) (bool, error) {
	return c.Dao.IsNameExist(ctx, name, excludeId)
}
