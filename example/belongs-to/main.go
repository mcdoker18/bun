package main

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/extra/bundebug"
)

// Profile belongs to User.
type Profile struct {
	ID     int64
	Lang   string
	UserID int64
}

type User struct {
	ID      int64
	Name    string
	Profile *Profile `bun:"rel:belongs-to"`
}

func main() {
	ctx := context.Background()

	sqldb, err := sql.Open("sqlite3", "file::memory:?cache=shared")
	if err != nil {
		panic(err)
	}

	db := bun.NewDB(sqldb, sqlitedialect.New())
	defer db.Close()

	db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose()))

	if err := createSchema(ctx, db); err != nil {
		panic(err)
	}

	var users []User
	if err := db.NewSelect().
		Model(&users).
		Column("user.*").
		Relation("Profile").
		Scan(ctx); err != nil {
		panic(err)
	}

	fmt.Println(len(users), "results")
	fmt.Println(users[0].ID, users[0].Name, users[0].Profile)
	fmt.Println(users[1].ID, users[1].Name, users[1].Profile)
	// Output: 2 results
	// 1 user 1 &{1 en 1}
	// 2 user 2 &{2 ru 2}
}

func createSchema(ctx context.Context, db *bun.DB) error {
	models := []interface{}{
		(*User)(nil),
		(*Profile)(nil),
	}
	for _, model := range models {
		if _, err := db.NewCreateTable().Model(model).Exec(ctx); err != nil {
			return err
		}
	}

	users := []*User{
		{ID: 1, Name: "user 1"},
		{ID: 2, Name: "user 2"},
	}
	if _, err := db.NewInsert().Model(&users).Exec(ctx); err != nil {
		return err
	}

	profiles := []*Profile{
		{ID: 1, Lang: "en", UserID: 1},
		{ID: 2, Lang: "ru", UserID: 2},
	}
	if _, err := db.NewInsert().Model(&profiles).Exec(ctx); err != nil {
		return err
	}

	return nil
}