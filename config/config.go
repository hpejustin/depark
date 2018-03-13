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

package config

import (
	"log"
	"io/ioutil"
	"fmt"
	"encoding/json"
	"os"
	"net/http"

	"depark/cache"
	"time"
)

var options = RunOptions{}

type RunOptions struct {
	Backend string `json: "backend"`
}

func init() {
	workDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("get current dir failed, %v", err)
	}
	bytes, err := ioutil.ReadFile(fmt.Sprintf(workDir + "/config/cfg.json"))
	if err != nil {
		fmt.Printf("read file error, %v", err)
		return
	}

	if err := json.Unmarshal(bytes, &options); err != nil {
		fmt.Printf("unmarshal error, %v", err)
		return
	}

	log.Printf("backend: %s", options.Backend)
	go Run()
	log.Println("init")
}

func Run() {
	for {
		resp, err := http.Get("http://" + options.Backend)
		if err == nil && resp.StatusCode == 200 {
			cache.FettleCache.Health = true
		}
		log.Printf("[%v] backend health status is %v.", time.Now(), cache.FettleCache.Health)
		time.Sleep(time.Second * 15)
	}
}