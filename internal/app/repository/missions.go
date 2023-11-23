package repository

import (
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"

	"lab1/internal/app/ds"
)

func (r *Repository) GetAllMissions(dateApproveStart, dateApproveEnd *time.Time, status string) ([]ds.Mission, error) {
	var missions []ds.Mission
	var err error

	if dateApproveStart != nil && dateApproveEnd != nil {
		err = r.db.Preload("Customer").Preload("Moderator").
			Where("LOWER(status) LIKE ?", "%"+strings.ToLower(status)+"%").
			Where("date_approve BETWEEN ? AND ?", *dateApproveStart, *dateApproveEnd).
			Find(&missions).Error
	} else if dateApproveStart != nil {
		err = r.db.Preload("Customer").Preload("Moderator").
			Where("LOWER(status) LIKE ?", "%"+strings.ToLower(status)+"%").
			Where("date_approve >= ?", *dateApproveStart).
			Find(&missions).Error
	} else if dateApproveEnd != nil {
		err = r.db.Preload("Customer").Preload("Moderator").
			Where("LOWER(status) LIKE ?", "%"+strings.ToLower(status)+"%").
			Where("date_approve <= ?", *dateApproveEnd).
			Find(&missions).Error
	} else {
		err = r.db.Preload("Customer").Preload("Moderator").
			Where("LOWER(status) LIKE ?", "%"+strings.ToLower(status)+"%").
			Find(&missions).Error
	}
	if err != nil {
		return nil, err
	}
	return missions, nil
}

func (r *Repository) GetDraftMission(customerId string) (*ds.Mission, error) {
	mission := &ds.Mission{}
	err := r.db.First(mission, ds.Mission{Status: ds.DRAFT, CustomerId: customerId}).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return mission, nil
}

func (r *Repository) CreateDraftMission(customerId string) (*ds.Mission, error) {
	mission := &ds.Mission{DateCreated: time.Now(), CustomerId: customerId, Status: ds.DRAFT}
	err := r.db.Create(mission).Error
	if err != nil {
		return nil, err
	}
	return mission, nil
}

func (r *Repository) GetMissionById(missionId, customerId string) (*ds.Mission, error) {
	mission := &ds.Mission{}
	err := r.db.Preload("Moderator").Preload("Customer").
		First(mission, ds.Mission{UUID: missionId, CustomerId: customerId}).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return mission, nil
}

func (r *Repository) GetFlight(missionId string) ([]ds.Module, error) {
	var modules []ds.Module

	err := r.db.Table("flights").
		Select("modules.*").
		Joins("JOIN modules ON flights.module_id = modules.uuid").
		Where(ds.Flight{MissionId: missionId}).
		Scan(&modules).Error

	if err != nil {
		return nil, err
	}
	return modules, nil
}

func (r *Repository) SaveMission(mission *ds.Mission) error {
	err := r.db.Save(mission).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) DeleteFromMission(missionId, moduleId string) error {
	err := r.db.Delete(&ds.Flight{MissionId: missionId, ModuleId: moduleId}).Error
	if err != nil {
		return err
	}
	return nil
}
