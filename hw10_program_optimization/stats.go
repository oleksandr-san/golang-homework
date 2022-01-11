package hw10programoptimization

import (
	"bufio"
	"bytes"
	"io"
	"strings"

	"github.com/valyala/fastjson"
)

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	suffix := []byte(".")
	suffix = append(suffix, []byte(domain)...)

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

		email := value.GetStringBytes("Email")
		if email != nil && bytes.HasSuffix(email, suffix) {
			domain := strings.ToLower(strings.SplitN(string(email), "@", 2)[1])
			domainStat[domain]++
		}
	}

	return domainStat, nil
}
