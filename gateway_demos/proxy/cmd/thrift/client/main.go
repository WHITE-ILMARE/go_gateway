package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/WHITE-ILMARE/go_gateway/gateway_demos/proxy/thrift_server_client/gen-go/thrift_gen"
	"github.com/apache/thrift/lib/go/thrift"
	"log"
	"os"
	"time"
)

func main() {
	// 真实thrift服务地址6001
	//addr := flag.String("addr", "127.0.0.1:6001", "input addr")
	// 代理thrift服务地址2002
	addr := flag.String("addr", "127.0.0.1:2002", "input addr")
	flag.Parse()
	if *addr == "" {
		flag.Usage()
		os.Exit(1)
	}
	for {
		tSocket, err := thrift.NewTSocket(*addr)
		if err != nil {
			log.Fatalln("tSocket error:", err)
		}
		transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
		transport, _ := transportFactory.GetTransport(tSocket)
		protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
		client := thrift_gen.NewFormatDataClientFactory(transport, protocolFactory)
		if err := transport.Open(); err != nil {
			log.Fatalln("Error opening:", *addr)
		}
		defer transport.Close()
		data := thrift_gen.Data{Text: "ping"}
		d, err := client.DoFormat(context.Background(), &data)
		if err != nil {
			fmt.Println("err:", err.Error())
		} else {
			fmt.Println("Text:", d.Text)
		}
		time.Sleep(40 * time.Millisecond)
	}
}
