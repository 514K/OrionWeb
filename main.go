package main

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"text/template"

	_ "github.com/lib/pq"
)

var connStr = "user=postgres password=Sudozo39! dbname=orion_1_3_0_0 sslmode=disable host=192.168.1.105"
var db, err = sql.Open("postgres", connStr)

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

		result, err := db.Query("SELECT * FROM users WHERE login = $1", string(request.Form["login"][0]))
		if err != nil {
			panic(err)
		}

		var v1, v2 string
		for result.Next() {
			result.Scan(&v1, &v2)
		}

		// fmt.Printf("%v %v\n", v1, v2)
		hasher := md5.New()
		hasher.Write([]byte(request.Form["password"][0]))
		// fmt.Printf("%v\n", hex.EncodeToString(hasher.Sum(nil)) == myhash)
		if v2 == hex.EncodeToString(hasher.Sum(nil)) {
			// fmt.Printf("%v\n", request.Form["checkbox"][0] == "on")
			if len(request.Form["checkbox"]) != 0 {
				// ТУТ ЗАПОМИНАЕМ
				fmt.Printf("check on\n")
			}
			cooka := http.Cookie{Name: "username", Value: request.Form["login"][0]}
			http.SetCookie(writer, &cooka)
			http.Redirect(writer, request, "/menu", http.StatusSeeOther)
		} else {
			templates["web/auth"].Execute(writer, "Не верный логин или пароль")
		}

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
		// println(request.URL.Query().Get("id"))

		id := request.URL.Query().Get("id")
		// Тут запрос на меню
		result, err := db.Query("SELECT * FROM menu")
		if err != nil {
			panic(err)
		}

		var v1, v2, v3, v4 string
		for result.Next() {
			result.Scan(&v1, &v2, &v3, &v4)
			fmt.Printf("%v %v %v %v\n", v1, v2, v3, v4)
			resp.Menus = append(resp.Menus, Menu{Id: v1, Name: v2, Select: id == v1})
		}

		tableid := request.URL.Query().Get("tableid")
		// Тут запрос на таблицы
		resp.Tables = []Table{}

		mId, _ := strconv.Atoi(id)
		result, err = db.Query("SELECT * FROM tables WHERE id IN (SELECT tableid FROM tablesofmenu WHERE menu = $1)", mId)
		if err != nil {
			panic(err)
		}

		// var v5, v6, v7, v8 string
		for result.Next() {
			result.Scan(&v1, &v2, &v3, &v4)
			resp.Tables = append(resp.Tables, Table{Id: v1, Name: v2, Select: tableid == v1})
			// fmt.Printf("%v\n", result)

			// fmt.Printf("%v %v %v %v\n", v5, v6, v7, v8)
		}

		resp.Fields = []Field{}

		tId, _ := strconv.Atoi(tableid)
		result, err = db.Query("SELECT insname FROM tables WHERE id = $1", tId)
		if err != nil {
			panic(err)
		}

		tabName := ""
		if result.Next() {
			result.Scan(&tabName)

			result, err = db.Query("SELECT * FROM " + tabName)
			if err != nil {
				panic(err)
			}
			head, _ := result.Columns()
			dat := [][]string{}
			vals := make([]interface{}, len(head))
			for result.Next() {
				for i := range head {
					vals[i] = new(string)
				}

				result.Scan(vals...)

				row := []string{}

				for i := range vals {
					row = append(row, *(vals[i].(*string)))
					fmt.Printf("Column %v: %s\n", i, *(vals[i].(*string)))
				}
				dat = append(dat, row)

				// row := []string{}

				// for i := range head {
				// 	val := values[i]

				// 	// b, ok := val.([]byte)
				// 	// fmt.Printf("%v\n", b)
				// 	// var v interface{}
				// 	// if ok {
				// 	// 	v = string(b)
				// 	// } else {
				// 	// 	v = val
				// 	// }

				// 	// fmt.Printf("%v\n", val)

				// 	a, _ := val.([]byte)
				// 	row = append(row, string(a[:]))
				// 	fmt.Printf("%v\n", val)

				// 	// test = append(test, []string{"1", "alex"})

				// 	// fmt.Println(col, v)
				// }
				// dat = append(dat, row)
				// row := []string{}
				// for range head {
				// 	values := make([]interface{}, len(head))
				// 	result.Scan(values)
				// 	row = append(row, v1)
				// }
				// values := make([]interface{}, len(head))
				// result.Scan(&values)

				// fmt.Printf("%v\n", values)
			}
			resp.Fields = append(resp.Fields, Field{Header: head, Data: dat})
			fmt.Printf("%v\n", resp.Fields)
		}

		// fmt.Printf("%v\n", tabName)

		// values := make([]interface{}, len(head))
		// valuePtrs := make([]interface{}, len(head))

		// tabName := ""
		// for result.Next() {
		// 	result.Scan(&tabName)
		// test := []string{}
		// for range head {
		// 	result.Scan(v1)
		// 	test = append(test, v1)
		// }
		// fmt.Printf("%v\n", test)
		// result.Scan(valuePtrs...)
		// for i, col := range head {
		// 	val := values[i]

		// 	b, ok := val.([]byte)
		// 	var v interface{}
		// 	if ok {
		// 		v = string(b)
		// 	} else {
		// 		v = val
		// 	}

		// 	fmt.Println("sas", col, v)
		// }
		// fmt.Printf("%v\n", len(head))
		// for i := 0; i < len(head); i++ {

		// }
		// result.Scan(&v1, &v2)

		// resp.Tables = append(resp.Tables, Table{Id: v1, Name: v2, Select: tableid == v1})
		// }
		// resp.Fields = append(resp.Fields, Field{Header: head, Data: test})

		// if tableid == "1" {
		// 	test := [][]string{}
		// 	test = append(test, []string{"1", "alex"})
		// 	test = append(test, []string{"2", "john"})
		// 	resp.Fields = append(resp.Fields, Field{Header: []string{"id", "name"}, Data: test})
		// 	// fmt.Println(test)
		// }
	}
	// Тут мапа с категориями
	// append(resp.Menus, Menu{id = "1", name = "test"})
	// resp.Menus.name = ""
	// fmt.Printf("%v\n", resp)

	templates["web/menu"].Execute(writer, resp)

	// for i, item := range resp.Menus {
	// 	fmt.Printf("%v,%v\n", i, item)
	// }
}

func main() {
	// test := md5.Sum([]byte("1"))
	// myhash := "c4ca4238a0b923820dcc509a6f75849b"

	// hasher := md5.New()
	// hasher.Write([]byte("1"))
	// fmt.Printf("%v\n", hex.EncodeToString(hasher.Sum(nil)) == myhash)

	if err != nil {
		fmt.Printf("sas\n")
		panic(err)
	}
	defer db.Close()

	http.Handle("/web/", http.StripPrefix("/web", http.FileServer(http.Dir("./web"))))

	loadTemplate()

	http.HandleFunc("/", welcomeHandler)
	http.HandleFunc("/menu", menuHandler)

	err = http.ListenAndServe(":80", nil)
	if err != nil {
		fmt.Printf("%v\n", err)
	}

	// fmt.Println(result.LastInsertId())  // не поддерживается
	// fmt.Println(result.RowsAffected())  // количество добавленных строк
}
