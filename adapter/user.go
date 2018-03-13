/*
Copyright 2017 The Depark Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package adapter

import (
	"net/http"
	"sync"

	"github.com/emicklei/go-restful"

	"depark/model"
	"depark/filters"
	"depark/dao"
)

type UserService struct {
	UserRepository *dao.UserDAO
	// This mutex guards all fields within this cache struct.
	mu sync.Mutex
}

func (u UserService) Register() *restful.WebService {

	ws := new(restful.WebService)
	ws.
		Path("/users").
		Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_XML) // you can specify this per route as well

	ws.Filter(filters.ServiceLogging).Filter(filters.ServiceDiagnosis)

	ws.Route(ws.GET("/").To(u.findAllUsers).
		// docs
		Doc("get all users").
		Operation("findAllUsers").
		Writes([]model.User{}).
		Returns(200, "OK", nil))

	ws.Route(ws.GET("/{user-id}").To(u.findUser).
		// docs
		Doc("get a user").
		Operation("findUser").
		Param(ws.PathParameter("user-id", "identifier of the user").DataType("string")).
		Writes(model.User{}). // on the response
		Returns(404, "Not Found", nil))

	ws.Route(ws.PUT("/{user-id}").To(u.updateUser).
		// docs
		Doc("update a user").
		Operation("updateUser").
		Param(ws.PathParameter("user-id", "identifier of the user").DataType("string")).
		Reads(model.User{})) // from the request

	ws.Route(ws.PUT("").To(u.createUser).
		// docs
		Doc("create a user").
		Operation("createUser").
		Reads(model.User{})) // from the request

	ws.Route(ws.DELETE("/{user-id}").To(u.removeUser).
		// docs
		Doc("delete a user").
		Operation("removeUser").
		Param(ws.PathParameter("user-id", "identifier of the user").DataType("string")))

	return ws
}



// GET http://localhost:8080/users
//
func (u UserService) findAllUsers(request *restful.Request, response *restful.Response) {
	list := []model.User{}
	for _, each := range u.UserRepository.List() {
		list = append(list, each)
	}
	response.WriteEntity(list)
}

// GET http://localhost:8080/users/1
//
func (u UserService) findUser(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("user-id")
	usr := u.UserRepository.Find(id)
	if len(usr.Id) == 0 {
		response.WriteErrorString(http.StatusNotFound, "User could not be found.")
	} else {
		response.WriteEntity(usr)
	}
}

// PUT http://localhost:8080/users/1
//
func (u *UserService) updateUser(request *restful.Request, response *restful.Response) {
	u.mu.Lock()
	defer u.mu.Unlock()
	id := request.PathParameter("user-id")
	user := model.User{}
	err := request.ReadEntity(&user)
	if err == nil {
		u.UserRepository.Update(user, id)
		response.WriteEntity(user)
	} else {
		response.WriteError(http.StatusInternalServerError, err)
	}
}

// PUT http://localhost:8080/users/
//
func (u *UserService) createUser(request *restful.Request, response *restful.Response) {
	user := model.User{}
	err := request.ReadEntity(&user)
	if err == nil {
		u.UserRepository.Add(user)
		response.WriteHeaderAndEntity(http.StatusCreated, user)
	} else {
		response.WriteError(http.StatusInternalServerError, err)
	}
}

// DELETE http://localhost:8080/users/1
//
func (u *UserService) removeUser(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("user-id")
	u.UserRepository.Remove(id)
}
