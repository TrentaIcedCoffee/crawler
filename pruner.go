package crawler

type Pruner interface {
	ShouldKeep(parent_url string, child_url string) (bool, error)
}

type SameDomain struct{}

func (this *SameDomain) ShouldKeep(parent_url string, child_url string) (bool, error) {
	return isSameHost(parent_url, child_url)
}
