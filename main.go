package main

import (
	"fmt"
	"time"
	"errors"
	"net/http"
	"encoding/json"

	"github.com/graphql-go/graphql"
)

type User struct {
	ID		int 	`json:"id"`
	Name	string	`json:"name"`
	Gender	string	`json:"gender"`
}

type Event struct {
	ID			int 		`json:"id"`
	UserID		int			`json:"user_id"`
	Name		string		`json:"name"`
	Description	string		`json:"description"`
	StartedAt	time.Time	`json:"started_at"`
}

var users = []User{
	{
		ID: 1,
		Name: "Dicky",
		Gender: "Male",
	},
	{
		ID: 2,
		Name: "Josephine",
		Gender: "Female",
	},
}

var events = []Event{
	{
		ID: 1,
		UserID: 1,
		Name: "Harbolnas 2020",
		Description: "Event jualan online terbesar 2020!",
		StartedAt: time.Now(),
	},
}

var userType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"gender": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

var eventType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Event",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
			},
			"user": &graphql.Field{
				Type: userType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					event, ok := p.Source.(Event)
					if ok {
						for _, user := range users {
							if event.UserID == user.ID {
								return user, nil
							}
						}
					}
					return nil, errors.New("user not exists.")
				},
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"description": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

var queryType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			// get user by id
			// url : http://localhost:8080/api?query={user(id:1){id,name}}
			"user": &graphql.Field{
				Type: userType,
				Description: "Get user by Id",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig {
						Type: graphql.Int,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id, ok := p.Args["id"].(int)
					if ok {
						for _, user := range users {
							if user.ID == id {
								return user, nil
							}
						}
					}
					return nil, errors.New("user not exists.")
				},
			},

			// get event by id
			// url : http://localhost:8080/api?query={event(id:1){id,user{name,gender},name,description}}
			"event": &graphql.Field{
				Type: eventType,
				Description: "Get event by Id",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id, ok := p.Args["id"].(int)
					if ok {
						for _, event := range events {
							if id == event.ID {
								return event, nil
							}
						}
					}
					return nil, errors.New("event not exists.")
				},
			},

			// get user list
			// url : http://localhost:8080/api?query={users{id,name}}
			"users": &graphql.Field{
				Type: graphql.NewList(userType),
				Description: "Get all users",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return users, nil
				},
			},
		},
	},
)

var mutationType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			// create new user
			// http://localhost:8080/api?query=mutation+_{createUser(name:"Jose",gender:"Female"){id,name,gender}}
			"createUser": &graphql.Field{
				Type: userType,
				Description: "insert a new user",
				Args: graphql.FieldConfigArgument{
					"name": &graphql.ArgumentConfig {
						Type: graphql.String,
					},
					"gender": &graphql.ArgumentConfig {
						Type: graphql.String,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					user := User{
						ID: users[len(users)-1].ID + 1, // auto increment id
						Name: p.Args["name"].(string),
						Gender: p.Args["gender"].(string),
					}
					users = append(users, user)

					return user, nil
				},
			},
		},
	},
)

var schema, _ = graphql.NewSchema(
	graphql.SchemaConfig{
		Query: queryType,
		Mutation: mutationType,
	},
)

func executeQuery(query string, schema graphql.Schema) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema: schema,
		RequestString: query,
	})

	if len(result.Errors) > 0 {
		fmt.Printf("errors: %v", result.Errors)
	}

	return result
}

func main() {
	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		result := executeQuery(r.URL.Query().Get("query"), schema)
		json.NewEncoder(w).Encode(result)
	})

	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)

}