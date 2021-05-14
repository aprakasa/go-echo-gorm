package main

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type (
	Store struct {
		DB *gorm.DB
	}

	Projects struct {
		ID          uint   `gorm:"primaryKey" json:"id"`
		Name        string `gorm:"size:255" json:"name"`
		Description string `gorm:"size:255" json:"description"`
		CreatedAt   *time.Time
		UpdatedAt   *time.Time
	}

	Response struct {
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}
)

func (s Store) createProject(c echo.Context) error {

	project := new(Projects)
	if err := c.Bind(&project); err != nil {
		return c.JSON(http.StatusBadRequest, Response{"unable to bind data", nil})
	}

	if err := s.DB.Create(&project).Error; err != nil {
		return c.JSON(http.StatusBadRequest, Response{err.Error(), nil})
	}

	return c.JSON(http.StatusCreated, Response{"project created", project})

}

func (s Store) getProjects(c echo.Context) error {

	var projects []Projects

	if err := s.DB.Find(&projects).Error; err != nil {
		return c.JSON(http.StatusBadRequest, Response{"unable to display projects", nil})
	}

	return c.JSON(http.StatusOK, Response{"project created", projects})

}

func (s Store) getProject(c echo.Context) error {

	id, _ := strconv.Atoi(c.Param("id"))

	var project Projects

	err := s.DB.First(&project, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return c.JSON(http.StatusNotFound, Response{"project not found", nil})
	} else if err != nil {
		return c.JSON(http.StatusBadRequest, Response{"unable to display project", nil})
	}

	return c.JSON(http.StatusOK, Response{"project created", project})

}

func (s Store) updateProject(c echo.Context) error {

	id, _ := strconv.Atoi(c.Param("id"))

	project := new(Projects)
	if err := c.Bind(&project); err != nil {
		return c.JSON(http.StatusBadRequest, Response{"unable to bind data", nil})
	}

	err := s.DB.Model(&project).Where("id=?", id).Updates(&project).First(&project, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return c.JSON(http.StatusNotFound, Response{"project not found", id})
	} else if err != nil {
		return c.JSON(http.StatusBadRequest, Response{err.Error(), nil})
	}

	return c.JSON(http.StatusOK, Response{"project updated", project})

}

func (s Store) deleteProject(c echo.Context) error {

	id, _ := strconv.Atoi(c.Param("id"))

	var project Projects

	err := s.DB.Model(&project).Where("id=?", id).Delete(&project, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return c.JSON(http.StatusNotFound, Response{"project not found", id})
	} else if err != nil {
		return c.JSON(http.StatusBadRequest, Response{err.Error(), nil})
	}

	return c.NoContent(http.StatusNoContent)

}

func main() {

	db, err := gorm.Open(sqlite.Open("crud.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&Projects{})

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))

	app := Store{
		DB: db,
	}

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, Response{"service is up and running", nil})
	})

	e.POST("/projects", app.createProject)
	e.GET("/projects", app.getProjects)
	e.GET("/projects/:id", app.getProject)
	e.PATCH("/projects/:id", app.updateProject)
	e.DELETE("/projects/:id", app.deleteProject)

	e.Logger.Fatal(e.Start("localhost:3333"))

}
