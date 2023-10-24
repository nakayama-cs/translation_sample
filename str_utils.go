package main

import (
	"regexp"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"golang.org/x/text/width"
)

type str_utils struct{}

var StrUtils str_utils

// toUpperCamelCase スペース区切りの英字をアッパーキャメルケース形式の文字列に変換する
func (str_utils) toUpperCamelCase(str string) string {
	return strings.ReplaceAll(cases.Title(language.Und).String(str), " ", "")
}

// toLowerCamelCase スペース区切りの英字をロアーキャメルケース形式の文字列に変換する
func (str_utils) toLowerCamelCase(str string) string {
	upperCamelCase := StrUtils.toUpperCamelCase(str)
	return strings.ToLower(upperCamelCase[0:1]) + upperCamelCase[1:]
}

// toSnakeCaseCase スペース区切りの英字をスネーク形式の文字列に変換する
func (str_utils) toSnakeCaseCase(str string) string {
	return strings.ReplaceAll(strings.ToLower(str), " ", "_")
}

func (str_utils) removeBrackets(str string) string {
	re := regexp.MustCompile(`[()（）]`)
	return re.ReplaceAllString(str, "")
}

// trimRightIndex 単語+数値の文字列が渡された時、単語と数値に分解して返却する
func (str_utils) trimRightIndex(str string) (string, string) {
	re := regexp.MustCompile(`([0-9０１２３４５６７８９]+)$`)
	matchers := re.FindSubmatch([]byte(str))
	if matchers == nil || len(matchers) < 2 {
		return str, ""
	}
	index := string(matchers[1])
	return strings.TrimRight(str, index), width.Fold.String(index)
}
