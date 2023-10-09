package crawler

import (
	"errors"
	"fmt"
	"net/http"
	net_url "net/url"
	"strings"
)

func filterOutSectionUrl(links []Link, errors []error) ([]Link, []error) {
	results := shortArray[Link]()
	for _, link := range links {
		if !strings.HasPrefix(link.Url, "#") {
			results = append(results, link)
		}
	}
	return results, errors
}

func mapToAbsUrl(links []Link, errors_ []error, base_url *net_url.URL) ([]Link, []error) {
	results := shortArray[Link]()
	for _, link := range links {
		child_url, err := net_url.Parse(link.Url)
		if err != nil {
			errors_ = append(errors_, errors.New(fmt.Sprintf("Error parsing child url %s, %v", link.Url, err)))
		} else {
			results = append(results, Link{Url: base_url.ResolveReference(child_url).String(), Text: link.Text})
		}
	}
	return results, errors_
}

func keepOnlySameHostname(links []Link, errors_ []error, base_hostname string) ([]Link, []error) {
	results := shortArray[Link]()
	for _, link := range links {
		child_url, err := net_url.Parse(link.Url)
		if err != nil {
			errors_ = append(errors_, errors.New(fmt.Sprintf("Error parsing child url %s, %v", link.Url, err)))
			continue
		}
		if base_hostname == child_url.Hostname() {
			results = append(results, link)
		}
	}
	return results, errors_
}

func scrapeLinks(url string, keep_only_same_domain bool) ([]Link, []error) {
	resp, err := requestUrl(url)
	if err != nil {
		return nil, []error{err}
	}
	defer resp.Body.Close()

	child_links, err := parseLinks(resp.Body)
	if err != nil {
		return nil, []error{errors.New(fmt.Sprintf("Failed to parse links in HTML of %s, error %v", url, err))}
	}

	base_url, err := net_url.Parse(url)
	if err != nil {
		return nil, []error{errors.New(fmt.Sprintf("Error parsing base url %s, %v", url, err))}
	}

	errs := shortArray[error]()
	child_links, errs = filterOutSectionUrl(child_links, errs)
	child_links, errs = mapToAbsUrl(child_links, errs, base_url)
	if keep_only_same_domain {
		child_links, errs = keepOnlySameHostname(child_links, errs, base_url.Hostname())
	}

	return child_links, errs
}

func scrapePage(url string) (string, string, error) {
	resp, err := requestUrl(url)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	title, content, err := parsePage(resp.Body)
	if err != nil {
		return "", "", errors.New(fmt.Sprintf("Failed to parse title of url %s, error %v", url, err))
	}

	return title, content, nil
}

func requestUrl(url string) (*http.Response, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to GET url %s, error %v", url, err))
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, errors.New(fmt.Sprintf("Failed to GET url %s with error code %v", url, resp.StatusCode))
	}

	return resp, nil
}
