package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

func GetDomainStatOptimised(r io.Reader, domain string) (DomainStat, error) {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)

	var (
		result       = make(DomainStat)
		dottedDomain = "." + domain
		err          error
		line         []byte
		user         User
	)

	for scanner.Scan() {
		line = scanner.Bytes()
		err = jsoniter.Unmarshal(line, &user)
		if err != nil {
			return nil, fmt.Errorf("can't unmarshal error: %w, line: %v", err, line)
		}
		if user.Email != "" {
			if strings.HasSuffix(user.Email, dottedDomain) {
				num := result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]
				num++
				result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])] = num
			}
		}
	}

	if err = scanner.Err(); err != nil {
		return nil, fmt.Errorf("line reading error: %w", err)
	}

	return result, nil
}
