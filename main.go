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

package main

import (
	"log"
	"net/http"
	"os"

	"github.com/emicklei/go-restful"
	"github.com/emicklei/go-restful-swagger12"

	"kube-service/proxy"
	_ "kube-service/config"
	"kube-service/svc"
)

func main() {

	defaultContainer := restful.NewContainer()
	svc.InitService(defaultContainer)
	defaultContainer.Filter(containerLogging)

	workDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("get current dir failed, %v", err)
	}
	// Optionally, you can install the Swagger Service which provides a nice Web UI on your REST API
	// You need to download the Swagger HTML5 assets and change the FilePath location in the config below.
	// Open http://localhost:8080/apidocs and enter http://localhost:8080/apidocs.json in the api input field.
	config := swagger.Config{
		WebServices:    defaultContainer.RegisteredWebServices(), // you control what services are visible
		WebServicesUrl: "http://localhost:8080",
		ApiPath:        "/apidocs.json",

		// Optionally, specify where the UI is located
		SwaggerPath:     "/apidocs/",
		SwaggerFilePath: workDir + "/swagger_dist"}

	swagger.RegisterSwaggerService(config, defaultContainer)
	server := &http.Server{Addr: ":8080", Handler: defaultContainer}

	log.Print("start listening on localhost:8080")
	go proxy.Run()
	log.Fatal(server.ListenAndServe())
}

func containerLogging(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	log.Printf("[container-filter] %s, %s\n", req.Request.Method, req.Request.URL)
	chain.ProcessFilter(req, resp)
}
