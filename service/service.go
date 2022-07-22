package service

import (
	"context"
	"database/sql"

	"github.com/stonear/go-template/entity"
	"github.com/stonear/go-template/helper"
	"github.com/stonear/go-template/repository"
)

type Service interface {
	Index(ctx context.Context) []entity.Person
	Show(ctx context.Context, id int) entity.Person
	Store(ctx context.Context, person entity.Person) (int, error)
	Update(ctx context.Context, person entity.Person) (int, error)
	Destroy(ctx context.Context, id int) (int, error)
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

func (s service) Index(ctx context.Context) []entity.Person {
	tx, err := s.DB.Begin()
	helper.Panic(err)
	persons := s.Repository.Index(ctx, tx)
	return persons
}

func (s service) Show(ctx context.Context, id int) entity.Person {
	tx, err := s.DB.Begin()
	helper.Panic(err)
	person := s.Repository.Show(ctx, tx, id)
	return person
}

func (s *service) Store(ctx context.Context, person entity.Person) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (s *service) Update(ctx context.Context, person entity.Person) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (s *service) Destroy(ctx context.Context, id int) (int, error) {
	//TODO implement me
	panic("implement me")
}
