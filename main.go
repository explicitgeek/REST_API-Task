package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type User struct {
	Email string
	Name  string
	Phone string
	Id    uint64
}

var nextUserID uint64

var Users = []User{
	{
		Email: "bharat@appointy.com", Name: "Bharat", Phone: "6377187695", Id: 1,
	},
	{
		Email: "ram@google.com", Name: "Ram", Phone: "6549871560", Id: 2,
	},
	{
		Email: "felix@google.com", Name: "Felix", Phone: "6599789126", Id: 3,
	},
}

func main() {

	http.HandleFunc("/users", func(res http.ResponseWriter, req *http.Request) {
		if req.Method == "POST" {
			res.Header().Set("Content-Type", "application/json")
			next_User := User{}
			json.NewDecoder(req.Body)
			nextUserID++
			next_User.Id = nextUserID
			Users = append(Users, next_User)

			new_User := User{Email: next_User.Email, Name: next_User.Name, Phone: next_User.Phone, Id: next_User.Id}
			json.NewEncoder(res).Encode(new_User)

		} else if req.Method == "GET" {
			res.Header().Set("Content-Type", "application/json")
			json.NewEncoder(res).Encode(Users)
		}
	})

	http.HandleFunc("/users/{id}", func(res http.ResponseWriter, req *http.Request) {
		if req.Method == "GET" {
			res.Header().Set("Content-Type", "application/json")

			fields := strings.Split(req.URL.String(), "/")
			id, err := strconv.ParseUint(fields[len(fields)-1], 10, 64)

			if err != nil {
				res.WriteHeader(http.StatusBadRequest)
			}

			fmt.Printf("Request to get user %v", id)

			flag := false
			for _, u := range Users {
				if u.Id == id {
					json.NewEncoder(res).Encode(u)
					flag = true
					break
				}
			}

			if !flag {
				res.WriteHeader(http.StatusBadRequest)
			}

		} else if req.Method == "PUT" {
			res.Header().Set("Content-Type", "application/json")

			json.NewDecoder(req.Body)
		} else if req.Method == "DELETE" {
			// get the user ID from the path
			fields := strings.Split(req.URL.String(), "/")
			id, err := strconv.ParseUint(fields[len(fields)-1], 10, 64)
			if err != nil {
				res.WriteHeader(http.StatusBadRequest)
			}

			log.Printf("Request to delete user %v", id)

			var tmp = []User{}
			for _, u := range Users {
				if id == u.Id {
					continue
				}
				tmp = append(tmp, u)
			}
			Users = tmp
		} else {
			res.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(res, "Unsupported method '%v' to %v\n", req.Method, req.URL)
			log.Printf("Unsupported method '%v' to %v\n", req.Method, req.URL)
		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
