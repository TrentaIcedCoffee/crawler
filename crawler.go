package crawler

import (
	"sync"

	"io"
	"time"
)

const kChannelMaxSize = 4_000_000 // Consumes 96 MB RAM.
const kChannelDefautSize = 10

type Link struct {
	Url  string
	Text string
}

type Config struct {
	Depth                    int
	Breadth                  int
	NumWorkers               int
	RequestThrottlePerWorker time.Duration
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
	go this.controller(len(urls), &wg, cs)
	wg.Add(1)
	go this.peeker(&wg, cs)

	for _, url := range urls {
		cs.tasks <- task{
			url:   url,
			depth: 0,
		}
		cs.total <- 1
	}

	links_c_closed := false
	errors_c_closed := false
	for {
		select {
		case link, ok := <-cs.links:
			if ok {
				this.logger.output("%s,%s", link.Url, toCsv(link.Text))
			} else {
				links_c_closed = true
			}
		case err, ok := <-cs.errors:
			if ok {
				this.logger.error("%v", err)
			} else {
				errors_c_closed = true
			}
		}

		if links_c_closed && errors_c_closed {
			wg.Wait()
			break
		}
	}

	return this
}

type task struct {
	url   string
	depth int
}

type channels struct {
	tasks          chan task
	pendingTaskCnt chan int
	links          chan Link
	errors         chan error
	finished       chan int
	total          chan int
}

func makeAllChannels() *channels {
	return &channels{
		tasks:          make(chan task, kChannelMaxSize),
		pendingTaskCnt: make(chan int, kChannelDefautSize),
		links:          make(chan Link, kChannelDefautSize),
		errors:         make(chan error, kChannelDefautSize),
		finished:       make(chan int, kChannelDefautSize),
		total:          make(chan int, kChannelDefautSize),
	}
}

func closeAllChannels(cs *channels) {
	close(cs.tasks)
	close(cs.pendingTaskCnt)
	close(cs.links)
	close(cs.errors)
	close(cs.finished)
	close(cs.total)
}

func (this *Crawler) worker(id int, wg *sync.WaitGroup, cs *channels) {
	this.logger.debug("Worker %d spawned", id)
	defer this.logger.debug("Worker %d exit", id)
	defer wg.Done()

	limiter := time.Tick(this.config.RequestThrottlePerWorker)

	for t := range cs.tasks {
		<-limiter
		links, errs := scrapeLinks(t.url)
		if this.config.Breadth > 0 && len(links) > this.config.Breadth {
			links = links[:this.config.Breadth]
		}
		for _, err := range errs {
			cs.errors <- err
		}
		for _, link := range links {
			if this.visited.has(hash(link.Url)) {
				continue
			}
			// Ideally we should check if the add inserts a new value, but this is fine.
			this.visited.add(hash(link.Url))
			cs.links <- link
			if t.depth+1 < this.config.Depth {
				cs.tasks <- task{
					url:   link.Url,
					depth: t.depth + 1,
				}
				cs.pendingTaskCnt <- 1
				cs.total <- 1
			}
		}

		cs.pendingTaskCnt <- -1
		cs.finished <- 1
	}
}

func (this *Crawler) controller(initial_task_cnt int, wg *sync.WaitGroup, cs *channels) {
	this.logger.debug("Controller spawned")
	defer this.logger.debug("Controller exit")
	defer wg.Done()

	cnt := initial_task_cnt
	for delta := range cs.pendingTaskCnt {
		cnt += delta
		if cnt == 0 {
			closeAllChannels(cs)
		}
	}
}

func (this *Crawler) peeker(wg *sync.WaitGroup, cs *channels) {
	this.logger.debug("Peeker spawned")
	defer this.logger.debug("Peeker exit")
	defer wg.Done()

	peek_limiter := time.Tick(500 * time.Millisecond)
	finished, total := 0, 0
	finished_closed, total_closed := false, false

	for {
		select {
		case delta, ok := <-cs.finished:
			if ok {
				finished += delta
			} else {
				finished_closed = true
			}
		case delta, ok := <-cs.total:
			if ok {
				total += delta
			} else {
				total_closed = true
			}
		case <-peek_limiter:
			this.logger.debug("Progress %d/%d. Queued %d", finished, total, len(cs.tasks))
			if finished_closed && total_closed {
				return
			}
		}
	}
}
