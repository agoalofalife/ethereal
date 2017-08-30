package ethereal

import (
	"github.com/gin-gonic/gin"
	"github.com/qor/i18n"
	"github.com/qor/i18n/backends/database"
	"net/http"
	"os"
)

//type App struct {
//	Db   *gorm.DB
//	I18n *i18n.I18n
//}

func Run() {
	envLoading()
	db := Database()
	I18n := i18n.New(
		database.New(db),
	)
	app = App{Db: Database(), I18n: I18n}

	if len(os.Args) > 1 {
		CliRun()
	} else {
		router := gin.Default()
		v1 := router.Group("/api/v1/users")
		{
			v1.GET("/", FetchAllUsers)
			//v1.POST("/", CreateTodo)
			//v1.GET("/:id", FetchSingleTodo)
			//v1.PUT("/:id", UpdateTodo)
			//v1.DELETE("/:id", DeleteTodo)
		}
		router.Run()
	}

}

func FetchAllUsers(c *gin.Context) {
	var users []User
	app.Db.Find(&users)
	if len(users) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No todo found!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": users})
}
