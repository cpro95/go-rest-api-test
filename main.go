package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

type Movie struct {
	IdMovie int     `json:"id"`
	C00     string  `json:"title"`
	C01     string  `json:"overview"`
	C03     string  `json:"tagline"`
	C05     float64 `json:"rating"`
	C07     int     `json:"year"`
	C08     string  `json:"poster"`
	C20     string  `json:"fanart"`
	C21     string  `json:"country"`
}

func main() {
	fmt.Println("Launching MyMovie REST API...........")
	r := gin.Default()
	r.Use(Cors())

	v1 := r.Group("api/v1")
	{
		v1.GET("/movies", GetMovies)
		v1.GET("/movies/search", GetMovie)
		v1.OPTIONS("/movies", OptionsMovie)
	}

	r.Run(":8080")
}

func OptionsMovie(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Methods", "DELETE, POST, PUT")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	c.Next()
}

func GetMovies(c *gin.Context) {
	// Connection to the database
	db := InitDb()

	// Close connection database
	defer db.Close()

	var movies []Movie
	var movie Movie

	// Query string parameters are parsed using the existing underlying request object.
	// The request responds to a url matching:  /search?id=123&name=string
	// c.Query("lastname") // shortcut for c.Request.URL.Query().Get("lastname")
	limitQuery := c.Request.URL.Query().Get("limit")
	offsetQuery := c.Request.URL.Query().Get("offset")
	// fmt.Println("id : " + limitQuery)
	// fmt.Println("name : " + offsetQuery)

	// query
	rows, err := db.Query("SELECT idMovie, c00, c01, c03, c05, c07, c08, c20, c21 FROM movie order by idMovie asc limit ? offset ?", limitQuery, offsetQuery)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		err := rows.Scan(&movie.IdMovie, &movie.C00, &movie.C01, &movie.C03, &movie.C05, &movie.C07, &movie.C08, &movie.C20, &movie.C21)
		if err != nil {
			panic(err)
		}
		movies = append(movies, movie)
	}
	rows.Close()

	// Display JSON result
	c.JSON(200, movies)

	// curl -i http://localhost:8080/api/v1/users
}

func GetMovie(c *gin.Context) {
	// Connection to the database
	db := InitDb()

	// Close connection database
	defer db.Close()

	var movie Movie
	var movies []Movie

	// Query string parameters are parsed using the existing underlying request object.
	// The request responds to a url matching:  /search?id=123&name=string
	// c.Query("lastname") // shortcut for c.Request.URL.Query().Get("lastname")
	idQuery := c.Request.URL.Query().Get("id")
	nameQuery := c.Request.URL.Query().Get("name")
	// fmt.Println("id : " + idQuery)
	// fmt.Println("name : " + nameQuery)
	// fmt.Println(idQuery)
	// fmt.Println(reflect.TypeOf(idQuery))
	// fmt.Println(reflect.TypeOf(nameQuery))

	var err error
	if idQuery == "" && nameQuery == "" {
		// fmt.Println("id, name : both arguments are empty")
		err = db.QueryRow("SELECT idMovie, c00, c01, c03, c05, c07, c08, c20, c21 FROM movie WHERE c00 LIKE ?", "%Interstellar%").Scan(&movie.IdMovie, &movie.C00, &movie.C01, &movie.C03, &movie.C05, &movie.C07, &movie.C08, &movie.C20, &movie.C21)
	} else if idQuery != "" && nameQuery == "" {
		// fmt.Println("id is ok, but name argument is empty")
		err = db.QueryRow("SELECT idMovie, c00, c01, c03, c05, c07, c08, c20, c21 FROM movie WHERE idMovie = ?", idQuery).Scan(&movie.IdMovie, &movie.C00, &movie.C01, &movie.C03, &movie.C05, &movie.C07, &movie.C08, &movie.C20, &movie.C21)
	} else if idQuery == "" && nameQuery != "" {
		// fmt.Println("id is blank and name is not empty")
		rows, err2 := db.Query("SELECT idMovie, c00, c01, c03, c05, c07, c08, c20, c21 FROM movie WHERE c00 LIKE ?", "%"+nameQuery+"%")
		if err2 != nil {
			panic(err2)
		}
		for rows.Next() {
			err3 := rows.Scan(&movie.IdMovie, &movie.C00, &movie.C01, &movie.C03, &movie.C05, &movie.C07, &movie.C08, &movie.C20, &movie.C21)
			if err3 != nil {
				panic(err3)
			}
			movies = append(movies, movie)
		}

	} else {
		// fmt.Println("id, name is full")
		err = db.QueryRow("SELECT idMovie, c00, c01, c03, c05, c07, c08, c20, c21 FROM movie WHERE idMovie = ? and c00 like ?", idQuery, "%"+nameQuery+"%").Scan(&movie.IdMovie, &movie.C00, &movie.C01, &movie.C03, &movie.C05, &movie.C07, &movie.C08, &movie.C20, &movie.C21)
	}

	// query
	// err2 := db.QueryRow("SELECT idMovie, c00, c01, c03, c05, c07, c08, c20, c21 FROM movie WHERE idMovie = ?", Idmovie).Scan(&movie.IdMovie, &movie.C00, &movie.C01, &movie.C03, &movie.C05, &movie.C07, &movie.C08, &movie.C20, &movie.C21)

	if err != nil {
		// Display JSON error
		c.JSON(404, gin.H{"error": "Movie not found"})
	} else {
		// fmt.Println(len(movies))
		if len(movies) == 0 {
			c.JSON(200, movie)
		} else {
			c.JSON(200, movies)
		}

	}

	// curl -i http://localhost:8080/api/v1/users/1
}

func InitDb() *sql.DB {
	// Openning file
	db, err := sql.Open("sqlite3", "./MyVideos99.db")

	// Error
	if err != nil {
		panic(err)
	}

	return db
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.Next()
	}
}
