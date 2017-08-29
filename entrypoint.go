package ethereal

import (
	"fmt"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/jinzhu/gorm"
	"github.com/qor/i18n"
	"github.com/qor/i18n/backends/database"
	"net/http"
	"os"
	"path"
	"runtime"
)

var app App

type App struct {
	Db   *gorm.DB
	I18n *i18n.I18n
}

//root mutation
var rootMutation = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootMutation",
	Fields: graphql.Fields{
		"createUser": &createUser,
	},
})

// root query
// we just define a trivial example here, since root query is required.
var rootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootQuery",
	Fields: graphql.Fields{
		"users": &UserField,
		"role":  &RoleField,
	},
})

// define schema, with our rootQuery and rootMutation
var schema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query:    rootQuery,
	Mutation: rootMutation,
})

func ConstructorDb() *gorm.DB {
	if app.Db == nil {
		envLoading()
		app.Db = Database()
	}
	return app.Db

}
func ConstructorI18N() *i18n.I18n {
	if app.I18n == nil {
		app.I18n = i18n.New(
			database.New(ConstructorDb()),
		)
	}
	return app.I18n
}

func Start() {
	//envLoading()
	//db := Database()
	//I18n := i18n.New(
	//	database.New(db),
	//)

	app = App{Db: ConstructorDb(), I18n: ConstructorI18N()}
	I18nGraphQL().Fill()
	if len(os.Args) > 1 {
		CliRun()
	} else {
		//http.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		//	result := executeQuery(r.URL.Query().Get("query"), schema)
		//	json.NewEncoder(w).Encode(result)
		//})
		h := handler.New(&handler.Config{
			Schema: &schema,
			Pretty: true,
		})
		http.Handle("/graphql", middlewareLocal(h))
		// Serve static files
		_, filename, _, _ := runtime.Caller(0)
		fs := http.FileServer(http.Dir(path.Dir(filename) + "/static"))
		http.Handle("/", fs)
		fmt.Println("Now server is running on port 8080")

		//fmt.Println("Create new todo: curl -g 'http://localhost:8080/graphql?query=mutation+_{createTodo(text:\"My+new+todo\"){id,text,done}}'")
		//fmt.Println("Update todo: curl -g 'http://localhost:8080/graphql?query=mutation+_{updateTodo(id:\"a\",done:true){id,text,done}}'")

		http.ListenAndServe(":8080", nil)
	}

}
