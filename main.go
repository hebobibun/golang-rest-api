package main

import (
	"api/config"
	"api/controller"
	"api/model"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	cfg := config.InitConfig()
	db := config.InitDB(*cfg)
	config.Migrate(db)
	userModel := model.UserModel{DB: db}
	userControll := controller.UserController{Mdl: userModel}
	itemModel := model.ItemModel{DB: db}
	itemControll := controller.ItemController{Mdl: itemModel}

	e.Pre(middleware.RemoveTrailingSlash()) // fungsi ini dijalankan sebelum routing
	e.Use(middleware.CORS()) // WAJIB DIPAKAI agar tidak terjadi masalah permission
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))

	e.POST("/register", userControll.Register())
	e.POST("/login", userControll.Login())
	
	auth := e.Group("/d")
	auth.Use(middleware.JWT([]byte(config.InitConfig().JWTKey)))
	
	auth.GET("/users", userControll.GetAll())
	auth.GET("/users", userControll.GetID())
	auth.PUT("/users", userControll.Update())
	auth.DELETE("/users", userControll.Delete())

	auth.POST("/items", itemControll.Insert())
	auth.GET("/items", itemControll.GetAll())
	auth.GET("/items/:id", itemControll.GetByID())
	auth.DELETE("/items/:id", itemControll.Delete())
	auth.PUT("/items/:id", itemControll.Update())

	if err := e.Start(":8000"); err != nil {
		log.Println(err.Error())
	}
}
