package repository

import (
	"context"
	"database/sql"

	"github.com/stonear/go-template/entity"
	"github.com/stonear/go-template/helper"
)

type Repository interface {
	Index(ctx context.Context, tx *sql.Tx) []entity.Person
	Show(ctx context.Context, tx *sql.Tx, id int) entity.Person
	Store(ctx context.Context, tx *sql.Tx, person entity.Person) (int, error)
	Update(ctx context.Context, tx *sql.Tx, person entity.Person) (int, error)
	Destroy(ctx context.Context, tx *sql.Tx, id int) (int, error)
}

func New() Repository {
	return &repository{}
}

type repository struct{}

func (r repository) Index(ctx context.Context, tx *sql.Tx) []entity.Person {
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

func (r repository) Show(ctx context.Context, tx *sql.Tx, id int) entity.Person {
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

func (r *repository) Store(ctx context.Context, tx *sql.Tx, person entity.Person) (int, error) {
	var id int
	query := "INSERT INTO person(name) VALUES ($1) RETURNING id"
	row := tx.QueryRowContext(ctx, query, person.Name)
	err := row.Scan(&id)
	helper.Panic(err)
	return id, err
}

func (r *repository) Update(ctx context.Context, tx *sql.Tx, person entity.Person) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (r *repository) Destroy(ctx context.Context, tx *sql.Tx, id int) (int, error) {
	//TODO implement me
	panic("implement me")
}
