package dao

import (
	"context"
	"github.com/go-pg/pg/v10/orm"
	"github.com/xuxusheng/time-frequency-be/internal/model"

	"time"
)

type ISubjectDao interface {
	Create(ctx context.Context, createdById int, name, description string) (*model.Subject, error)
	Get(ctx context.Context, id int) (*model.Subject, error)
	Update(ctx context.Context, id int, name, description string) (*model.Subject, error)
	Delete(ctx context.Context, id int) error
	IsNameExist(ctx context.Context, name string, excludeId int) (bool, error)
}

func NewSubjectDao(db orm.DB) *SubjectDao {
	return &SubjectDao{db: db}
}

type SubjectDao struct {
	db orm.DB
}

func (s SubjectDao) Create(ctx context.Context, createdById int, name, description string) (*model.Subject, error) {
	subject := model.Subject{
		Name:        name,
		Description: description,
		CreatedById: createdById,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	_, err := s.db.ModelContext(ctx, &subject).Returning("*").Insert()
	if err != nil {
		return nil, err
	}
	return &subject, nil
}

func (s SubjectDao) Get(ctx context.Context, id int) (*model.Subject, error) {
	subject := model.Subject{Id: id}
	err := s.db.ModelContext(ctx, &subject).WherePK().Select()
	if err != nil {
		return nil, err
	}
	return &subject, err
}

func (s SubjectDao) Update(ctx context.Context, id int, name, description string) (*model.Subject, error) {
	subject := model.Subject{
		Id: id, Name: name, Description: description, UpdatedAt: time.Now(),
	}
	_, err := s.db.ModelContext(ctx, &subject).Column("name", "description", "updated_at").WherePK().Returning("*").Update()
	if err != nil {
		return nil, err
	}
	return &subject, err
}

func (s SubjectDao) Delete(ctx context.Context, id int) error {
	_, err := s.db.ModelContext(ctx, &model.Subject{Id: id}).WherePK().Delete()
	return err
}

func (s SubjectDao) IsNameExist(ctx context.Context, name string, excludeId int) (bool, error) {
	db := s.db.ModelContext(ctx, &model.Subject{})
	if excludeId != 0 {
		db = db.Where("id != ?", excludeId)
	}
	return db.Where("name = ?", name).Exists()
}
