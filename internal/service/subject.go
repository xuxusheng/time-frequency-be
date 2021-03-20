package service

import (
	"context"
	"github.com/pkg/errors"
	"github.com/xuxusheng/time-frequency-be/internal/dao"
	"github.com/xuxusheng/time-frequency-be/internal/model"
	"github.com/xuxusheng/time-frequency-be/internal/pkg/cerror"
)

type ISubject interface {
	Create(ctx context.Context, createById int, name, description string) (*model.Subject, error)
	Get(ctx context.Context, id int) (*model.Subject, error)
	Update(ctx context.Context, id int, name, description string) (*model.Subject, error)
	Delete(ctx context.Context, id int) error
	IsNameExist(ctx context.Context, name string, excludeId int) (bool, error)
}

func NewSubject(dao dao.ISubject) *Subject {
	return &Subject{Dao: dao}
}

type Subject struct {
	Dao dao.ISubject
}

func (s Subject) Create(ctx context.Context, createById int, name, description string) (*model.Subject, error) {
	d := s.Dao
	// 科目名称是否重复
	is, err := d.IsNameExist(ctx, name, 0)
	if err != nil {
		return nil, err
	}
	if is {
		return nil, cerror.BadRequest.WithMsg("科目名称已存在")
	}
	subject, err := d.Create(ctx, createById, name, description)
	if err != nil {
		return nil, err
	}
	return subject, nil
}

func (s Subject) Get(ctx context.Context, id int) (*model.Subject, error) {
	return s.Dao.Get(ctx, id)
}

func (s Subject) Update(ctx context.Context, id int, name, description string) (*model.Subject, error) {
	d := s.Dao
	// 判断科目名称是否重复
	is, err := d.IsNameExist(ctx, name, id)
	if err != nil {
		return nil, err
	}
	if is {
		return nil, errors.New("科目名称已存在")
	}
	subject, err := d.Update(ctx, id, name, description)
	if err != nil {
		return nil, err
	}
	return subject, nil
}

func (s Subject) Delete(ctx context.Context, id int) error {
	return s.Dao.Delete(ctx, id)
}

func (s Subject) IsNameExist(ctx context.Context, name string, excludeId int) (bool, error) {
	return s.Dao.IsNameExist(ctx, name, excludeId)
}
