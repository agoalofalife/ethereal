package ethereal

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ethereal-go/ethereal/root/app"
	"github.com/ethereal-go/ethereal/root/middleware"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"log"
	"net/http"
	"path"
	"runtime"
)

var App app.Application
var Middleware *middleware.Middleware

const runServer = "Now server is running on port "

func Start() {
	// Config variables
	var (
		debug string = GetCnf("GRAPHQL.DEBUG").(string)
		host  string = GetCnf("HOST.PORT").(string)
	)

	App := app.Application{
		Db:              ConstructorDb(),
		I18n:            ConstructorI18N(),
		GraphQlQuery:    startQueries(),
		GraphQlMutation: startMutations(),
		Context:         context.Background(),
		Config:          ConstructorConfig(),
	}
	// link itself
	App.Context = CtxStruct(&App, &App)
	mid := ConstructorMiddleware()
	mid.LoadApplication(&App)

	// root query
	var rootQuery = graphql.NewObject(graphql.ObjectConfig{
		Name:   "RootQuery",
		Fields: App.GraphQlQuery,
	})

	//root mutation
	var rootMutation = graphql.NewObject(graphql.ObjectConfig{
		Name:   "RootMutation",
		Fields: App.GraphQlMutation,
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
	http.Handle("/graphql", mid.GetHandler(h))

	if debug == "true" {
		_, filename, _, _ := runtime.Caller(0)
		fs := http.FileServer(http.Dir(path.Dir(filename) + "/static"))
		http.Handle("/", fs)
	}

	fmt.Println(runServer + host)
	err := http.ListenAndServe(":"+host, nil)
	if err != nil {
		fmt.Println(err.Error())
	}
}
