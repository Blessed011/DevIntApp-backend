package repository

import (
	"lab1/internal/app/ds"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func New(dsn string) (*Repository, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &Repository{
		db: db,
	}, nil
}

func (r *Repository) GetAllModules() ([]ds.Module, error) {
	var modules []ds.Module

	err := r.db.Find(&modules).Error

	if err != nil {
		return nil, err
	}

	return modules, nil
}

func (r *Repository) GetModuleByID(id string) (*ds.Module, error) {
	module := &ds.Module{}

	err := r.db.Where("module_id = ?", id).First(module).Error
	if err != nil {
		return nil, err
	}

	return module, nil
}

func (r *Repository) GetModuleByName(name string) ([]ds.Module, error) {
	var modules []ds.Module

	err := r.db.
		Where("LOWER(modules.name) LIKE ?", "%"+strings.ToLower(name)+"%").Where("is_deleted = ?", false).
		Find(&modules).Error

	if err != nil {
		return nil, err
	}

	return modules, nil
}

func (r *Repository) DeleteModule(id string) error {
	err := r.db.Exec("UPDATE modules SET is_deleted = ? WHERE module_id = ?", true, id).Error

	if err != nil {
		return err
	}

	return nil
}
