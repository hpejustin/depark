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

package model

type User struct {
	Id   string `json: "id"`
	Name string `json: "name"`
	/*TODO: User Auth*/
	/*Token string `json: "token"`*/
}

type Cluster struct {
	Id   string `json: "id"`
	Name string `json: "name"`
}

type Deployment struct {
	Id          string        `json: "id"`
	Name        string        `json: "name"`
	Conf        Configuration `json: "conf"`
	ClusterName Cluster       `json: "cluster"`
	Namespace   string        `json: "namespace"`
}

type Configuration struct {
	Id   string `json: "id"`
	Name string `json: "name"`
	/*resource url*/
	Source string `json: "source"`
	/*format:cluster/ns/cfg*/
	/*TODO: Support multiple kubernetes*/
	Target string `json: "source"`
}
