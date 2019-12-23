package postgres

import (
	"context"
	"fmt"
	"github.com/go-pg/pg/v9"
	_ "github.com/lib/pq"
)

type dbLogger struct{}

func (d dbLogger) BeforeQuery(c context.Context, q *pg.QueryEvent) (context.Context, error) {
	return c, nil
}

func (d dbLogger) AfterQuery(c context.Context, q *pg.QueryEvent) error {
	fmt.Println(q.FormattedQuery())
	return nil
}

func New(opts *pg.Options) *pg.DB {
	db := pg.Connect(opts)
	db.AddQueryHook(dbLogger{})

	return db
}
