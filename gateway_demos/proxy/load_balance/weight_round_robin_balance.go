package load_balance

import (
	"errors"
	"strconv"
)

type WeightNode struct {
	addr            string
	weight          int // 节点基础权重，需要初始化
	currentWeight   int // 节点当前权重，不初始化，默认为0
	effectiveWeight int // 有效权重，需要初始化为weight，为简化计算，本例中不改变effectiveWeight
}

/**
1. currentWeight = currentWeight + effectiveWeight
2. 选出最大的currentWeight节点作为选中节点
3. currentWeight = currentWeight - totalWeight
此算法中，effectiveWeight一直是初始的weight没有改变，第一步给每个节点加effectiveWeight，
一共加了totalWeight，第三步给被选中的节点减去totalWeight，所以一轮结束，totalWeight不改变。
即使effectiveWeight变化，currentWeight之和也不会变化
*/

type WeightRoundRobinBalance struct {
	curIndex int
	rss      []*WeightNode
	//conf LoadBalanceConf
}

func (r *WeightRoundRobinBalance) Add(params ...string) error {
	if len(params) != 2 {
		return errors.New("params need len 2")
	}
	weight, err1 := strconv.ParseInt(params[1], 10, 64)
	if err1 != nil {
		return err1
	}
	weightNode := &WeightNode{weight: int(weight), addr: params[0]}
	weightNode.effectiveWeight = weightNode.weight
	r.rss = append(r.rss, weightNode)
	return nil
}

func (r *WeightRoundRobinBalance) Next() string {
	total := 0
	var best *WeightNode
	for i := 0; i < len(r.rss); i++ {
		w := r.rss[i]
		// step 1 计算totalWeight
		total += w.effectiveWeight
		// step 2 变更currentWeight
		w.currentWeight += w.effectiveWeight
		// step 3 effectiveWeight默认为weight，通讯异常-1，成功+1，上限为weight大小
		if w.effectiveWeight < w.weight {
			w.effectiveWeight++
		}
		// step 4 选择最大currentWeight节点
		if best == nil || w.currentWeight > best.currentWeight {
			best = w
		}
	}
	if best == nil {
		return ""
	}
	// step 5 变更best的currentWeight
	best.currentWeight -= total
	return best.addr
}
