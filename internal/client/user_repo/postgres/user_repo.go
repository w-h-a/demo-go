package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/w-h-a/demo-go/api/user"
	userrepo "github.com/w-h-a/demo-go/internal/client/user_repo"
)

var DRIVER string

func init() {
	// TODO: register with otel and set driver

	DRIVER = "postgres"
}

// pgUserRepo is an implementation of UserRepo
type pgUserRepo struct {
	options userrepo.Options
	conn    *sql.DB
}

// Create inserts a new user into the db
func (ur *pgUserRepo) Create(ctx context.Context, dto user.CreateUserDTO) (user.User, error) {
	u := user.User{
		ID:    uuid.NewString(),
		Name:  dto.Name,
		Email: dto.Email,
	}

	query := `INSERT INTO users (id, name, email) VALUES ($1, $2, $3)`

	if _, err := ur.conn.ExecContext(ctx, query, u.ID, u.Name, u.Email); err != nil {
		return user.User{}, err
	}

	return u, nil
}

// GetByID retrieves a user from the db given their ID.
func (ur *pgUserRepo) GetByID(ctx context.Context, id string) (user.User, error) {
	query := `SELECT id, name, email FROM users WHERE id = $1`

	row := ur.conn.QueryRowContext(ctx, query, id)

	var u user.User

	if err := row.Scan(&u.ID, &u.Name, &u.Email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user.User{}, userrepo.ErrUserNotFound
		}
		return user.User{}, err
	}

	return u, nil
}

// GetByEmail retrieves a user from the db given their email.
func (ur *pgUserRepo) GetByEmail(ctx context.Context, email string) (user.User, error) {
	query := `SELECT id, name, email FROM users WHERE email = $1`

	row := ur.conn.QueryRowContext(ctx, query, email)

	var u user.User

	err := row.Scan(&u.ID, &u.Name, &u.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user.User{}, userrepo.ErrUserNotFound
		}
		return user.User{}, err
	}

	return u, nil
}

// GetAll retrieves all users that satisfy GetAllOptions from the db
func (ur *pgUserRepo) GetAll(ctx context.Context, opts ...userrepo.GetAllOption) ([]user.User, error) {
	// TODO: use opts to get sort, filters, etc

	query := `SELECT id, name, email FROM users`

	rows, err := ur.conn.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var us []user.User

	for rows.Next() {
		var u user.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email); err != nil {
			return nil, err
		}
		us = append(us, u)
	}

	return us, nil
}

// NewUserRepo creates a new pgUserRepo
func NewUserRepo(opts ...userrepo.Option) userrepo.UserRepo {
	options := userrepo.NewOptions(opts...)

	// TODO: validate options

	ur := &pgUserRepo{
		options: options,
	}

	conn, err := sql.Open(DRIVER, ur.options.Location)
	if err != nil {
		// log
		panic(err)
	}

	if err := conn.Ping(); err != nil {
		// log
		panic(err)
	}

	// TODO: start otel

	ur.conn = conn

	query := `
    CREATE TABLE IF NOT EXISTS users (
        id UUID PRIMARY KEY,
        name VARCHAR(100) NOT NULL,
        email VARCHAR(100) NOT NULL UNIQUE
    );
    `

	if _, err := ur.conn.Exec(query); err != nil {
		// log
		panic(err)
	}

	return ur
}
