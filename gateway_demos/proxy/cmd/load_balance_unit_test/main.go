package load_balance_unit_test

import (
	"fmt"
	"github.com/WHITE-ILMARE/go_gateway/gateway_demos/proxy/load_balance"
	"testing"
)

// 运行时需同时test本文件及相应的Load Balance定义文件
func TestRandomBalance(t *testing.T) {
	rb := load_balance.NewConsistentHashBalance(10, nil)
	rb.Add("127.0.0.1:2003") // 0
	rb.Add("127.0.0.1:2004") // 1
	rb.Add("127.0.0.1:2005") // 2
	rb.Add("127.0.0.1:2006") // 3
	rb.Add("127.0.0.1:2007") // 4

	fmt.Println(rb.Get("http://127.0.0.1:2002/base/getinfo"))
	fmt.Println(rb.Get("http://127.0.0.1:2002/base/error"))
	fmt.Println(rb.Get("http://127.0.0.1:2002/base/getinfo"))
	fmt.Println(rb.Get("http://127.0.0.1:2002/base/changepwd"))

	fmt.Println(rb.Get("127.0.0.1"))
	fmt.Println(rb.Get("192.168.0.1"))
	fmt.Println(rb.Get("127.0.0.1"))
}
