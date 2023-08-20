package crawler

import (
	"sync"

	"io"
	"time"
)

type Link struct {
	Url   string
	Text  string
	Title string
}

type Config struct {
	Depth                    int
	Breadth                  int
	NumWorkers               int
	RequestThrottlePerWorker time.Duration
	SameHostname             bool
}

type Crawler struct {
	config  *Config
	visited *concurrentSet
	logger  *logger
}

func NewCrawler(config *Config) *Crawler {
	return &Crawler{
		config:  config,
		visited: newConcurrentSet(),
		logger:  &logger{output_stream: nil, error_stream: nil},
	}
}

func (this *Crawler) OutputTo(output_stream io.Writer) *Crawler {
	this.logger.output_stream = output_stream
	return this
}

func (this *Crawler) ErrorTo(error_stream io.Writer) *Crawler {
	this.logger.error_stream = error_stream
	return this
}

func (this *Crawler) Crawl(urls []string) *Crawler {
	cs := makeAllChannels()

	var wg sync.WaitGroup

	for i := 0; i < this.config.NumWorkers; i++ {
		wg.Add(1)
		go this.worker(i, &wg, cs)
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
