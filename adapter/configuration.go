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
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/emicklei/go-restful"

	"depark/dao"
	"depark/filters"
	"depark/model"
)

type ConfigurationService struct {
	ConfigurationRepository *dao.ConfigurationDAO
	// This mutex guards all fields within this cache struct.
	dataSource []model.Configuration
	mu         sync.Mutex
}

func (svc ConfigurationService) Register() *restful.WebService {

	ws := new(restful.WebService)
	ws.
		Path("/configurations").
		Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_XML) // you can specify this per route as well

	ws.Filter(filters.ServiceLogging).Filter(filters.ServiceDiagnosis)

	ws.Route(ws.GET("/").To(svc.findAll).
		// docs
		Doc("get all").
		Operation("findAll").
		Writes([]model.Configuration{}).
		Returns(200, "OK", nil))

	ws.Route(ws.GET("/{id}").To(svc.find).
		// docs
		Doc("get one").
		Operation("find").
		Param(ws.PathParameter("id", "identifier").DataType("string")).
		Writes(model.Configuration{}). // on the response
		Returns(404, "Not Found", nil))

	ws.Route(ws.PUT("/{id}").To(svc.update).
		// docs
		Doc("update").
		Operation("update").
		Param(ws.PathParameter("id", "identifier").DataType("string")).
		Reads(model.Configuration{})) // from the request

	ws.Route(ws.PUT("").To(svc.create).
		// docs
		Doc("create").
		Operation("create").
		Reads(model.Configuration{})) // from the request

	ws.Route(ws.DELETE("/{id}").To(svc.remove).
		// docs
		Doc("delete").
		Operation("remove").
		Param(ws.PathParameter("id", "identifier").DataType("string")))

	svc.dataSource = []model.Configuration{
		{Id: "1", Name: "Demo1", Source: "http://demo.jfrogchina.com/artifactory/kube-config/1.0/app.cfg", Target: "http://39.106.21.94:8080/api/v1/namespaces/devops/configmaps/app-config"},
		{Id: "2", Name: "Demo2", Source: "http://demo.jfrogchina.com/artifactory/kube-config/1.0/app.cfg", Target: "http://39.106.21.94:8080/api/v1/namespaces/devops/configmaps/app-config"},
	}

	return ws
}

func (svc *ConfigurationService) findAll(request *restful.Request, response *restful.Response) {
	// TODO
	//configurationList := []model.Configuration{
	//	{Id: "1", Name: "Demo1", Source: "http://demo.jfrogchina.com/artifactory/kube-config/app.cfg", Target: "kube-dev/default/kube-conf-dev"},
	//	{Id: "2", Name: "Demo2", Source: "http://demo.jfrogchina.com/artifactory/kube-config/app.cfg", Target: "kube-prd/default/kube-conf-prd"},
	//}
	response.WriteEntity(svc.dataSource)
}

func (svc *ConfigurationService) find(request *restful.Request, response *restful.Response) {
	// TODO
}

func (svc *ConfigurationService) update(request *restful.Request, response *restful.Response) {
	svc.mu.Lock()
	defer svc.mu.Unlock()
	id := request.PathParameter("id")
	log.Printf("configuration of id: %s will be changed.\n", id)
	configuration := model.Configuration{}
	err := request.ReadEntity(&configuration)
	if err == nil {
		for _, each := range svc.dataSource {
			if each.Id == id {
				each.Source = configuration.Source
				each.Target = configuration.Target
			}
		}
		// Update configuration to Kubernetes according to the new configuration
		err = svc.ConfigurationRepository.Update(configuration)
		if err == nil {
			response.WriteEntity(configuration)
		} else {
			fmt.Sprintf("Update configuration to Kuberntes failed, err %s", err)
		}
	} else {
		response.WriteError(http.StatusInternalServerError, err)
	}
}

func (svc *ConfigurationService) create(request *restful.Request, response *restful.Response) {
	// TODO
}

func (svc *ConfigurationService) remove(request *restful.Request, response *restful.Response) {
	// TODO
}
