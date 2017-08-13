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
	"kube-service/model"
	"time"
	"log"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
)

type UserDAO struct {}

func NewUserDAO() *UserDAO  {
	return &UserDAO {}
}

func (dao UserDAO) List() map[string]model.User {
	var result map[string]model.User
	list := make(chan map[string]model.User)
	timeout := make(chan string)
	go func() {
		time.Sleep(time.Second * 5)
		list <- query()
	}()

	go func() {
		time.Sleep(time.Second * 15)
		timeout <- "TIME_OUT"
	}()

	select {
	case result = <- list:
		log.Printf("result: %#v", result)
	case <- timeout:
		log.Println("request timeout")
	}
	return result
}

func (dao UserDAO) Find(id string) model.User {
	var result model.User
	userChan := make(chan model.User)
	timeout := make(chan string)
	go func() {
		time.Sleep(time.Second * 5)
		userChan <- getOne(id)
	}()

	go func() {
		time.Sleep(time.Second * 15)
		timeout <- "TIME_OUT"
	}()
	select {
	case result = <- userChan:
		log.Printf("result: %#v", result)
	case <- timeout:
		log.Println("request timeout")
	}
	return result
}

func (dao UserDAO) Add(user model.User) model.User {
	var result model.User
	userChan := make(chan model.User)
	timeout := make(chan string)
	go func() {
		time.Sleep(time.Second * 5)
		userChan <- create(user)
	}()

	go func() {
		time.Sleep(time.Second * 15)
		timeout <- "TIME_OUT"
	}()
	select {
	case result = <- userChan:
		log.Printf("result: %#v", result)
	case <- timeout:
		log.Println("request timeout")
	}
	return result
}

func (dao UserDAO) Remove(id string) bool {
	var result bool
	userChan := make(chan bool)
	timeout := make(chan string)
	go func() {
		time.Sleep(time.Second * 5)
		userChan <- delete(id)
	}()

	go func() {
		time.Sleep(time.Second * 15)
		timeout <- "TIME_OUT"
	}()
	select {
	case result = <- userChan:
		log.Printf("result: %#v", result)
	case <- timeout:
		log.Println("request timeout")
	}
	return result
}

func (dao UserDAO) Update(user model.User, id string) model.User {
	var result model.User
	userChan := make(chan model.User)
	timeout := make(chan string)
	go func() {
		time.Sleep(time.Second * 5)
		userChan <- change(user, id)
	}()

	go func() {
		time.Sleep(time.Second * 15)
		timeout <- "TIME_OUT"
	}()
	select {
	case result = <- userChan:
		log.Printf("result: %#v", result)
	case <- timeout:
		log.Println("request timeout")
	}
	return result
}

func change(user model.User, id string) model.User {
	db, err := con()
	CheckErr(err)
	defer db.Close()
	stmt, err := db.Prepare("update account set name=? where id=?")
	CheckErr(err)
	res, err := stmt.Exec(user.Name, id)
	CheckErr(err)
	_, err = res.RowsAffected()
	CheckErr(err)
	return user
}

func delete(id string) bool {
	db, err := con()
	CheckErr(err)
	defer db.Close()
	stmt, err := db.Prepare("delete from account where id=?")
	CheckErr(err)
	res, err := stmt.Exec(id)
	CheckErr(err)
	n, err := res.RowsAffected()
	CheckErr(err)
	return n == int64(1)

}

func create(user model.User) model.User {
	db, err := con()
	CheckErr(err)
	defer db.Close()

	//stmt, err := db.Prepare(fmt.Sprintf("insert into account values('%s','%s');", user.Id, user.Name))
	stmt, err := db.Prepare("insert account set id=?, name=?")
	CheckErr(err)
	res, err := stmt.Exec(user.Id, user.Name)
	CheckErr(err)
	n, err := res.RowsAffected()
	CheckErr(err)
	if n != int64(1) {
		panic("create fail")
	}
	return user
}

func getOne(id string) model.User {
	db, err := con()
	CheckErr(err)
	defer db.Close()
	user := model.User{}
	rows, err := db.Query(fmt.Sprintf("select * from account where id = %s;", id))
	CheckErr(err)
	for rows.Next() {
		var id, name string
		err = rows.Scan(&id, &name)
		if err != nil {
			log.Println("scan failed.")
		}
		user.Id = id
		user.Name = name
	}
	return user
}

func query() (userMap map[string]model.User) {
	db, err := con()
	CheckErr(err)
	defer db.Close()
	userMap =  make(map[string]model.User)
	rows, err := db.Query("select * from account;")
	CheckErr(err)
	for rows.Next() {
		var id, name string
		err = rows.Scan(&id, &name)
		if err != nil {
			log.Println("scan failed.")
		}
		userMap[id] = model.User{Id: id, Name: name}
	}
	return
}

func con() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:1234@tcp(127.0.0.1:3306)/demo?charset=utf8")
	if err != nil {
		return nil, err
	}
	return db, nil
}

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}


