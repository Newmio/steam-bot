package main

import (
	"bot/internal/app"

	_ "github.com/mattn/go-sqlite3"
)

// chromedriver --port=9515
// docker system prune -a --volumes

func main() {
	// strs := []string{
	// 	"+ 3 4",
	// 	"+ 3 4 - 1 2",
	// 	"+ 3 4 - 1 2 + 2 3",
	// }

	// for _, value := range strs {
	// 	fmt.Println(parse(value))
	// }

	app.Init()
}

// func parse(str string) string {
// 	var nums []string
// 	var s []string
// 	str = strings.Replace(str, " ", "", -1)

// 	for _, value := range str {

// 		_, err := strconv.Atoi(string(value))
// 		if err != nil {
// 			s = append(s, string(value))

// 		} else {
// 			nums = append(nums, string(value))
// 		}
// 	}

// 	var result string

// 	for i := 0; i < len(nums); i += 2{

// 	}

// 	for _, value := range s{
// 		result = strings.Replace(result, " ", value, 1)
// 	}

// 	return result
// }
