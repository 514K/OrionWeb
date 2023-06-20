package main

import (
	"fmt"
	"net/http"
	"text/template"
)

// var responses = make([]*RSVP, 0, 10)
var templates = make(map[string]*template.Template, 3)

type Menu struct {
	Id     string
	Name   string
	Select bool
}

type Table struct {
	Id     string
	Name   string
	Select bool
}

type Field struct {
	Header []string
	Data   [][]string
}

type Resp struct {
	Username string
	Menus    []Menu
	Tables   []Table
	Fields   []Field
}

func loadTemplate() {
	templateNames := [...]string{"web/auth", "web/menu"}

	for i, item := range templateNames {
		t, err := template.ParseFiles("web/layout.html", item+".html")
		if err == nil {
			templates[item] = t
			fmt.Println("Loaded template ", i, item)
		} else {
			panic(err)
		}
	}
}

func welcomeHandler(writer http.ResponseWriter, request *http.Request) {

	if request.Method == http.MethodPost {

		// ТУТ ДОБАВИТЬ ПРОВЕРКУ ЛОГИНА/ПАРОЛЯ
		request.ParseForm()
		// fmt.Printf("%v\n", request.Form["login"])
		// fmt.Printf("%v\n", request.Form["password"])

		// templates["web/menu"].Execute(writer, request.Form["login"][0])

		cooka := http.Cookie{Name: "username", Value: request.Form["login"][0]}
		http.SetCookie(writer, &cooka)
		http.Redirect(writer, request, "/menu", http.StatusSeeOther)
	} else {
		templates["web/auth"].Execute(writer, nil)
	}
}

func menuHandler(writer http.ResponseWriter, request *http.Request) {
	// request.ParseForm()
	username, _ := request.Cookie("username")
	// fmt.Printf("%v\n", username.Value)
	// fmt.Printf("%v\n", request.Form["password"])

	resp := Resp{}
	resp.Username = username.Value
	resp.Menus = []Menu{}

	if request.Method == http.MethodGet {
		println(request.URL.Query().Get("id"))

		id := request.URL.Query().Get("id")
		// Тут запрос на меню
		resp.Menus = append(resp.Menus, Menu{Id: "1", Name: "Настройки организации", Select: id == "1"})
		resp.Menus = append(resp.Menus, Menu{Id: "2", Name: "Настройки пользователей", Select: id == "2"})

		tableid := request.URL.Query().Get("tableid")
		// Тут запрос на таблицы
		resp.Tables = []Table{}
		if id == "1" {
			resp.Tables = append(resp.Tables, Table{Id: "1", Name: "Сотрудники", Select: tableid == "1"})
			resp.Tables = append(resp.Tables, Table{Id: "2", Name: "Подразделения", Select: tableid == "2"})
		} else if id == "2" {
			resp.Tables = append(resp.Tables, Table{Id: "3", Name: "Пользователи", Select: tableid == "3"})
		}

		resp.Fields = []Field{}
		if tableid == "1" {
			test := [][]string{}
			test = append(test, []string{"1", "alex"})
			test = append(test, []string{"2", "john"})
			resp.Fields = append(resp.Fields, Field{Header: []string{"id", "name"}, Data: test})
			fmt.Println(test)
		}
	}
	// Тут мапа с категориями
	// append(resp.Menus, Menu{id = "1", name = "test"})
	// resp.Menus.name = ""
	// fmt.Printf("%v\n", resp)

	templates["web/menu"].Execute(writer, resp)

	for i, item := range resp.Menus {
		fmt.Printf("%v,%v\n", i, item)
	}
}

func main() {
	loadTemplate()

	http.HandleFunc("/", welcomeHandler)
	http.HandleFunc("/menu", menuHandler)

	err := http.ListenAndServe(":80", nil)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
}
