package repository

import (
	"errors"
	"strings"

	"lab1/internal/app/ds"

	"gorm.io/gorm"
)

func (r *Repository) GetModuleByID(id string) (*ds.Module, error) {
	module := &ds.Module{UUID: id}
	err := r.db.First(module, "is_deleted = ?", false).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return module, nil
}

func (r *Repository) AddModule(module *ds.Module) error {
	err := r.db.Create(&module).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetModulesByName(moduleName string) ([]ds.Module, error) {
	var modules []ds.Module

	err := r.db.Where("LOWER(modules.name) LIKE ?", "%"+strings.ToLower(moduleName)+"%").
		Where("is_deleted = ?", false).
		Find(&modules).Error
	if err != nil {
		return nil, err
	}
	return modules, nil
}

func (r *Repository) SaveModule(module *ds.Module) error {
	err := r.db.Save(module).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) AddToMission(missionId, moduleId string) error {
	Flight := ds.Flight{MissionId: missionId, ModuleId: moduleId}
	err := r.db.Create(&Flight).Error
	if err != nil {
		return err
	}
	return nil
}
