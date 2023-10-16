package crawler

import (
	"sync"

	"io"
	"time"
)

type Link struct {
	Depth   int
	Url     string
	Text    string
	Title   string
	Content string
}

type Config struct {
	Depth            int
	Breadth          int
	NumWorkers       int
	RequestThrottler time.Duration
}

type Crawler struct {
	config  *Config
	visited *concurrentSet
	logger  *logger
	pruner  Pruner
}

func NewCrawler(config *Config, output_stream io.Writer, error_stream io.Writer, pruner Pruner) *Crawler {
	return &Crawler{
		config:  config,
		visited: newConcurrentSet(),
		logger:  &logger{output_stream: output_stream, error_stream: error_stream},
		pruner:  pruner,
	}
}

func (this *Crawler) Crawl(urls []string) *Crawler {
	throttlers := makeThrottlers(this.config.RequestThrottler, urls)
	cs := makeAllChannels()

	var wg sync.WaitGroup

	for i := 0; i < this.config.NumWorkers; i++ {
		wg.Add(1)
		go this.worker(i, &wg, cs, throttlers)
	}
	wg.Add(1)
	go this.peeker(&wg, cs)
	wg.Add(1)
	go this.emitter(&wg, cs)

	for _, url := range urls {
		this.addTask(cs, task{
			taskType: crawlingTask,
			url:      url,
			depth:    0,
		})
	}

	pendingTasks := 0
	noTask := true // True if no task has been created.
	for delta := range cs.pendingTaskCnt {
		pendingTasks += delta
		if noTask && pendingTasks != 0 {
			noTask = false
		}
		if !noTask && pendingTasks == 0 {
			closeAllChannels(cs)
			wg.Wait()
			break
		}
	}

	return this
}
