package crawler

type Pruner interface {
	ShouldKeep(url string) (bool, error)
}

type NoOpPruner struct{}

func (this *NoOpPruner) ShouldKeep(url string) (bool, error) {
	return true, nil
}
