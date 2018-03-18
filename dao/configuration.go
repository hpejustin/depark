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

package dao

import (
	"bytes"
	"depark/model"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type ConfigurationDAO struct{}

func NewConfigurationDAO() *ConfigurationDAO {
	return &ConfigurationDAO{}
}

func (dao *ConfigurationDAO) Update(configuration model.Configuration) error {
	source := configuration.Source
	target := configuration.Target
	filePath := download(source)
	updateConfigMap(filePath, target)
	return nil
}

func download(source string) string {
	log.Print(source)
	client := http.Client{Timeout: time.Second * 5}
	req, err := http.NewRequest(http.MethodGet, source, nil)
	req.Header.Add("Authorization", "Basic YWRtaW46amZyb2djaGluYQ==")
	if err != nil {
		log.Fatal(err)
	}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile("app.cfg", body, 0755)
	if err != nil {
		log.Fatal(err)
	}
	return "app.cfg"
}

func updateConfigMap(filePath string, target string) error {
	client := http.Client{Timeout: time.Second * 5}
	data, err := ioutil.ReadFile("template/configMap-template.json")
	if err != nil {
		log.Fatal(err)
	}
	req, err := http.NewRequest(http.MethodPut, target, bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", "Basic YWRtaW46amZyb2djaGluYQ==")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if resp.Status != 200 {
		log.Fatal(fmt.Sprintf("Update configMap %s according to %s failed", target, filePath))
	}
	return nil
}
