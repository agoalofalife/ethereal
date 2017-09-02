package ethereal

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"github.com/justinas/alice"
	"github.com/qor/i18n"
	"log"
	"net/http"
	"os"
	"path"
	"runtime"
)

var App Application
var mutations GraphQlMutations
var queries GraphQlQueries

type GraphQlMutations map[string]*graphql.Field
type GraphQlQueries map[string]*graphql.Field

func (g GraphQlMutations) Add(name string, field *graphql.Field) GraphQlMutations {

	mutations[name] = field
	return g
}

func (g GraphQlQueries) Add(name string, field *graphql.Field) {
	queries[name] = field
}

//
//func startMutations() GraphQlMutations {
//	mutations = map[string]*graphql.Field{
//		"users": &UserField,
//		"role":  &RoleField,
//	}
//
//}

// Base structure
type Application struct {
	// library gorm for work database
	Db *gorm.DB
	// localization application
	I18n            *i18n.I18n
	Middleware      *Middleware
	GraphQlMutation graphql.Fields
	GraphQlQuery    graphql.Fields
}

func Start() {
	// First we have to determine the mode of operation
	// - cli console
	// - api server
	// Secondly, we must determine the sequence of actions

	App = Application{
		Db:         ConstructorDb(),
		I18n:       ConstructorI18N(),
		Middleware: ConstructorMiddleware(),
		GraphQlQuery: graphql.Fields{
			"users": &UserField,
			"role":  &RoleField,
		},
		GraphQlMutation: graphql.Fields{
			"createUser": &createUser,
		},
	}

	App.Middleware.LoadApplication()

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

	I18nGraphQL().Fill()

	if len(os.Args) > 1 {
		CliRun()
	} else {
		h := handler.New(&handler.Config{
			Schema: &schema,
			Pretty: true,
		})

		// here can add middleware
		http.Handle("/graphql", alice.New(App.Middleware.includeMiddleware...).Then(h))

		http.HandleFunc("/auth0/login", func(w http.ResponseWriter, r *http.Request) {
			claims := EtherealClaims{
				jwt.StandardClaims{
					ExpiresAt: 15000,
					Issuer:    "test",
				},
			}
			// TODO add choose crypt via configuration!
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

			tokenString, err := token.SignedString(JWTKEY())
			fmt.Println(tokenString, err)
		})

		// Serve static files, if variable env debug in true.
		if os.Getenv("DEBAG") != "" && os.Getenv("DEBAG") == "true" {
			_, filename, _, _ := runtime.Caller(0)
			fs := http.FileServer(http.Dir(path.Dir(filename) + "/static"))
			http.Handle("/", fs)
		}

		if os.Getenv("SERVER_PORT") == "" {
			os.Setenv("SERVER_PORT", "8080")
		}
		fmt.Println("Now server is running on port " + os.Getenv("SERVER_PORT"))
		http.ListenAndServe(":"+os.Getenv("SERVER_PORT"), nil)
	}
}

/**
/ Load environment variables
*/
func envLoading() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
