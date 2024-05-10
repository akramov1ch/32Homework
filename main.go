package main

import (
    "context"
    "database/sql"
    "fmt"
    "log"
    "math/rand"
    "time"
    _ "github.com/lib/pq"
    "golang.org/x/sync/errgroup"
)

func main() {
    db, err := sql.Open("postgres", "user=user password=password dbname=database sslmode=disable")
    if err != nil {
        log.Fatalf("Xatolik: %v", err)
    }
    defer db.Close()
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    g, ctx := errgroup.WithContext(ctx)
    g.Go(func() error {
        _, err := db.ExecContext(ctx, "INSERT INTO large_dataset (data) VALUES ($1)", rand.Int())
        return err
    })
    g.Go(func() error {
        rows, err := db.QueryContext(ctx, "SELECT data FROM large_dataset")
        if err != nil {
            return err
        }
        defer rows.Close()
        for rows.Next() {
            var data int
            err = rows.Scan(&data)
            if err != nil {
                return err
            }
            fmt.Println("Data:", data)
        }
        return nil
    })
    g.Go(func() error {
        _, err := db.ExecContext(ctx, "UPDATE large_dataset SET data = $1 WHERE data = $2", rand.Int(), rand.Int())
        return err
    })
    if err := g.Wait(); err != nil {
        log.Fatal(err)
    }
}
