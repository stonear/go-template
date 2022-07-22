package service

import (
	"context"
	"database/sql"

	"github.com/stonear/go-template/entity"
	"github.com/stonear/go-template/helper"
	"github.com/stonear/go-template/repository"
)

type Service interface {
	Create(ctx context.Context, person entity.Person) (int, error)
	Update(ctx context.Context, person entity.Person) (int, error)
	Delete(ctx context.Context, id int) (int, error)
	FindById(ctx context.Context, id int) entity.Person
	FindAll(ctx context.Context) []entity.Person
}

func New(repo repository.Repository, db *sql.DB) Service {
	return &service{
		Repository: repo,
		DB:         db,
	}
}

type service struct {
	Repository repository.Repository
	DB         *sql.DB
}

func (s *service) Create(ctx context.Context, person entity.Person) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (s *service) Update(ctx context.Context, person entity.Person) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (s *service) Delete(ctx context.Context, id int) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (s service) FindById(ctx context.Context, id int) entity.Person {
	//TODO implement me
	panic("implement me")
}

func (s service) FindAll(ctx context.Context) []entity.Person {
	tx, err := s.DB.Begin()
	helper.Panic(err)
	persons := s.Repository.FindAll(ctx, tx)
	return persons
}
