package load_balance

type LoadBalance interface {
	Add(...string) error
	Get(key string) (string, error)
	// Update()
}
