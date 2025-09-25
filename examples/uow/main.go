package main

import (
	"context"
	"crypto/rand"
	"database/sql"
	"log"
	"math/big"

	_ "github.com/lib/pq"

	"github.com/xfrr/go-cqrsify/uow"
	uowpg "github.com/xfrr/go-cqrsify/uow/postgres"
)

// NOTE: Postgres DB and Tables must be created before running this example.
// Execute the init.sh script in this folder to deploy the Postgres docker container
// and create the required resources.
func main() {
	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:5432/cqrsify_uow_example?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	mgr := uowpg.NewManager(db)
	u := uow.New(mgr, wireBase(db), wireTx, uow.Config{EnableSavepoints: true})

	if err = registerUserAndCreateOrder(context.Background(), u, "Jane", nil); err != nil {
		panic(err)
	}

	if err = registerUserAndCreateOrder(context.Background(), u, "Jane", sql.ErrConnDone); err != nil {
		log.Printf("Expected failure: %v", err)
	}
}

// Domain-agnostic registry for the app:
type Repos struct {
	Users  UserRepo
	Orders OrderRepo
}

type User struct {
	ID   int64
	Name string
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
	_, err := r.exec.ExecContext(ctx, `INSERT INTO users(id, name) VALUES ($1,$2)`, u.ID, u.Name)
	return err
}

type OrderRepo interface {
	Create(ctx context.Context, orderID int64, userID int64, amount float64) error
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

func (r *orderRepo) Create(ctx context.Context, orderID int64, userID int64, amount float64) error {
	_, err := r.exec.ExecContext(ctx, `INSERT INTO orders(id,user_id,amount) VALUES ($1,$2,$3)`, orderID, userID, amount)
	return err
}

// wireBase constructs a Repos with "no tx" repositories
func wireBase(db *sql.DB) Repos {
	return Repos{Users: newUserRepoExec(db), Orders: newOrderRepoExec(db)}
}

// wireTx binds a Tx to a Repos with tx-scoped repositories
func wireTx(tx uow.Tx) Repos {
	// Adapt wrapped tx to *sql.Tx to reuse exec methods
	tw, ok := uowpg.Unwrap(tx)
	if !ok {
		return Repos{}
	}
	// Return a Repos with tx-bound repos
	return Repos{Users: newUserRepoExec(tw), Orders: newOrderRepoExec(tw)}
}

func registerUserAndCreateOrder(ctx context.Context, u *uow.UnitOfWork[Repos], name string, err error) error {
	return u.Do(ctx, func(ctx context.Context, r Repos) error {
		userID := randomID()

		// Register a new user
		userErr := r.Users.Register(ctx, &User{
			ID:   userID,
			Name: name,
		})
		if userErr != nil {
			return userErr
		}

		log.Printf("User %s (%d) registered successfully", name, userID)

		// Create an order for the user
		orderID := randomID()
		orderAmount := 99.99
		orderErr := r.Orders.Create(ctx, orderID, userID, orderAmount)
		if orderErr != nil {
			return orderErr
		}

		if err != nil {
			return err
		}

		log.Printf("Order (%d) created successfully", orderID)
		return nil
	})
}

func randomID() int64 {
	x := 100000
	nBig, err := rand.Int(rand.Reader, big.NewInt(int64(x)))
	if err != nil {
		log.Fatal(err)
	}
	return nBig.Int64()
}
