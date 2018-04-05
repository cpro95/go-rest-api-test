package main

import (
	"database/sql"
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
	r := gin.Default()
	r.Use(Cors())

	v1 := r.Group("api/v1")
	{
		// v1.POST("/users", PostUser)
		v1.GET("/movies", GetMovies)
		v1.GET("/movies/search", GetMovie)
		// v1.PUT("/users/:id", UpdateUser)
		// v1.DELETE("/users/:id", DeleteUser)
		v1.OPTIONS("/movies", OptionsMovie)
		v1.OPTIONS("/movies/:id", OptionsMovie)
	}

	r.Run(":8080")
}

func OptionsMovie(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Methods", "DELETE, POST, PUT")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	c.Next()
}

// func PostUser(c *gin.Context) {
// 	db := InitDb()
// 	defer db.Close()

// 	var user Users
// 	c.Bind(&user)

// 	if user.Firstname != "" && user.Lastname != "" {
// 		// INSERT INTO "users" (name) VALUES (user.Name);
// 		db.Create(&user)

// 		// Display status
// 		c.JSON(201, gin.H{"success": user})
// 	} else {
// 		// Display error
// 		c.JSON(422, gin.H{"error": "Fields are empty"})
// 	}

// 	// curl -i -X POST -H "Content-Type: application/json" -d "{\"firstname\":\"Thea\",\"lastname\":\"Queen\" }" http://localhost:8080/api/v1/users
// }

func GetMovies(c *gin.Context) {
	// Connection to the database
	db := InitDb()

	// Close connection database
	defer db.Close()

	var movies []Movie
	var movie Movie

	// // SELECT * FROM movie
	// // db.Table("movie").Where("c00 like ?", "%interstellar%").Select("idMovie, c00, c01, c03, c05, c07, c08, c20, c21").Limit(10).Find(&movies)
	// db.Table("movie").Select("idFile, c00, c01, c03, c05, c07, c08, c20, c21").Find(&movies)
	// //db.Table("movie").Limit(10).Find(&movies)
	// fmt.Println(movies)

	// query
	rows, err := db.Query("SELECT idMovie, c00, c01, c03, c05, c07, c08, c20, c21 FROM movie")
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

// func UpdateUser(c *gin.Context) {
// 	// Connection to the database
// 	db := InitDb()
// 	// Close connection database
// 	defer db.Close()

// 	// Get id user
// 	id := c.Params.ByName("id")
// 	var user Users
// 	// SELECT * FROM users WHERE id = 1;
// 	db.First(&user, id)

// 	if user.Firstname != "" && user.Lastname != "" {
// 		if user.Id != 0 {
// 			var newUser Users
// 			c.Bind(&newUser)

// 			result := Users{
// 				Id:        user.Id,
// 				Firstname: newUser.Firstname,
// 				Lastname:  newUser.Lastname,
// 			}
// 			// UPDATE users SET firstname='newUser.Firstname', lastname='newUser.Lastname' WHERE id = user.Id;
// 			db.Save(&result)

// 			// Display modified data in JSON message "success"
// 			c.JSON(200, gin.H{"success": result})
// 		} else {
// 			// Display JSON error
// 			c.JSON(404, gin.H{"error": "User not found"})
// 		}

// 	} else {
// 		// Display JSON error
// 		c.JSON(422, gin.H{"error": "Fields are empty"})
// 	}

// 	// curl -i -X PUT -H "Content-Type: application/json" -d "{ \"firstname\":\"Thea\", \"lastname\": \"Merlyn\" }" http://localhost:8080/api/v1/users/1
// }

// func DeleteUser(c *gin.Context) {
// 	// Connection to the database
// 	db := InitDb()

// 	// Close connection database
// 	defer db.Close()

// 	// Get id user
// 	id := c.Params.ByName("id")
// 	var user Users
// 	// SELECT * FROM users WHERE id = 1;
// 	db.First(&user, id)

// 	if user.Id != 0 {
// 		// DELETE FROM users WHERE id = user.Id
// 		db.Delete(&user)
// 		// Display JSON result
// 		c.JSON(200, gin.H{"success": "User #" + id + " deleted"})
// 	} else {
// 		// Display JSON error
// 		c.JSON(404, gin.H{"error": "User not found"})
// 	}

// 	// curl -i -X DELETE http://localhost:8080/api/v1/users/1
// }

// func InitDb() *gorm.DB {
// 	// Openning file
// 	db, err := gorm.Open("sqlite3", "./MyVideos99.db")
// 	db.LogMode(true)

// 	// Error
// 	if err != nil {
// 		panic(err)
// 	}

// 	// Creating the table
// 	// if !db.HasTable(&Users{}) {
// 	// 	db.CreateTable(&Users{})
// 	// 	db.Set("gorm:table_options", "ENGINE=InnoDB").CreateTable(&Users{})
// 	//

// 	return db
// }

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
