package main

import (
	"encoding/csv"
	"log"
	"os"
	"strings"
)

type inOut struct {
	input  string
	output string
}

type staticDict struct {
	values []inOut
}

func newStaticDict() *staticDict {
	f, err := os.Open("./data/jp_en_dictionary.csv")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	r := csv.NewReader(f)
	rows, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	values := make([]inOut, 0)
	for _, cols := range rows {
		if len(cols) < 2 {
			// 列数が不正の場合、そのデータは無視する
			continue
		}
		values = append(values, inOut{
			input:  cols[0],
			output: cols[1],
		})
	}

	return &staticDict{values: values}
}

func (self *staticDict) find(word string) []string {

	preProcess := func(str string) string {
		// 半角の丸括弧は全角に変換
		result := strings.ReplaceAll(str, "(", "（")
		result = strings.ReplaceAll(result, ")", "）")

		// スペースを削除
		result = strings.ReplaceAll(result, " ", "")
		result = strings.ReplaceAll(result, "　", "")
		return result
	}

	// 前処理の実施
	word = preProcess(word)

	results := make([]string, 0)
	for _, v := range self.values {
		if v.input == word {
			results = append(results, v.output)
		}
	}
	return results
}
