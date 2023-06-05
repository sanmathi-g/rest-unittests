package main_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	controllers "github.com/sana/rest/controller"
	"github.com/sana/rest/database"
	. "github.com/smartystreets/goconvey/convey"
)

func TestMain(t *testing.T) {

	r := gin.Default()

	Convey("testing all the functions", t, func() {
		database.DatabaseConnection()
		r.GET("/movies/:id", controllers.ReadMovie)
		r.GET("/movies", controllers.ReadMovies)
		r.POST("/movies", controllers.CreateMovie)
		r.PUT("/moviess/:id", controllers.UpdateMovie)
		r.DELETE("/moviess/:id", controllers.DeleteMovie)

		Convey("when calling GET /movies/:id", func() {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/movies/:1", nil)
			r.ServeHTTP(w, req)

			Convey("Then respond should have status 200 ok", func() {
				So(w.Code, ShouldEqual, http.StatusOK)
			})
		})
		Convey("when calling GET /movies/:id with valid movie ID", func() {
			movie := database.Movie{Title: "TestMovie", Year: "2021"}
			database.DB.Create(&movie)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/movies/:id"+strconv.Itoa(int(movie.ID)), nil)
			r.ServeHTTP(w, req)

			Convey("Then respond should have status 200 ok", func() {
				So(w.Code, ShouldEqual, http.StatusOK)
			})
			Convey("then the response should contain the movie data", func() {
				var responseMovie database.Movie
				err := json.Unmarshal(w.Body.Bytes(), &responseMovie)

				So(err, ShouldBeNil)
				So(responseMovie.ID, ShouldEqual, movie.ID)
				So(responseMovie.Title, ShouldEqual, movie.Title)
				So(responseMovie.Year, ShouldEqual, movie.Year)

			})
			database.DB.Delete(&movie)
		})
		Convey("when calling GET /movies/:id with invalid movie ID", func() {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/movies/999", nil)
			r.ServeHTTP(w, req)

			Convey("Then respond should have status 404 Not Found", func() {
				So(w.Code, ShouldEqual, http.StatusNotFound)
			})
		})
		Convey("When Creating GET /movies", func() {
			//Create some movies for testing
			movies := []database.Movie{
				{Title: "Movie1", Year: "2021"},
				{Title: "Movie2", Year: "1998"},
			}
			database.DB.Create(&movies)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/movies", nil)
			r.ServeHTTP(w, req)

			Convey("Then the response should have status 200 OK", func() {
				So(w.Code, ShouldEqual, http.StatusOK)
			})
			Convey("Then the response should contain the list of movies", func() {
				var responseMovie []database.Movie
				err := json.Unmarshal(w.Body.Bytes(), &responseMovie)
				So(err, ShouldBeNil)
			})
			database.DB.Delete(&movies)
		})
		Convey("when calling GET /movies", func() {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/movies", nil)
			r.ServeHTTP(w, req)

			Convey("Then respond should have status 200 ok", func() {
				So(w.Code, ShouldEqual, http.StatusOK)
			})
		})
		Convey("when calling POST /movies", func() {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/movies", nil)
			r.ServeHTTP(w, req)

			Convey("Then respond should have status 400 bad request ok", func() {
				So(w.Code, ShouldEqual, http.StatusBadRequest)
			})
		})
		Convey("when calling POST /movies with valid data", func() {
			//create a new payload
			payload := `{"title":"NewMovie","year":"2021"}`

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/movies", strings.NewReader(payload))
			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, req)

			Convey("Then respond should have status 200 ok", func() {
				So(w.Code, ShouldEqual, http.StatusOK)
			})
			Convey("then the response shouls contain the created movie", func() {
				var responseMovie database.Movie
				err := json.Unmarshal(w.Body.Bytes(), &responseMovie)
				So(err, ShouldBeNil)

				So(responseMovie.Title, ShouldEqual, "New Movie")
				So(responseMovie.Year, ShouldEqual, "2021")
			})
			database.DB.Delete(&database.Movie{}, "title=?", "New Movie")
		})
		Convey("when calling delete /movies:id with a valid ID", func() {
			movie := database.Movie{Title: "Test Movie", Year: "2021"}
			database.DB.Create(&movie)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("DELETE", "/movies/:id"+strconv.Itoa(int(movie.ID)), nil)
			r.ServeHTTP(w, req)

			Convey("Then respond should have status 200 ok", func() {
				So(w.Code, ShouldEqual, http.StatusOK)
			})

			Convey("Then the response should contain the sucess message", func() {
				response := make(map[string]string)
				err := json.Unmarshal(w.Body.Bytes(), &response)
				So(err, ShouldBeNil)

				So(response["message"], ShouldEqual, "Movie deleted sucessfully")
			})

			Convey("then the movie should be deleted from the database", func() {
				var DeleteMovie database.Movie
				res := database.DB.First(&DeleteMovie, movie.ID)
				So(res.RowsAffected, ShouldEqual, 0)
			})
		})
		Convey("when calling DELETE /movies:id with invalid movie ID", func() {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("DELETE", "/movies/999", nil)
			r.ServeHTTP(w, req)

			Convey("then the response have status 404 not found", func() {
				So(w.Code, ShouldEqual, http.StatusNotFound)
			})
		})
		database.DB.Migrator().DropTable(&database.Movie{})
	})

}
