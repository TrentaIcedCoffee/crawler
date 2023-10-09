package crawler

import (
	"encoding/csv"
	"fmt"
	"strconv"
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
	limiter        <-chan time.Time
	tasks          chan task
	pendingTaskCnt chan int
	links          chan Link
	errors         chan error
	finished       chan int
	total          chan int
}

func makeAllChannels(request_throttle time.Duration) *channels {
	return &channels{
		limiter:        time.Tick(request_throttle),
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

	for t := range cs.tasks {
		switch t.taskType {
		case crawlingTask:
			this.handleCrawlingTask(&t, cs)
		case pageTitleTask:
			this.handlePageTitleTask(&t, cs)
		}
	}
}

func (this *Crawler) handleCrawlingTask(t *task, cs *channels) {
	<-cs.limiter
	all_links, errs := scrapeLinks(t.url, this.config.SameHostname)
	for _, err := range errs {
		cs.errors <- err
	}

	// Pruning links.
	links := shortArray[Link]()
	for _, link := range all_links {
		should_keep, err := this.pruner.ShouldKeep(link.Url)
		if err != nil {
			cs.errors <- err
			continue
		}
		if should_keep {
			links = append(links, link)
		}
	}

	// Limiting to max breadth.
	if this.config.Breadth > 0 && len(links) > this.config.Breadth {
		links = links[:this.config.Breadth]
	}

	for _, link := range links {
		if !this.visited.add(link.Url) {
			continue
		}

		link.Depth = t.depth

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
	<-cs.limiter
	title, content, err := scrapePage(t.link.Url)
	if err != nil {
		cs.errors <- err
	}
	t.link.Title = title
	t.link.Content = content
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

	csv_writer := csv.NewWriter(this.logger.output_stream)

	// Writing the header of result csv.
	err := csv_writer.Write([]string{
		"Depth",
		"Url",
		"Text",
		"Title",
		"Content",
	})
	if err != nil {
		panic(fmt.Sprintf("FATAL error in writing csv header, %v", err))
	}
	csv_writer.Flush()
	if err := csv_writer.Error(); err != nil {
		panic(fmt.Sprintf("FATAL error in flushing csv header, %v", err))
	}

	links_c_closed := false
	errors_c_closed := false
	for {
		select {
		case link, ok := <-cs.links:
			if ok {
				err := csv_writer.Write([]string{strconv.Itoa(link.Depth), link.Url, link.Text, link.Title, link.Content})
				if err != nil {
					this.logger.error("Error in writing result, %v", err)
				}
				csv_writer.Flush()
				if err := csv_writer.Error(); err != nil {
					this.logger.error("Error in flushing result, %v", err)
				}
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
