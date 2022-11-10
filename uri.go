package miniutils

import (
	neturl "net/url"
	"strings"
)

// GetUrlQueryValue 提取网址的查询字符串。即问号之后的部分。
func GetUrlQueryValue(key, url string) string {
	if strings.Contains(url, `%2Fdp%2F`) {
		url, _ = neturl.QueryUnescape(url)
	}
	urlSplit := strings.Split(url, `?`)
	if len(urlSplit) != 2 {
		return ""
	}
	queryString := urlSplit[1]
	querySplit := strings.Split(queryString, `&`)
	qValue := ""
	for _, queryItem := range querySplit {
		itemSplit := strings.Split(queryItem, `=`)
		if itemSplit[0] == key {
			qValue = itemSplit[1]
		}
	}
	return qValue
}

// GetUrl. 获取Url链接全路径。常用于爬虫链接格式化。例: Get("/product?id=90", "https://www.site.com")
// url: startwith http, /, // ; base must startwith http
func GetUrl(url, base string) string {
	if strings.Index(base, "http") != 0 {
		panic("base must start with http. the base is : " + base)
	}
	if strings.Index(url, "http") == 0 {
		return url
	}
	if strings.Index(url, "//") == 0 {
		return strings.Split(base, "/")[0] + url
	}
	if strings.Index(url, "/") == 0 {
		return GetBaseUrl(base) + url
	}
	if url != "" {
		return GetBaseUrl(base) + "/" + url
	}
	return ""
}

// GetDomainByUrl. 获取url网址的域名。
// the arg url startwith http, //, / ; return like: "www.baidu.com", "baidu.com", ""
func GetDomainByUrl(url string) string {
	urlS := strings.Split(url, "/")
	if strings.Index(url, "http") == 0 || strings.Index(url, "//") == 0 {
		return urlS[2]
	}
	if strings.Index(url, "/") == 0 {
		return urlS[1]
	}
	return ""
}

// GetKeywordByDomain GetKeywordByDomain("www.baidu.com") -> baidu
func GetKeywordByDomain(domain string) string {
	domain = strings.ToLower(domain)
	domainSplit := strings.Split(domain, `.`)
	if strings.Contains(domain, "www") {
		return domainSplit[1]
	}
	return domainSplit[0]
}

// GetBaseUrl. GetBaseUrl("https://www.baidu.com/hello?word=hiiii") -> "https://www.baidu.com"
// url must starwith http; return like: https://www.baidu.com
func GetBaseUrl(url string) string {
	if strings.Index(url, "http") != 0 {
		panic("url must start with http. the url is : " + url)
	}
	urlS := strings.Split(url, "/")
	return urlS[0] + "//" + urlS[2]
}
