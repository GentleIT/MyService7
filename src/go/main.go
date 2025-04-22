package main

import (
	"encoding/json"
	"log"
	"net/http"
	// _ "github.com/lib/pq"
)

type Auth struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type token struct {
	Token int // `json:"token"` // Не уверен что здесь нужно добавлять тег. Я ведь отправляю а не читаю
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

func main() {
	http.HandleFunc("/authorize", Authorization)
	http.ListenAndServe("localhost:8080", nil)
}

func SignUpCheck(log string, pass string) bool { // Это должно работать через базу данных
	if log == "Admin" && pass == "12345" {
		return true
	}
	return false
}

/*
	- Нужно чтоб через постмэн отправился логин и пароль на /login
	- Сервер проверил данные
	- Ответил Входом или не входом сообщением ок не ок.
	- После этого загорается флаг isAuthinticated и можно использовать что-то далее.
	= ВСЁ ЭТА КОММУНИКАЦИЯ ПРОХОДИТ ЧЕРЕЗ JSON.
*/
