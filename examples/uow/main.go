package main

import (
	"context"
	"database/sql"
	"log"

	_ "github.com/lib/pq"

	"github.com/xfrr/go-cqrsify/uow"

	pg "github.com/xfrr/go-cqrsify/uow/postgres"
)

// Domain-agnostic registry for the app:
type Repos struct {
	Users  UserRepo
	Orders OrderRepo
}

type User struct {
	ID          int64
	Email, Name string
}

type UserRepo interface {
	Register(ctx context.Context, u *User) error
}

// Example repositories (simplified)
type userRepo struct {
	exec interface {
		ExecContext(context.Context, string, ...any) (sql.Result, error)
		QueryRowContext(context.Context, string, ...any) *sql.Row
	}
}

func newUserRepoExec(exec interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
	QueryRowContext(context.Context, string, ...any) *sql.Row
}) UserRepo {
	return &userRepo{exec: exec}
}

func (r *userRepo) Register(ctx context.Context, u *User) error {
	_, err := r.exec.ExecContext(ctx, `INSERT INTO users(email,name) VALUES ($1,$2)`, u.Email, u.Name)
	return err
}

type OrderRepo interface {
	Create(ctx context.Context, orderID int64) error
}

type orderRepo struct {
	exec interface {
		ExecContext(context.Context, string, ...any) (sql.Result, error)
	}
}

func newOrderRepoExec(exec interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
}) OrderRepo {
	return &orderRepo{exec: exec}
}

func (r *orderRepo) Create(ctx context.Context, orderID int64) error {
	_, err := r.exec.ExecContext(ctx, `INSERT INTO orders(id) VALUES ($1)`, orderID)
	return err
}

// wireBase constructs a Repos with "no tx" repositories
func wireBase(db *sql.DB) Repos {
	return Repos{Users: newUserRepoExec(db), Orders: newOrderRepoExec(db)}
}

// wireTx binds a Tx to a Repos with tx-scoped repositories
func wireTx(tx uow.Tx) Repos {
	// Adapt wrapped tx to *sql.Tx to reuse exec methods
	tw, ok := pg.Unwrap(tx)
	if !ok {
		return Repos{}
	}
	// Return a Repos with tx-bound repos
	return Repos{Users: newUserRepoExec(tw), Orders: newOrderRepoExec(tw)}
}

func registerUserAndCreateOrderSuccess(ctx context.Context, r Repos, email, name string) error {
	userErr := r.Users.Register(ctx, &User{Email: email, Name: name})
	if userErr != nil {
		return userErr
	}

	orderErr := r.Orders.Create(ctx, 12345)
	if orderErr != nil {
		return orderErr
	}
	return nil
}

func registerUserAndCreateOrderFail(ctx context.Context, r Repos, email, name string) error {
	userErr := r.Users.Register(ctx, &User{Email: email, Name: name})
	if userErr != nil {
		return userErr
	}

	// This will fail due to duplicate order ID
	orderErr := r.Orders.Create(ctx, 12345)
	if orderErr != nil {
		return orderErr
	}
	return nil
}

// NOTE: The "cqrsify_uow_example" database must exist, and the "users" table must be created beforehand:
// Example table creation SQL:
//
// DROP TABLE IF EXISTS users;
// CREATE DATABASE cqrsify_uow_example;
// CREATE TABLE users (id SERIAL PRIMARY KEY, email TEXT UNIQUE NOT NULL, name TEXT NOT NULL);
func main() {
	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:5432/cqrsify_uow_example?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	mgr := pg.NewManager(db)
	u := uow.New(mgr, wireBase(db), wireTx, uow.Config{EnableSavepoints: true})

	if err = registerUserAndCreateOrderSuccess(context.Background(), u.Repos(), "jane@example.com", "Jane"); err != nil {
		panic(err)
	}

	log.Println("User registered and order created successfully")

	if err = registerUserAndCreateOrderFail(context.Background(), u.Repos(), "jane@example.com", "Jane"); err != nil {
		log.Printf("Expected failure: %v", err)
	}

	log.Println("Transaction rolled back due to error, no changes must have been applied")
}
