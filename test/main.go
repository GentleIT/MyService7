package main

import (
	"io"
	"log"
	"net/http"
)

func main() {

	// t := time.Now()
	// log.Print(t.Format("2006-01-02"))

	resp, _ := http.Get("https://wttr.in/Алматы?format=3")
	body, _ := io.ReadAll(resp.Body)
	log.Println(string(body))

	// resp, err := http.Get("https://gobyexample.com")
	// if err != nil {
	// 	panic(err)
	// }
	// defer resp.Body.Close()

	// fmt.Println("Response status:", resp.Status)
	// fmt.Println(resp)

	// scanner := bufio.NewScanner(resp.Body)
	// for i := 0; scanner.Scan() && i < 5; i++ {
	// 	fmt.Println(scanner.Text())
	// }

	// if err := scanner.Err(); err != nil {
	// 	panic(err)
	// }
}
