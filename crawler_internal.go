package crawler

import (
	"sync"
	"time"
)

const kChannelMaxSize = 4_000_000 // Consumes 96 MB RAM.
const kChannelDefautSize = 100

type taskType int

const (
	crawlingTask taskType = iota
	pageTitleTask
)

type task struct {
	taskType taskType
	url      string
	depth    int
	link     Link
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

func (this *Crawler) addTask(cs *channels, task task) {
	cs.tasks <- task
	cs.pendingTaskCnt <- 1
	cs.total <- 1
}

func (this *Crawler) finishedTask(cs *channels) {
	cs.pendingTaskCnt <- -1
	cs.finished <- 1
}

func (this *Crawler) worker(id int, wg *sync.WaitGroup, cs *channels) {
	this.logger.debug("Worker %d spawned", id)
	defer this.logger.debug("Worker %d exit", id)
	defer wg.Done()

	limiter := time.Tick(this.config.RequestThrottlePerWorker)

	for t := range cs.tasks {
		<-limiter
		switch t.taskType {
		case crawlingTask:
			this.handleCrawlingTask(&t, cs)
		case pageTitleTask:
			this.handlePageTitleTask(&t, cs)
		}
	}
}

func (this *Crawler) handleCrawlingTask(t *task, cs *channels) {
	links, errs := scrapeLinks(t.url, this.config.SameHostname)
	if this.config.Breadth > 0 && len(links) > this.config.Breadth {
		links = links[:this.config.Breadth]
	}
	for _, err := range errs {
		cs.errors <- err
	}
	for _, link := range links {
		if !this.visited.add(link.Url) {
			continue
		}
		this.addTask(cs, task{taskType: pageTitleTask, link: link})
		if t.depth+1 < this.config.Depth {
			this.addTask(cs, task{
				taskType: crawlingTask,
				url:      link.Url,
				depth:    t.depth + 1,
			})
		}
	}
	this.finishedTask(cs)
}

func (this *Crawler) handlePageTitleTask(t *task, cs *channels) {
	title, err := scrapeTitle(t.link.Url)
	if err != nil {
		cs.errors <- err
	}
	t.link.Title = title
	cs.links <- t.link
	this.finishedTask(cs)
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

func (this *Crawler) emitter(wg *sync.WaitGroup, cs *channels) {
	this.logger.debug("Emitter spawned")
	defer this.logger.debug("Emitter exit")
	defer wg.Done()

	links_c_closed := false
	errors_c_closed := false
	for {
		select {
		case link, ok := <-cs.links:
			if ok {
				this.logger.output("%s,%s,%s", link.Url, toCsvRow(link.Text), toCsvRow(link.Title))
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
			return
		}
	}

}
