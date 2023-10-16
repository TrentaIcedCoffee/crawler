# Concurrent Web Crawler

Concurrent web crawler in Go using goroutines.

## Usage

In your main func

```go
crawler.NewCrawler(&crawler.Config{
		Depth:            3,
		Breadth:          0,                      // Use breadth = 0 to get all links on a page.
		NumWorkers:       50,
		RequestThrottler: 100 * time.Millisecond, // 10 requests per second for each **domain**.
	}, os.Stdout, os.Stderr, &crawler.SameDomain{}).Crawl([]string{
		"https://example.com",
		"https://another.com"
	})
```

Output is in CSV format of `<depth>,<url>,<text>,<page_title>`. Can be directly loaded to Pandas. Feel free to use.

## Design

![design](/graphes/design.jpg)

## How it works

![diagram](/graphes/diagram.jpg)
