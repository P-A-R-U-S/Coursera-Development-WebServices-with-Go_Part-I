package main

import (
	"bufio"
	"fmt"
	"os"
)

func main()  {

	in := bufio.NewScanner(os.Stdin)
	for in.Scan() {
		txt := in.Text()
		fmt.Println("-", txt)
	}
}
