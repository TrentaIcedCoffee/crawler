package crawler

import (
	"sync"

	"io"
	"time"
)

const kChannelMaxSize = 4_000_000 // Consumes 96 MB RAM.

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
	config *Config
	logger *logger
}

func NewCrawler(config *Config) *Crawler {
	return &Crawler{
		config: config,
		logger: &logger{output_stream: nil, error_stream: nil},
	}
}

func (crawler *Crawler) OutputTo(output_stream io.Writer) *Crawler {
	crawler.logger.output_stream = output_stream
	return crawler
}

func (crawler *Crawler) ErrorTo(error_stream io.Writer) *Crawler {
	crawler.logger.error_stream = error_stream
	return crawler
}

func (crawler *Crawler) Crawl(urls []string) *Crawler {
	c := taskChannels{
		Tasks:           make(chan task, kChannelMaxSize),
		PendingTaskCnt:  make(chan int),
		FinishedTaskCnt: make(chan int),
		TotalTaskCnt:    make(chan int),
		Links:           make(chan Link),
		Errors:          make(chan error),
	}

	var wg sync.WaitGroup

	for i := 0; i < crawler.config.NumWorkers; i++ {
		wg.Add(1)
		go worker(i, crawler.config, &wg, &c, crawler.logger)
	}
	wg.Add(1)
	go controller(len(urls), &wg, &c, crawler.logger)
	wg.Add(1)
	go peeker(&c, &wg, crawler.logger)

	for _, url := range urls {
		c.Tasks <- task{
			Url:   url,
			Depth: 0,
		}
		c.TotalTaskCnt <- 1
	}

	links_c_closed := false
	errors_c_closed := false
	for {
		select {
		case link, ok := <-c.Links:
			if ok {
				crawler.logger.output("%s, %s", link.Url, link.Text)
			} else {
				links_c_closed = true
			}
		case err, ok := <-c.Errors:
			if ok {
				crawler.logger.error("%v", err)
			} else {
				errors_c_closed = true
			}
		default:
		}

		if links_c_closed && errors_c_closed {
			wg.Wait()
			break
		}
	}

	return crawler
}

type task struct {
	Url   string
	Depth int
}

type taskChannels struct {
	Tasks           chan task
	PendingTaskCnt  chan int
	FinishedTaskCnt chan int
	TotalTaskCnt    chan int
	Links           chan Link
	Errors          chan error
}

func closeAllChannels(ch *taskChannels) {
	close(ch.Tasks)
	close(ch.PendingTaskCnt)
	close(ch.FinishedTaskCnt)
	close(ch.TotalTaskCnt)
	close(ch.Links)
	close(ch.Errors)
}

func worker(id int, config *Config, wg *sync.WaitGroup, c *taskChannels, logger *logger) {
	logger.debug("Worker %d spawned", id)
	defer logger.debug("Worker %d exit", id)
	defer wg.Done()

	limiter := time.Tick(config.RequestThrottlePerWorker)

	for t := range c.Tasks {
		<-limiter
		links, errs := scrapeLinks(t.Url)
		if config.Breadth > 0 && len(links) > config.Breadth {
			links = links[:config.Breadth]
		}
		for _, err := range errs {
			c.Errors <- err
		}
		for _, link := range links {
			c.Links <- link
			if t.Depth+1 < config.Depth {
				c.Tasks <- task{
					Url:   link.Url,
					Depth: t.Depth + 1,
				}
				c.PendingTaskCnt <- 1
				c.TotalTaskCnt <- 1
			}
		}

		c.PendingTaskCnt <- -1
		c.FinishedTaskCnt <- 1
	}
}

func controller(initial_task_cnt int, wg *sync.WaitGroup, c *taskChannels, logger *logger) {
	logger.debug("Controller spawned")
	defer logger.debug("Controller exit")
	defer wg.Done()

	cnt := initial_task_cnt
	for delta := range c.PendingTaskCnt {
		cnt += delta
		if cnt == 0 {
			closeAllChannels(c)
		}
	}
}

func peeker(c *taskChannels, wg *sync.WaitGroup, logger *logger) {
	logger.debug("Peeker spawned")
	defer logger.debug("Peeker exit")
	defer wg.Done()

	peek_limiter := time.Tick(500 * time.Millisecond)
	finished, total := 0, 0

	finished_c_closed, total_c_closed := false, false
	for {
		select {
		case delta, ok := <-c.FinishedTaskCnt:
			if ok {
				finished += delta
			} else {
				finished_c_closed = true
			}
		case delta, ok := <-c.TotalTaskCnt:
			if ok {
				total += delta
			} else {
				total_c_closed = true
			}
		case <-peek_limiter:
			logger.debug("Progress %d/%d. Pending %d", finished, total, len(c.Tasks))
			if finished_c_closed && total_c_closed {
				return
			}
		}
	}
}
