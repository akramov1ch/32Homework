package main

import (
    "context"
    "database/sql"
    "fmt"
    "log"
    "math/rand"
    "sync"
    "time"

    _ "github.com/lib/pq" 
)

func main() {
    db, err := sql.Open("postgres", "user=user password=password dbname=database sslmode=disable")
    if err != nil {
        log.Fatalf("Xatolik: %v", err)
    }
    defer db.Close()

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    var wg sync.WaitGroup
    wg.Add(3)

    go func() {
        defer wg.Done()
        _, err := db.ExecContext(ctx, "INSERT INTO large_dataset (data) VALUES ($1)", rand.Int())
        if err != nil {
            log.Printf("Insert xatoligi: %v", err)
        }
    }()

    go func() {
        defer wg.Done()
        rows, err := db.QueryContext(ctx, "SELECT data FROM large_dataset")
        if err != nil {
            log.Printf("Select xatoligi: %v", err)
            return
        }
        defer rows.Close()

        for rows.Next() {
            var data int
            err = rows.Scan(&data)
            if err != nil {
                log.Printf("Scan xatoligi: %v", err)
                return
            }
            fmt.Println("Data:", data)
        }
    }()

    go func() {
        defer wg.Done()
        _, err := db.ExecContext(ctx, "UPDATE large_dataset SET data = $1 WHERE data = $2", rand.Int(), rand.Int())
        if err != nil {
            log.Printf("Update xatoligi: %v", err)
        }
    }()

    wg.Wait()
}
