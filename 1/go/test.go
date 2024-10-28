package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("请输入一些内容：")
	input, _ := reader.ReadString('\n')
	fmt.Println("你输入的内容是：", input)
}
