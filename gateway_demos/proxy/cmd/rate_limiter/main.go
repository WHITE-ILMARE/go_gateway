package main

// 这个例子是time/rate库的使用示例
// Wait(), Reserve(), Allow()操作都会去请求资源，只不过返回值不同
// Wait()返回error， Reserve()返回等待时间，Allow()返回bool
// 由于限流，for循环不会执行得很快，而是会阻塞执行，也体现了限流的意义
// 三个方法底层都是调用ResereN()实现的

import (
	"context"
	"fmt"
	"golang.org/x/time/rate"
	"log"
	"time"
)

func main() {
	l := rate.NewLimiter(1, 5)
	log.Println(l.Limit(), l.Burst())
	for i := 0; i < 100; i++ {
		// 阻塞等待，直到取得一个token
		log.Println("before wait")
		c, _ := context.WithTimeout(context.Background(), time.Second*2)
		if err := l.Wait(c); err != nil {
			log.Println("Limiter wait err:" + err.Error())
		}
		log.Println("after wait")
		// 返回需要等待多久才有新的token，这样就可以等待指定时间执行任务
		r := l.Reserve()
		log.Println("reserve Delay:", r.Delay())

		// 判断当前是否可以取到token
		a := l.Allow()
		fmt.Println("Allow:", a)
	}
}
