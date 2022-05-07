package main

import (
	"fmt"
)

func main() {
	// 基础知识补充
	// go中，一个汉字占3个字节
	testByteArray := []byte("中")
	fmt.Printf("%v\n", testByteArray)
	// len()参数为字符串时返回的也是字节数而不是字符数
	fmt.Printf("len('中')=%d\n", len("中"))
	// 基础知识补充结束
}
