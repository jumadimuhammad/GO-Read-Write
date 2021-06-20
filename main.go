package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

type Handler struct {
	Bio []Bio
	Api ApiResponse
}

type Bio struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
}

type ApiResponse struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func main() {
	var handler Handler

	mux := http.NewServeMux()
	mux.HandleFunc("/user", handler.ListFunc)
	mux.HandleFunc("/user/add", handler.AddFunc)
	mux.HandleFunc("/user/delete", handler.DeleteFunc)

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	fmt.Println("Server running on port ", server.Addr)

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}

}

func (b *Handler) AddFunc(res http.ResponseWriter, req *http.Request) {
	method := req.Method
	name := req.FormValue("name")
	phone := req.FormValue("phone")
	address := req.FormValue("address")

	if method != "POST" {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte("Error method invalid"))
		return
	}

	if name == "" || phone == "" || address == "" {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte("Error data invalid"))
		return
	}

	b.List()
	b.Add(name, phone, address)
	b.Save()

	b.Bio[len(b.Bio)-1].Id = len(b.Bio)

	b.Api.Code = http.StatusCreated
	b.Api.Data = b.Bio[len(b.Bio)-1]
	b.Api.Message = "Success"

	mar, _ := json.Marshal(b.Api)
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(mar))
}

func (b *Handler) ListFunc(res http.ResponseWriter, req *http.Request) {
	method := req.Method

	if method != "GET" {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte("Error method invalid"))
		return
	}

	b.List()

	b.Api.Code = http.StatusOK
	b.Api.Data = b.Bio
	b.Api.Message = "Success"

	mar, _ := json.Marshal(b.Api)
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(mar))
}

func (b *Handler) DeleteFunc(res http.ResponseWriter, req *http.Request) {
	id, _ := strconv.Atoi(req.FormValue("id"))
	method := req.Method

	if method != "DELETE" {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte("Error method invalid"))
		return
	}

	if id == 0 {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte("Error data invalidate"))
		return
	}

	b.List()

	if len(b.Bio) < id {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte("Error id not found"))
		return
	}

	b.Delete(id)
	b.Save()

	b.Api.Code = http.StatusOK
	b.Api.Data = b.Bio
	b.Api.Message = "Success"

	mar, _ := json.Marshal(b.Api)
	res.Write([]byte(mar))
}

func (b *Handler) List() {
	file, _ := ioutil.ReadFile("data/files.json")
	_ = json.Unmarshal(file, &b.Bio)

	for i := range b.Bio {
		b.Bio[i].Id = i + 1
	}
}

func (b *Handler) Delete(id int) {
	if len(b.Bio) == id {
		b.Bio = b.Bio[0 : id-1]
	} else {
		a := b.Bio[0 : id-1]
		c := b.Bio[id:len(b.Bio)]

		a = append(a, c...)
		b.Bio = a
	}
}

func (b *Handler) Add(name, phone, address string) {
	b.Bio = append(b.Bio, Bio{
		Id:      len(b.Bio)+1,
		Name:    name,
		Phone:   phone,
		Address: address,
	})
}

func (b *Handler) Save() {
	mar, _ := json.Marshal(b.Bio)
	ioutil.WriteFile("data/files.json", mar, 0777)
	b.List()
}
