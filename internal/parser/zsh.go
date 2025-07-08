package parser

import (
	"regexp"
)

type ZshParser struct {
	Path string
}

func ParseLine(line string) (command string, timestamp int64, skip bool) {
	return "", 0, false
}

func NewZshParser(path string) *ZshParser {
	return &ZshParser{Path: path}
}

func (p *ZshParser) GetHistoryPath() []string {
	return []string{p.Path}
}

func stripZshWeirdness(line string) string {
	// Regex: `: \d+:\d;(.*)`
	// Input:  ": 1666062975:0;echo hello"
	// Output: "echo hello"
	re := regexp.MustCompile(`: \d+:\d;(.*)`)
	matches := re.FindStringSubmatch(line)
	if len(matches) == 2 {
		return matches[1]
	}
	return line
}

