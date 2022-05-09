package load_balance

type LbType int

const (
	LbRandom LbType = iota // 高级，自增序数
	LbRoundRobin
	LbWeightRoundRobin
	LbConsistentHash
)

func LoadBalanceFactory(lbType LbType) LoadBalance {
	switch lbType {
	case LbRandom:
		return &RandomBalance{}
	case LbRoundRobin:
		return &RoundRobinBalance{}
	case LbWeightRoundRobin:
		return &RoundRobinBalance{}
	case LbConsistentHash:
		return NewConsistentHashBalance(10, nil)
	default:
		return &RandomBalance{}
	}
}
