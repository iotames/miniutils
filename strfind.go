package miniutils

import (
	"regexp"
	"strings"
)

type Strfind struct {
	BodyStr   string
	regStr    string
	matchList [][]string
}

// NewStrfind 字符串提取器。使用正则表达式从字符串中提取信息。
func NewStrfind(bodyStr string) *Strfind {
	return &Strfind{BodyStr: bodyStr}
}

// SetRegexp 设置一个正则表达式。例：匹配数字 `(?s:(\d+))` OR `(\d+)`
func (s *Strfind) SetRegexp(rege string) *Strfind {
	s.regStr = rege
	return s
}

func (s *Strfind) DoFind() *Strfind {
	re := regexp.MustCompile(s.regStr)
	s.matchList = re.FindAllStringSubmatch(s.BodyStr, -1)
	return s
}

func (s *Strfind) GetOne(matchFull bool) string {
	if matchFull {
		return s.matchList[0][0]
	}

	if len(s.matchList) > 0 && len(s.matchList[0]) > 1 {
		return s.matchList[0][1]
	}
	return ""
}

func (s *Strfind) GetAll(matchFull bool) []string {
	var strList []string
	var matchStr string
	for i := 0; i < len(s.matchList); i++ {
		matchStr = s.matchList[i][1]
		if matchFull {
			matchStr = s.matchList[i][0]
		}
		strList = append(strList, matchStr)
	}
	return strList
}

// re = regexp.MustCompile(`<script type=\"application/ld\+json\">(?s:(.+?))</script>`)
func (s *Strfind) SetRegexpBeginEnd(begin string, end string) *Strfind {
	s.SetRegexp(begin + `(?s:(.+?))` + end)
	return s
}

func GetMapStringValue(key string, dictMap map[string]interface{}) string {
	iValue, ok := dictMap[key]
	value := ""
	if ok {
		value = strings.TrimSpace(iValue.(string))
	}
	return value
}

func ReplaceAllString(originalstr, oldstr, newstr string) string {
	return string(regexp.MustCompile(oldstr).ReplaceAll([]byte(originalstr), []byte(newstr)))
}
