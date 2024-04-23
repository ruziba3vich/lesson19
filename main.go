package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
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

	for _, dbName := range dbNames {
		name := "sql/db/" + dbName + ".sql"
		sqlFile, err := os.ReadFile(name)
		if err != nil {
			log.Fatal(err)
		}

		_, err = db.Exec(string(sqlFile))
		if err != nil {
			log.Fatal(err)
		}
	}

	var wg sync.WaitGroup
	for i := range dbNames {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if err := LoadFromDb(db, dbNames[i], &buildings[i]); err != nil {
				log.Println(err)
			}
		}(i)
	}
	wg.Wait()

	var fullnameOfAllUsers []string
	ch := make(chan string)

	for _, building := range buildings {
		go FanIn(building, ch)
	}

	go func() {
		for _, building := range buildings {
			FanIn(building, ch)
		}
		close(ch)
	}()

	for fullname := range ch {
		fullnameOfAllUsers = append(fullnameOfAllUsers, fullname)
	}

	if err := WriteToFile("fullnames.txt", fullnameOfAllUsers); err != nil {
		log.Fatal(err)
	}
}

func FanIn(participants []string, ch chan<- string) {
	for _, participant := range participants {
		ch <- participant
	}
}

func LoadFromDb(db *sql.DB, dbName string, data *[]string) error {
	query := fmt.Sprintf("SELECT fullname FROM %s;", dbName)
	rows, err := db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var fullname string
		if err := rows.Scan(&fullname); err != nil {
			return err
		}
		*data = append(*data, fullname)
	}
	if err := rows.Err(); err != nil {
		return err
	}

	return nil
}

func WriteToFile(filename string, data []string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, line := range data {
		_, err := fmt.Fprintln(file, line)
		if err != nil {
			return err
		}
	}

	return nil
}
