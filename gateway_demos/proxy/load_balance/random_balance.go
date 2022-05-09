package load_balance

import (
	"errors"
	"math/rand"
)

// RandomBalance 实现随机负载均衡
type RandomBalance struct {
	curIndex int
	// 目标服务器数组
	rss []string
	// 观察者主体，服务发现时会用到
	//conf LoadBalanceConf
}

func (r *RandomBalance) Add(params ...string) error {
	if len(params) == 0 {
		return errors.New("param len 1 at least")
	}
	addr := params[0]
	r.rss = append(r.rss, addr)
	return nil
}

func (r *RandomBalance) Next() string {
	if len(r.rss) == 0 {
		return ""
	}
	r.curIndex = rand.Intn(len(r.rss))
	return r.rss[r.curIndex]
}

// 简单包装了Next()
func (r *RandomBalance) Get(key string) (string, error) {
	return r.Next(), nil
}
