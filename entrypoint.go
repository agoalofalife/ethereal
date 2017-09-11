package ethereal

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/jinzhu/gorm"
	"github.com/justinas/alice"
	"github.com/qor/i18n"
	"log"
	"net/http"
	"os"
	"path"
	"runtime"
)

var App Application

const runServer = "Now server is running on port "

// Base structure
type Application struct {
	Db              *gorm.DB
	I18n            *i18n.I18n
	Middleware      *Middleware
	GraphQlMutation graphql.Fields
	GraphQlQuery    graphql.Fields
	Context         context.Context
	Config          *Config
}

func Start() {
	// Config variables
	var (
		debug string = GetCnf("GRAPHQL.DEBUG").(string)
		host  string = GetCnf("HOST.PORT").(string)
	)

	App = Application{
		Db:              ConstructorDb(),
		I18n:            ConstructorI18N(),
		Middleware:      ConstructorMiddleware(),
		GraphQlQuery:    startQueries(),
		GraphQlMutation: startMutations(),
		Context:         context.Background(),
		Config:          ConstructorConfig(),
	}
	// link itself
	CtxStruct(&App, App)
	App.Middleware.LoadApplication(&App)

	//root mutation
	var rootMutation = graphql.NewObject(graphql.ObjectConfig{
		Name:   "RootMutation",
		Fields: App.GraphQlMutation,
	})

	// root query
	var rootQuery = graphql.NewObject(graphql.ObjectConfig{
		Name:   "RootQuery",
		Fields: App.GraphQlQuery,
	})

	// define schema, with our rootQuery and rootMutation
	var schema, _ = graphql.NewSchema(graphql.SchemaConfig{
		Query:    rootQuery,
		Mutation: rootMutation,
	})
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		opts := handler.NewRequestOptions(r)
		result := graphql.Do(graphql.Params{
			Schema:         schema,
			OperationName:  opts.OperationName,
			VariableValues: opts.Variables,
			RequestString:  opts.Query,
			Context:        App.Context,
		})
		if len(result.Errors) > 0 {
			log.Printf("wrong result, unexpected errors: %v", result.Errors)
			return
		}
		json.NewEncoder(w).Encode(result)
	})

	// here can add middleware
	http.Handle("/graphql", alice.New(App.Middleware.includeMiddleware...).Then(h))

	if debug == "true" {
		_, filename, _, _ := runtime.Caller(0)
		fs := http.FileServer(http.Dir(path.Dir(filename) + "/static"))
		http.Handle("/", fs)
	}

	fmt.Println(runServer + host)
	http.ListenAndServe(":"+host, nil)

}
