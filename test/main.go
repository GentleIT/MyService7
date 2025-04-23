package main

import (
	"log"
	"time"
)

func main() {

	t := time.Now()
	log.Print(t.Format("2006-01-02"))

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
