package hw10programoptimization

import (
	"bufio"
	"io"
	"regexp"
	"strings"

	"github.com/valyala/fastjson"
)

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	re, err := regexp.Compile("\\." + domain)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(r)
	var parser fastjson.Parser

	domainStat := make(DomainStat)
	for scanner.Scan() {
		content := scanner.Bytes()
		if len(content) == 0 {
			continue
		}

		value, err := parser.ParseBytes(scanner.Bytes())
		if err != nil {
			return nil, err
		}

		if email := value.GetStringBytes("Email"); email != nil && re.Match(email) {
			domain := strings.ToLower(strings.SplitN(string(email), "@", 2)[1])
			domainStat[domain]++
		}
	}

	return domainStat, nil
}
