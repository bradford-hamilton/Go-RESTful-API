package controllers

import (
	// standard library packages
	"encoding/json"
	"fmt"
	"net/http"
	// third party packages
	"github.com/julienschmidt/httprouter"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	// my local models
	"github.com/bradford-hamilton/go-rest-service/models"
)

type (
	// UserController represents the controller for operating on the User resource
	UserController struct {
		session *mgo.Session
	}
)

func NewUserController(session *mgo.Session) *UserController {
	return &UserController{session}
}

// GET all users
func (uc UserController) GetAllUsers(res http.ResponseWriter, req *http.Request, p httprouter.Params) {

	var results []models.User

	if err := uc.session.DB("go-rest-service").C("users").Find(nil).All(&results); err != nil {
		res.WriteHeader(404)
		return
	}
	fmt.Println(results)

	// write content-type, status code payload
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(200)
	json.NewEncoder(res).Encode(results)
}

// GET user
func (uc UserController) GetUser(res http.ResponseWriter, req *http.Request, p httprouter.Params) {
	// get id
	id := p.ByName("id")
	// verify id is ObjectId, otherwise return
	if !bson.IsObjectIdHex(id) {
		res.WriteHeader(404)
		return
	}
	// get id
	oid := bson.ObjectIdHex(id)

	u := models.User{}
	// fetch user
	if err := uc.session.DB("go-rest-service").C("users").FindId(oid).One(&u); err != nil {
		res.WriteHeader(404)
		return
	}
	// write content-type, status code, payload
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(200)
	json.NewEncoder(res).Encode(u)
}

// POST a user
func (uc UserController) CreateUser(res http.ResponseWriter, req *http.Request, p httprouter.Params) {
	// Stub an user to be populated from the body
	u := models.User{}
	// Populate the user data
	json.NewDecoder(req.Body).Decode(&u)
	// Add an Id
	u.ID = bson.NewObjectId()
	// write the user to mongodb
	uc.session.DB("go-rest-service").C("users").Insert(u)
	// Write content-type, statuscode, payload
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(201)
	json.NewEncoder(res).Encode(u)
}

// RemoveUser removes an existing user resource
func (uc UserController) RemoveUser(res http.ResponseWriter, req *http.Request, p httprouter.Params) {
	// get specific user
	id := p.ByName("id")
	// verify ObjectId otherwise return
	if !bson.IsObjectIdHex(id) {
		res.WriteHeader(404)
		return
	}
	// get user
	oid := bson.ObjectIdHex(id)
	// remove user
	if err := uc.session.DB("go-rest-service").C("users").RemoveId(oid); err != nil {
		res.WriteHeader(404)
		return
	}

	response := make(map[string]string)
	response["message"] = "User was successfully deleted!"
	// Write content-type, statuscode, payload
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(200)
	json.NewEncoder(res).Encode(response)
}
