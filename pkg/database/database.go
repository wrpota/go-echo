package database

import "gorm.io/gorm"

var _ Repository = (*dbRepository)(nil)

type Repository interface {
	i()
	GetDb() *gorm.DB
	DbClose() error
}

type dbRepository struct {
	engine *gorm.DB
}

func NewRepository() *dbRepository {

	return &dbRepository{
		engine: engine,
	}
}

func (r *dbRepository) i() {}

func (r *dbRepository) GetDb() *gorm.DB {
	return r.engine
}

func (r *dbRepository) DbClose() error {
	sqlDB, err := r.engine.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}
