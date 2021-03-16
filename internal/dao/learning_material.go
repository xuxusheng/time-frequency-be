package dao

import (
	"context"
	"github.com/go-pg/pg/v10/orm"
	"github.com/xuxusheng/time-frequency-be/internal/model"
	"time"
)

type ILearningMaterial interface {
	Create(ctx context.Context, createdBy, subjectId int, name, description, md5, filePath string) (*model.LearningMaterial, error)
	Get(ctx context.Context, id int) (*model.LearningMaterial, error)
	Update(ctx context.Context, id, updatedBy int, name, description string) (*model.LearningMaterial, error)
	Delete(ctx context.Context, id int) error
}

func NewLearningMaterialDao(db orm.DB) *LearningMaterialDao {
	return &LearningMaterialDao{db: db}
}

type LearningMaterialDao struct {
	db orm.DB
}

func (l LearningMaterialDao) Create(ctx context.Context, createdById, subjectId int, name, description, md5, filePath string) (*model.LearningMaterial, error) {
	lm := model.LearningMaterial{
		Name:        name,
		Description: description,
		Md5:         md5,
		FilePath:    filePath,
		SubjectId:   subjectId,
		CreatedById: createdById,
		UpdatedById: createdById,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	_, err := l.db.ModelContext(ctx, &lm).Returning("*").Insert()
	if err != nil {
		return nil, err
	}
	return &lm, err
}

func (l LearningMaterialDao) Get(ctx context.Context, id int) (*model.LearningMaterial, error) {
	lm := model.LearningMaterial{Id: id}
	err := l.db.ModelContext(ctx, &lm).WherePK().Select()
	if err != nil {
		return nil, err
	}
	return &lm, err
}

func (l LearningMaterialDao) Update(ctx context.Context, id, updatedBy int, name, description string) (*model.LearningMaterial, error) {
	lm := model.LearningMaterial{
		Id:          id,
		Name:        name,
		Description: description,
		UpdatedById: updatedBy,
		UpdatedAt:   time.Now(),
	}
	_, err := l.db.ModelContext(ctx, &lm).Column("name", "description", "updated_by_id", "updated_at").WherePK().Returning("*").Update()
	if err != nil {
		return nil, err
	}
	return &lm, err
}

func (l LearningMaterialDao) Delete(ctx context.Context, id int) error {
	_, err := l.db.ModelContext(ctx, &model.LearningMaterial{Id: id}).WherePK().Delete()
	return err
}
