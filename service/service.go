package service

import (
	"context"
	"database/sql"
	"github.com/stonear/go-template/database"

	"github.com/stonear/go-template/entity"
	"github.com/stonear/go-template/helper"
	"github.com/stonear/go-template/repository"
)

type Service interface {
	Index(ctx context.Context) []entity.Person
	Show(ctx context.Context, id int) entity.Person
	Store(ctx context.Context, person entity.Person) (int, error)
	Update(ctx context.Context, id int, person entity.Person) (entity.Person, error)
	Destroy(ctx context.Context, id int) error
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
	defer database.Commit(tx)
	persons := s.Repository.Index(ctx, tx)
	err = tx.Commit()
	helper.Panic(err)
	return persons
}

func (s service) Show(ctx context.Context, id int) entity.Person {
	tx, err := s.DB.Begin()
	helper.Panic(err)
	defer database.Commit(tx)
	person := s.Repository.Show(ctx, tx, id)
	err = tx.Commit()
	helper.Panic(err)
	return person
}

func (s *service) Store(ctx context.Context, person entity.Person) (int, error) {
	tx, err := s.DB.Begin()
	helper.Panic(err)

	id, err := s.Repository.Store(ctx, tx, person)
	helper.Panic(err)
	err = tx.Commit()

	return id, err
}

func (s *service) Update(ctx context.Context, id int, person entity.Person) (entity.Person, error) {
	tx, err := s.DB.Begin()
	helper.Panic(err)
	person, err = s.Repository.Update(ctx, tx, id, person)
	helper.Panic(err)
	err = tx.Commit()
	return person, err
}

func (s *service) Destroy(ctx context.Context, id int) error {
	tx, err := s.DB.Begin()
	helper.Panic(err)
	err = s.Repository.Destroy(ctx, tx, id)
	helper.Panic(err)
	err = tx.Commit()
	return err
}
