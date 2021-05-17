//
// main.go
//
// May 2021, Prashant Desai
//

package main

import (
	"WardrobeManagerMS/pkg/api"
	"WardrobeManagerMS/pkg/repository"
	"fmt"
)

func main() {

	fmt.Printf("Version for definitions : %s\n", api.Version)
	fmt.Printf("Version for repository : %s\n", repository.Version)
}
