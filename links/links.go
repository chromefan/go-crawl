package links

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
)

var pattern = map[string]string{
	"content_url":       `<span><a href="(.+?)" target="_blank">.+?</a>.+?</span>`,
	"category_url":      `<div class="cont">\s<a([\s\S]+?)"abcd">`,
	"category_text":      `href="(.+?)">([^>]+?)</a>`,
	"article_html":      `>(.+?)</h1>([\s\S]+?)</div>`,
	"article_title":     ">(.+?)</h1>",
	"article_content":   `id="contson[\d\w]+">([\s\S]+?)$`,
	"article_tags_list": `<div class="tag">([\s\S]+?)</div>`,
	"article_tags":      `>([^><]+?)</a>`,
	"article_info_html": `<p class="source">([\s\S.]+?)</a>\s</p>`,
	"article_info":      `>([^<>：].+?)`,
}

func ExtractCategory(url string) (map[string][]string, error) {
	baseUrl := getHost(url)
	doc, err := getHtml(url)
	if err != nil {
		log.Print(err)
	}
	category :=make(map[string][]string)
	var categoryHtml string
	reg := regexp.MustCompile(pattern["category_url"])
	result := reg.FindAllStringSubmatch(doc, -1)

	for _, text := range result {
		categoryHtml = text[1]
		break
	}
	reg = regexp.MustCompile(pattern["category_text"])
	result = reg.FindAllStringSubmatch(categoryHtml, -1)
	for _, text := range result {
		category["url"] = append(category["url"],baseUrl+text[1])
		category["name"] = append(category["name"],text[2])
	}
	fmt.Println(category)
	return category, err
}
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
