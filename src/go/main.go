package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

type Auth struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type token struct {
	Token int `json:"token"` // Не уверен что здесь нужно добавлять тег. Я ведь отправляю а не читаю
}

type Equipment struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Driver string `json:"driver"`
}

var db *sql.DB // Эту практику...

func main() {
	var err error // И эту практику подсмотрел из одного примера.
	connStr := "user=postgres password=пароль123 dbname=postgres sslmode=disable"

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Println("Не удалось подключиться к базе postgres")
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Println(err)
	} else {
		log.Println("Подключение с базой данных установлено")
	}

	http.HandleFunc("/authorize", Authorization)
	http.HandleFunc("/equipment", AddEquipment)
	http.ListenAndServe("localhost:8080", nil)
}

func Authorization(w http.ResponseWriter, r *http.Request) {
	var authTemp Auth

	if r.Method == http.MethodPost {
		err := json.NewDecoder(r.Body).Decode(&authTemp)
		if err != nil {
			w.Write([]byte("400"))
			log.Println(err)
		}
		if enter := SignUpCheck(authTemp.Login, authTemp.Password); enter { // Почему-то у меня чувство, будто я наговнокодил
			byteToken, _ := json.MarshalIndent(token{Token: 123456789}, "", "") // Не думаю что indent стоит использовать, но всё же
			w.Write([]byte("You've successfully entered to the system\n"))
			w.Write(byteToken)
		}
	}
}

func SignUpCheck(log string, pass string) bool { // Это должно работать через базу данных
	if log == "Admin" && pass == "12345" {
		return true
	}
	return false
}

func AddEquipment(w http.ResponseWriter, r *http.Request) {
	var newEquipment Equipment

	if r.Method == http.MethodPost && r.Header.Get("token") == "123456789" {
		w.Write([]byte("Succesfully entered to AddEquipment"))
		err := json.NewDecoder(r.Body).Decode(&newEquipment)
		if err != nil {
			log.Println(err)
		}

		db.Exec("INSERT INTO equipment (name, driver) values ($1, $2)", newEquipment.Name, newEquipment.Driver)
	} else {
		w.Write([]byte("400"))
	}
}

/*
--- Начало
	- Нужно чтоб через постмэн отправился логин и пароль на /login
	- Сервер проверил данные
	- Ответил Входом или не входом сообщением ок не ок.
	- После этого загорается флаг isAuthinticated и можно использовать что-то далее.
	= ВСЁ ЭТА КОММУНИКАЦИЯ ПРОХОДИТ ЧЕРЕЗ JSON.
--- Законечено
	- Добавить новый /equipment который принимает JSON с нужными полями для техники
	- Добавить базу данных Postgres для сохранения добавленных данных в ней (нужно отправку реализовать)
	- Доступ к эндпоинту /equipment должен проходить через токен
*/
