package service

import (
	"context"
	"github.com/pkg/errors"
	"github.com/xuxusheng/time-frequency-be/internal/dao"
	"github.com/xuxusheng/time-frequency-be/internal/model"
)

type ILearningMaterialSvc interface {
	Create(ctx context.Context, createdById, subjectId int, name, description, md5, filePath string) (*model.LearningMaterial, error)
	Get(ctx context.Context, id int) (*model.LearningMaterial, error)
	Update(ctx context.Context, id, updatedById int, name, description string) (*model.LearningMaterial, error)
	Delete(ctx context.Context, id int) error
	IsNameExist(ctx context.Context, name string, excludeId int) (bool, error)
}

func NewLearningMaterialSvc(dao dao.LearningMaterialDao) *LearningMaterialSvc {
	return &LearningMaterialSvc{Dao: dao}
}

type LearningMaterialSvc struct {
	Dao dao.LearningMaterialDao
}

func (l LearningMaterialSvc) Create(ctx context.Context, createdById, subjectId int, name, description, md5, filePath string) (*model.LearningMaterial, error) {
	d := l.Dao
	// 判断资料名称是否存在
	is, err := d.IsNameExist(ctx, name, 0)
	if err != nil {
		return nil, err
	}
	if is {
		return nil, errors.New("资料名称已存在")
	}
	lm, err := d.Create(ctx, createdById, subjectId, name, description, md5, filePath)
	if err != nil {
		return nil, err
	}
	return lm, nil
}

func (l LearningMaterialSvc) Get(ctx context.Context, id int) (*model.LearningMaterial, error) {
	return l.Dao.Get(ctx, id)
}

func (l LearningMaterialSvc) Update(ctx context.Context, id, updatedById int, name, description string) (*model.LearningMaterial, error) {
	d := l.Dao
	// 判断资料名称是否存在
	is, err := d.IsNameExist(ctx, name, id)
	if err != nil {
		return nil, err
	}
	if is {
		return nil, errors.New("资料名称已存在")
	}
	lm, err := d.Update(ctx, id, updatedById, name, description)
	if err != nil {
		return nil, err
	}
	return lm, nil
}

func (l LearningMaterialSvc) Delete(ctx context.Context, id int) error {
	return l.Dao.Delete(ctx, id)
}

func (l LearningMaterialSvc) IsNameExist(ctx context.Context, name string, excludeId int) (bool, error) {
	return l.Dao.IsNameExist(ctx, name, excludeId)
}
