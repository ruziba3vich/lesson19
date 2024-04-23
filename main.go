package main

import (
	"database/sql"
	"fmt"
	"lesson19/packages"
	"log"
	"sync"

	_ "github.com/lib/pq"
)

func main() {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", "localhost", 5432, "postgres", "Dost0n1k", "participantsDB")
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	dbNames := []string{"building1", "building2", "building3", "building4", "building5"}
	buildings := make([][]string, len(dbNames))

	var wg sync.WaitGroup
	for i, dbName := range dbNames {
		wg.Add(1)
		go func(i int, dbName string) {
			defer wg.Done()
			if err := packages.LoadFromDb(db, dbName, &buildings[i]); err != nil {
				log.Printf("Error loading data from %s: %v", dbName, err)
			}
		}(i, dbName)
	}
	wg.Wait()

	fullnameOfAllUsers := make([]string, 0)

	ch := make(chan string)
	go func() {
		defer close(ch)
		for _, building := range buildings {
			packages.FanIn(building, ch)
		}
	}()

	for fullname := range ch {
		fullnameOfAllUsers = append(fullnameOfAllUsers, fullname)
	}

	if err := packages.WriteToFile("fullnames.txt", fullnameOfAllUsers); err != nil {
		log.Fatalf("Error writing to file: %v", err)
	}
}
