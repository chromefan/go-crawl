// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 138.
//!+Extract

// Package links provides a link-extraction function.
package links

import (
	"fmt"
	"golang.org/x/net/html"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
)

var pattern = map[string]string{
	"content_url":       `<span><a href="(.+?)" target="_blank">.+?</a>.+?</span>`,
	"category_url":      `<a href=\".[^\"]+?\">.+?</a>`,
	"article_html":      `>(.+?)</h1>([\s\S]+?)</div>`,
	"article_title":     ">(.+?)</h1>",
	"article_content":   `id="contson[\d\w]+">([\s\S]+?)$`,
	"article_tags_list": `<div class="tag">([\s\S]+?)</div>`,
	"article_tags":      `>([^><]+?)</a>`,
	"article_info_html": `<p class="source">([\s\S.]+?)</a>\s</p>`,
	"article_info":      `>([^<>：].+?)`,
}

// Extract makes an HTTP GET request to the specified URL, parses
// the response as HTML, and returns the links in the HTML document.
func ExtractList(url string) ([]string, error) {
	baseUrl := getHost(url)
	doc, err := getHtml(url)
	var links []string
	reg := regexp.MustCompile(pattern["content_url"])
	result := reg.FindAllStringSubmatch(doc, -1)
	for _, text := range result {
		links = append(links, baseUrl+text[1])
	}
	return links, err
}
func ExtractArticle(url string) (map[string]string, error) {
	doc, err := getHtml(url)
	contents := make(map[string]string)
	if err !=nil {
		return contents,err
	}
	var articleHtml string
	var tagList string
	var info string

	contents["url"] = url

	reg := regexp.MustCompile(pattern["article_html"])
	result := reg.FindAllStringSubmatch(doc, -1)
	for _, text := range result {
		contents["title"] = text[1]
		articleHtml = text[2]
		break
	}
	reg = regexp.MustCompile(pattern["article_content"])
	result = reg.FindAllStringSubmatch(articleHtml, -1)
	for _, text := range result {
		contents["content"] = text[1]
		break
	}
	reg = regexp.MustCompile(pattern["article_tags_list"])
	result = reg.FindAllStringSubmatch(doc, -1)
	for _, text := range result {
		tagList = text[1]
		break
	}
	reg = regexp.MustCompile(pattern["article_tags"])
	result = reg.FindAllStringSubmatch(tagList, -1)
	for _, text := range result {
		contents["tags"] += text[1] + ","
	}
	reg = regexp.MustCompile(pattern["article_info_html"])
	result = reg.FindAllStringSubmatch(doc, -1)
	for _, text := range result {
		info = text[1]
		break
	}
	reg = regexp.MustCompile(pattern["article_info"])
	result = reg.FindAllStringSubmatch(info, -1)
	for i, text := range result {
		if i == 0 {
			contents["dynasty"] = text[1]
		} else {
			contents["author"] = text[1]
		}
	}
	fmt.Println(contents)
	return contents, err
}
func getHtml(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return "", fmt.Errorf("getting %s: %s", url, resp.Status)
	}
	htmlDoc, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return "", fmt.Errorf("parsing %s as HTML: %v", url, err)
	}
	doc := string(htmlDoc)
	return doc, nil
}
func getHost(urlText string) string {
	u, err := url.Parse(urlText)
	if err != nil {
		panic(err)
	}
	baseUrl := u.Scheme + "://" + u.Host
	return baseUrl
}

//!-Extract

// Copied from gopl.io/ch5/outline2.
func forEachNode(n *html.Node, pre, post func(n *html.Node)) {
	if pre != nil {
		pre(n)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		forEachNode(c, pre, post)
	}
	if post != nil {
		post(n)
	}
}
