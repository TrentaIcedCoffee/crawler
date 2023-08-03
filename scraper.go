package crawler

import (
	"errors"
	"fmt"
	"net/http"
	net_url "net/url"
	"strings"
)

func FilterOutSectionUrl(links []Link, errors []error) ([]Link, []error) {
	var results []Link
	for _, link := range links {
		if !strings.HasPrefix(link.Url, "#") {
			results = append(results, link)
		}
	}
	return results, errors
}

func MapToAbsUrl(links []Link, errors_ []error, base_url *net_url.URL) ([]Link, []error) {
	var results []Link
	for _, link := range links {
		child_url, err := net_url.Parse(link.Url)
		if err != nil {
			errors_ = append(errors_, errors.New(fmt.Sprintf("Error parsing child url %s, %v", link.Url, err)))
		} else {
			results = append(results, Link{base_url.ResolveReference(child_url).String(), link.Text})
		}
	}
	return results, errors_
}

func KeepOnlySameHostname(links []Link, errors_ []error, base_hostname string) ([]Link, []error) {
	var results []Link
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
	resp, err := http.Get(url)
	if err != nil {
		return nil, []error{errors.New(fmt.Sprintf("Failed to GET url %s, error %v", url, err))}
	}
	defer resp.Body.Close()

	child_links, err := parseLinks(resp.Body)
	if err != nil {
		return nil, []error{errors.New(fmt.Sprintf("Failed to get links in HTML of %s, error %v", url, err))}
	}

	base_url, err := net_url.Parse(url)
	if err != nil {
		return nil, []error{errors.New(fmt.Sprintf("Error parsing base url %s, %v", url, err))}
	}

	var errs []error
	child_links, errs = FilterOutSectionUrl(child_links, errs)
	child_links, errs = MapToAbsUrl(child_links, errs, base_url)
	if keep_only_same_domain {
		child_links, errs = KeepOnlySameHostname(child_links, errs, base_url.Hostname())
	}

	return child_links, errs
}
