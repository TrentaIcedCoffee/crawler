# Concurrent Web Crawler

Concurrent web crawler in Go using goroutines.

## Usage

In your main func

```go
crawler.NewCrawler(&crawler.Config{
		Depth:            3,
		Breadth:          0,                      // Using breadth = 0 to get all links on a page.
		NumWorkers:       100,
		RequestThrottler: 100 * time.Millisecond, // 10 requests per second for each **domain**.
	}, outputFile, errorFile, &crawler.SameDomain{}).Crawl([]string{
		"https://example.com",
		"https://another.com"
	})
```

Output is in CSV format of `<depth>,<url>,<text>,<page_title>,<page_content>`. Can be directly loaded to Pandas.

## Design

![design](/graphes/design.jpg)

## How it works

![diagram](/graphes/diagram.jpg)
