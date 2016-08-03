package main

import (
	// standard library packages
	"net/http"

	// third party http package
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	mgo "gopkg.in/mgo.v2"

	// my local controllers directory
	"github.com/bradford-hamilton/go-rest-service/controllers"
)

func getSession() *mgo.Session {
	// connect to local mongo db
	session, err := mgo.Dial("mongodb://localhost")
	// check for connection error
	if err != nil {
		panic(err)
	}
	return session
}

func main() {
	// instantiate new router
	router := httprouter.New()

	// get an instance of UserController
	uc := controllers.NewUserController(getSession())

	router.GET("/user", uc.GetAllUsers)

	router.GET("/user/:id", uc.GetUser)

	router.POST("/user", uc.CreateUser)

	router.DELETE("/user/:id", uc.RemoveUser)

	// cors
	c := cors.New(cors.Options{
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
	})
	handler := c.Handler(router)
	// listen on 1337
	http.ListenAndServe("localhost:1337", handler)

}
