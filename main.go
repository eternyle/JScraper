package main

import (
        "bufio"
        "fmt"
        "net/http"
        "net/url"
        "os"
        "strings"

        "github.com/PuerkitoBio/goquery"
)

func main() {
        scanner := bufio.NewScanner(os.Stdin)
        for scanner.Scan() {
                urls := extractUrls(scanner.Text())
                for _, u := range urls {
                        jsLinks, err := findJsLinks(u)
                        if err != nil {
                                fmt.Fprintf(os.Stderr, "Error finding .js links in %s: %v\n", u, err)
                                continue
                        }
                        for _, j := range jsLinks {
                                if strings.HasPrefix(j, "//") {
                                        j = "https:" + j
                                }
                                fmt.Println(j)
                        }
                }
        }
        if err := scanner.Err(); err != nil {
                fmt.Fprintln(os.Stderr, "Error reading standard input:", err)
        }
}


func extractUrls(input string) []string {
        var urls []string
        for _, s := range strings.Fields(input) {
                if u, err := url.Parse(s); err == nil && (u.Scheme == "http" || u.Scheme == "https") {
                        urls = append(urls, u.String())
                }
        }
        return urls
}

func findJsLinks(u string) ([]string, error) {
        var jsLinks []string

        client := &http.Client{}
        req, err := http.NewRequest("GET", u, nil)
        if err != nil {
                return nil, err
        }

        resp, err := client.Do(req)
        if err != nil {
                return nil, err
        }
        defer resp.Body.Close()

        if resp.StatusCode != http.StatusOK {
                return nil, fmt.Errorf("response status code is not OK: %v", resp.StatusCode)
        }

        doc, err := goquery.NewDocumentFromReader(resp.Body)
        if err != nil {
                return nil, err
        }

        // Extract .js links from the initial response
        doc.Find("script[src]").Each(func(i int, s *goquery.Selection) {
                if j, ok := s.Attr("src"); ok && strings.HasSuffix(j, ".js") {
                        if strings.HasPrefix(j, "http") || strings.HasPrefix(j, "//") {
                                jsLinks = append(jsLinks, j)
                        } else {
                                jsLinks = append(jsLinks, fmt.Sprintf("%s/%s", u, j))
                        }
                }
        })

        // Extract .js links from subsequent HTTP requests made by the page
        doc.Find("a[href]").Each(func(i int, s *goquery.Selection) {
                href, ok := s.Attr("href")
                if !ok || !strings.HasPrefix(href, "http") {
                        return
                }
                req, err := http.NewRequest("GET", href, nil)
                if err != nil {
                        return
                }
                resp, err := client.Do(req)
                if err != nil {
                        return
                }
                defer resp.Body.Close()
                if resp.StatusCode != http.StatusOK {
                        return
                }
                if strings.HasSuffix(href, ".js") {
                        jsLinks = append(jsLinks, href)
                        return
                }
                doc, err := goquery.NewDocumentFromReader(resp.Body)
                if err != nil {
                        return
                }

                doc.Find("script[src]").Each(func(i int, s *goquery.Selection) {
                        if j, ok := s.Attr("src"); ok && strings.HasSuffix(j, ".js") {
                                if strings.HasPrefix(j, "http") || strings.HasPrefix(j, "//") {
                                        jsLinks = append(jsLinks, j)
                                } else if strings.HasPrefix(j, "/") {
                                        jsLinks = append(jsLinks, fmt.Sprintf("%s%s", u, j))
                                } else {
                                        jsLinks = append(jsLinks, fmt.Sprintf("%s/%s", href, j))
                                }
                        }
                })
        })
        return jsLinks, nil
}
