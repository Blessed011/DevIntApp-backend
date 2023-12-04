package repository

import (
	"lab1/internal/app/ds"
)

func (r *Repository) AddUser(user *ds.User) error {
	return r.db.Create(user).Error
}

func (r *Repository) GetUserByLogin(login string) (*ds.User, error) {
	user := &ds.User{}
	if err := r.db.Where("login = ?", login).
		First(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}
