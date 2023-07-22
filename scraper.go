package crawler

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

func scrapeLinks(url string) ([]Link, []error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, []error{errors.New(fmt.Sprintf("Failed to GET url %s, error %v", url, err))}
	}
	defer resp.Body.Close()

	links, err := parseLinks(resp.Body)
	if err != nil {
		return nil, []error{errors.New(fmt.Sprintf("Failed to get links in HTML of %s, error %v", url, err))}
	}

	domain, err := getDomain(url)
	if err != nil {
		return nil, []error{errors.New(fmt.Sprintf("Failed to get domain of %s, error %v", url, err))}
	}

	var abs_link []Link
	var errs []error
	for _, link := range links {
		if strings.HasPrefix(link.Url, "#") {
			// Skip urls referring to a section.
			continue
		}
		abs_url, err := joinUrl(domain, link.Url)
		if err != nil {
			errs = append(errs, err)
		} else {
			abs_link = append(abs_link, Link{Url: abs_url, Text: link.Text})
		}
	}

	return dedupLinks(abs_link), errs
}

func dedupLinks(links []Link) []Link {
	unique_urls := make(map[string]struct{})

	var unique_links []Link
	for _, link := range links {
		if _, found := unique_urls[link.Url]; !found {
			unique_urls[link.Url] = struct{}{}
			unique_links = append(unique_links, link)
		}
	}

	return unique_links
}
