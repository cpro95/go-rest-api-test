package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

type Movie struct {
	id    int
	title string
}

func main() {
	db, err := sql.Open("sqlite3", "./MyVideos99.db")
	checkErr(err)

	//query
	rows, err := db.Query("SELECT idMovie, c00 from movie where c00 like ?", "%Inter%")
	checkErr(err)
	var movies []Movie
	var movie Movie
	for rows.Next() {
		err = rows.Scan(&movie.id, &movie.title)
		checkErr(err)
		movies = append(movies, movie)
	}
	fmt.Println(movies)
	rows.Close()
}
