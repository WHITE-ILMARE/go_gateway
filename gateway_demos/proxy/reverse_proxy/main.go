package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

func main() {
	// 基础知识补充
	// go中，一个汉字占3个字节
	testByteArray := []byte("中")
	fmt.Printf("%v\n", testByteArray)
	// len()参数为字符串时返回的也是字节数而不是字符数
	fmt.Printf("len('中')=%d\n", len("中"))
	// 基础知识补充结束

	modify := func(resp *http.Response) error {
		// 请求以下命令：curl 'http://127.0.0.1:2002/error
		if resp.StatusCode != 200 {
			// 获取内容
			oldPayload, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}
			// []byte顾名思义，字节数组，每个元素都是一个字节，8位，存的是0-255之间的整数，即一个ASCII码
			// 只有字节流才能在网络中传输，所以需要转化成字节数组再操作
			newPayLoad := []byte("StatusCode error:" + string(oldPayload))
			resp.Body = ioutil.NopCloser(bytes.NewBuffer(newPayLoad))
			resp.ContentLength = int64(len(newPayLoad))
			resp.Header.Set("Content-Length", strconv.FormatInt(int64(len(newPayLoad)), 10))
		}
		return nil
	}
}
