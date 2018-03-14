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

package filters

import (
	"depark/cache"
	"log"

	"github.com/emicklei/go-restful"
)

// ContainerLogging is a filter for container
func ContainerLogging(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	log.Printf("[container-filter] %s, %s\n", req.Request.Method, req.Request.URL)
	chain.ProcessFilter(req, resp)
}

// ServiceLogging is a filter for service
func ServiceLogging(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	log.Printf("[service-filter] %s, %s\n", req.Request.Method, req.Request.URL)
	chain.ProcessFilter(req, resp)
}

// ServiceDiagnosis provide state of backend system
func ServiceDiagnosis(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	if !cache.FettleCache.Health {
		resp.AddHeader("Content-Type", "application/json")
		resp.WriteHeaderAndEntity(203, "backend is not ready")
		return
	}
	chain.ProcessFilter(req, resp)
}
