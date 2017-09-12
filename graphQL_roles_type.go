package ethereal

import (
	"github.com/graphql-go/graphql"
	"fmt"
)

var roleType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Role",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type:        graphql.String,
			Description: string(ConstructorI18N().T(GetCnf("L18N.LOCALE").(string), "graphQL.RoleType.id")),
		},
		"name": &graphql.Field{
			Type:        graphql.String,
			Description: string(ConstructorI18N().T(GetCnf("L18N.LOCALE").(string), "graphQL.RoleType.name")),
		},
		"display_name": &graphql.Field{
			Type:        graphql.String,
			Description: string(ConstructorI18N().T(GetCnf("L18N.LOCALE").(string), "graphQL.RoleType.display_name")),
		},
		"description": &graphql.Field{
			Type:        graphql.String,
			Description: string(ConstructorI18N().T(GetCnf("L18N.LOCALE").(string), "graphQL.RoleType.description")),
		},
	},
})

var RoleField = graphql.Field{
	Type:        graphql.NewList(roleType),
	Description: string(ConstructorI18N().T(GetCnf("L18N.LOCALE").(string), "graphQL.Role.Description")),
	Args: graphql.FieldConfigArgument{
		"id": &graphql.ArgumentConfig{
			Type: graphql.String,
		},
	},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			var roles []Role
			App.Db.Find(&roles)

			idQuery, isOK := params.Args["id"].(string)
			if isOK {
				for _, role := range roles {
					if string(role.ID) == idQuery {
						return role, nil
					}
				}
			}
		fmt.Println(roles[1].ID, roles[1].Name)
			return roles, nil
	},
}
