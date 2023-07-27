package crawler

import (
	"sync"
	"sync/atomic"

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
	config   *Config
	visited  *concurrentSet
	logger   *logger
	finished uint64
	total    uint64
}

func NewCrawler(config *Config) *Crawler {
	return &Crawler{
		config:   config,
		visited:  newConcurrentSet(),
		logger:   &logger{output_stream: nil, error_stream: nil},
		finished: 0,
		total:    0,
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
		Tasks:          make(chan task, kChannelMaxSize),
		PendingTaskCnt: make(chan int),
		Links:          make(chan Link),
		Errors:         make(chan error),
	}

	var wg sync.WaitGroup

	for i := 0; i < crawler.config.NumWorkers; i++ {
		wg.Add(1)
		go worker(i, crawler, &wg, &c)
	}
	wg.Add(1)
	go controller(len(urls), &wg, &c, crawler.logger)
	wg.Add(1)
	go peeker(crawler, &c, &wg, crawler.logger)

	for _, url := range urls {
		c.Tasks <- task{
			Url:   url,
			Depth: 0,
		}
		atomic.AddUint64(&crawler.total, 1)
	}

	links_c_closed := false
	errors_c_closed := false
	for {
		select {
		case link, ok := <-c.Links:
			if ok {
				crawler.logger.output("%s,%s", link.Url, toCsv(link.Text))
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
	Tasks          chan task
	PendingTaskCnt chan int
	Links          chan Link
	Errors         chan error
}

func closeAllChannels(ch *taskChannels) {
	close(ch.Tasks)
	close(ch.PendingTaskCnt)
	close(ch.Links)
	close(ch.Errors)
}

func worker(id int, crawler *Crawler, wg *sync.WaitGroup, c *taskChannels) {
	crawler.logger.debug("Worker %d spawned", id)
	defer crawler.logger.debug("Worker %d exit", id)
	defer wg.Done()

	limiter := time.Tick(crawler.config.RequestThrottlePerWorker)

	for t := range c.Tasks {
		<-limiter
		links, errs := scrapeLinks(t.Url)
		if crawler.config.Breadth > 0 && len(links) > crawler.config.Breadth {
			links = links[:crawler.config.Breadth]
		}
		for _, err := range errs {
			c.Errors <- err
		}
		for _, link := range links {
			if crawler.visited.has(hash(link.Url)) {
				continue
			}
			// Ideally we should check if the add inserts a new value, but this is fine.
			crawler.visited.add(hash(link.Url))
			c.Links <- link
			if t.Depth+1 < crawler.config.Depth {
				c.Tasks <- task{
					Url:   link.Url,
					Depth: t.Depth + 1,
				}
				c.PendingTaskCnt <- 1
				atomic.AddUint64(&crawler.total, 1)
			}
		}

		c.PendingTaskCnt <- -1
		atomic.AddUint64(&crawler.finished, 1)
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

func peeker(crawler *Crawler, c *taskChannels, wg *sync.WaitGroup, logger *logger) {
	logger.debug("Peeker spawned")
	defer logger.debug("Peeker exit")
	defer wg.Done()

	peek_limiter := time.Tick(500 * time.Millisecond)

	for {
		<-peek_limiter
		logger.debug("Progress %d/%d. Queued %d", crawler.finished, crawler.total, len(c.Tasks))
		if len(c.Tasks) == 0 && crawler.finished == crawler.total {
			break
		}
	}
}
