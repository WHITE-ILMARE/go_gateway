package main

import "log"

func main() {
	rb := load_balance.LoadBalanceFactory(LbConsistentHash)
	if err := rb.Add("http://127.0.0.1:2003/base", "10"); err != nil {
		log.Println(err)
	}
	if err := rb.Add("http://127.0.0.1:2004/base", "20"); err != nil {
		log.Println(err)
	}
	proxy := NewMul
}
