package ethereal

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"log"
	"os"
)

type App struct {
	Db *gorm.DB
}

var app App

func Run() {
	envLoading()
	app = App{Db: Database()}
	if len(os.Args) > 1 {
		CliRun()
	} else {
		router := gin.Default()
		v1 := router.Group("/api/v1/todos")
		{
			v1.POST("/", func(context *gin.Context) {

			})
			//v1.POST("/", CreateTodo)
			//v1.GET("/", FetchAllTodo)
			//v1.GET("/:id", FetchSingleTodo)
			//v1.PUT("/:id", UpdateTodo)
			//v1.DELETE("/:id", DeleteTodo)
		}
		router.Run()
	}

}

func envLoading() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

//func CreateTodo(c *gin.Context) {
//	completed, _ := strconv.Atoi(c.PostForm("completed"))
//	todo := Todo{Title: c.PostForm("title"), Completed: completed}
//	db := Database()
//	db.Save(&todo)
//	c.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "message": "Todo item created successfully!", "resourceId": todo.ID})
//}
