package crawler

import (
	"sync"

	"time"
)

type Link struct {
	Url  string
	Text string
}

type Config struct {
	Depth      int
	Breadth    int
	NumWorkers int
	IsDebug    bool
}

func DefaultConfig() *Config {
	return &Config{
		Depth:      3,
		Breadth:    0,
		NumWorkers: 10,
		IsDebug:    true,
	}
}

func Crawl(urls []string, config *Config) {
	c := taskChannels{
		Limiter:        time.Tick(200 * time.Millisecond),
		Tasks:          make(chan task, kChannelMaxSize),
		PendingTaskCnt: make(chan int),
		Links:          make(chan Link),
		Errors:         make(chan error),
	}

	var wg sync.WaitGroup

	logger := logger{isDebug: config.IsDebug}

	for i := 0; i < config.NumWorkers; i++ {
		wg.Add(1)
		go worker(i, config, &wg, &c, &logger)
	}
	go controller(len(urls), &wg, &c, &logger)

	for _, url := range urls {
		c.Tasks <- task{
			Url:   url,
			Depth: 0,
		}
	}

	links_c_closed := false
	errors_c_closed := false
	for {
		select {
		case link, ok := <-c.Links:
			if ok {
				logger.Output("%s, %s", link.Url, link.Text)
			} else {
				links_c_closed = true
			}
		case err, ok := <-c.Errors:
			if ok {
				logger.Error("%v", err)
			} else {
				errors_c_closed = true
			}
		default:
		}

		if links_c_closed && errors_c_closed {
			break
		}
	}
}

type task struct {
	Url   string
	Depth int
}

type taskChannels struct {
	Limiter        <-chan time.Time
	Tasks          chan task
	PendingTaskCnt chan int
	Links          chan Link
	Errors         chan error
}

const kChannelMaxSize = 1_000_000

func worker(id int, config *Config, wg *sync.WaitGroup, c *taskChannels, logger *logger) {
	logger.Debug("Worker %d spawned", id)
	defer logger.Debug("Worker %d exit", id)
	defer wg.Done()

	for t := range c.Tasks {
		if t.Depth < config.Depth {
			<-c.Limiter
			links, errs := scrapeLinks(t.Url)
			if config.Breadth > 0 {
				links = links[:config.Breadth]
			}
			for _, err := range errs {
				c.Errors <- err
			}
			for _, link := range links {
				c.Links <- link
				c.Tasks <- task{
					Url:   link.Url,
					Depth: t.Depth + 1,
				}
				c.PendingTaskCnt <- 1
			}
		}

		c.PendingTaskCnt <- -1
	}
}

func controller(initial_task_cnt int, wg *sync.WaitGroup, c *taskChannels, logger *logger) {
	logger.Debug("Controller spawned")
	defer logger.Debug("Controller exit")

	cnt := initial_task_cnt
	for delta := range c.PendingTaskCnt {
		cnt += delta
		if cnt == 0 {
			close(c.Tasks)
			wg.Wait()
			close(c.PendingTaskCnt)
			close(c.Links)
			close(c.Errors)
		}
	}
}
