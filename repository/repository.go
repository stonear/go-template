package repository

import (
	"context"
	"database/sql"

	"github.com/stonear/go-template/entity"
	"github.com/stonear/go-template/helper"
)

type Repository interface {
	Save(ctx context.Context, tx *sql.Tx, person entity.Person) (int, error)
	Update(ctx context.Context, tx *sql.Tx, person entity.Person) (int, error)
	Delete(ctx context.Context, tx *sql.Tx, id int) (int, error)
	FindById(ctx context.Context, tx *sql.Tx, id int) entity.Person
	FindAll(ctx context.Context, tx *sql.Tx) []entity.Person
}

func New() Repository {
	return &repository{}
}

type repository struct{}

func (r *repository) Save(ctx context.Context, tx *sql.Tx, person entity.Person) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (r *repository) Update(ctx context.Context, tx *sql.Tx, person entity.Person) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (r *repository) Delete(ctx context.Context, tx *sql.Tx, id int) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (r repository) FindById(ctx context.Context, tx *sql.Tx, id int) entity.Person {
	query := "SELECT id, name FROM person WHERE id = $1"
	rows, err := tx.QueryContext(ctx, query, id)
	helper.Panic(err)
	person := entity.Person{}
	if rows.Next() {
		err := rows.Scan(&person.ID, &person.Name)
		helper.Panic(err)
	}
	return person
}

func (r repository) FindAll(ctx context.Context, tx *sql.Tx) []entity.Person {
	query := "SELECT id, name FROM person"
	rows, err := tx.QueryContext(ctx, query)
	helper.Panic(err)

	var persons []entity.Person
	for rows.Next() {
		person := entity.Person{}
		err := rows.Scan(&person.ID, &person.Name)
		helper.Panic(err)
		persons = append(persons, person)
	}

	return persons
}
