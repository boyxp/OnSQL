package tokenizer

import "regexp"
import "strings"

func Tokenize(condition string) []string {
	condition = strings.TrimSpace(condition)
	condition = regexp.MustCompile(`\s*([=!><]+|\(|\)|\sand\s|\sor\s)\s*`).ReplaceAllString(condition, " $1 ")
	condition = regexp.MustCompile(`\sin\s*\(\s*\?\s*\)`).ReplaceAllString(condition, " in ?")
	tokens   := regexp.MustCompile(`\s+`).Split(condition, -1)
	return tokens
}
