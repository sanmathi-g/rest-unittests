package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sana/rest/database"
	"gorm.io/gorm"
)

type MovieService interface {
	Create(value interface{}) (tx *gorm.DB)
	Find(dest interface{}, conds ...interface{}) (tx *gorm.DB)
	Model(value interface{}) (tx *gorm.DB)
	Delete(value interface{}, conds ...interface{}) (tx *gorm.DB)
}

func CreateMovie(c *gin.Context) {
	var movie *database.Movie
	err := c.ShouldBind(&movie)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	res := database.DB.Create(movie)
	if res.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "error creating a movie",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"movie": movie,
	})
	return
}

func ReadMovie(c *gin.Context) {
	var movie database.Movie
	id := c.Param("id")
	res := database.DB.Find(&movie, id)
	if res.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "movie not found",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"movie": movie,
	})
	return
}

func ReadMovies(c *gin.Context) {
	var movies []database.Movie
	res := database.DB.Find(&movies)
	if res.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": errors.New(" not found"),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"books": movies,
	})
	return
}
func UpdateMovie(c *gin.Context) {
	var movie database.Movie
	id := c.Param("id")
	err := c.ShouldBind(&movie)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	var updateBook database.Movie
	res := database.DB.Model(&updateBook).Where("id = ?", id).Updates(movie)

	if res.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "movie not updated",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"movie": movie,
	})
	return
}

func DeleteMovie(c *gin.Context) {
	var movie database.Movie
	id := c.Param("id")
	res := database.DB.Find(&movie, id)
	if res.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "movie not found",
		})
		return
	}
	database.DB.Delete(&movie)
	c.JSON(http.StatusOK, gin.H{
		"message": "movie deleted successfully",
	})
	return
}
