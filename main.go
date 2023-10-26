package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

	"golang.org/x/exp/slices"
)

type config struct {
	CloudConfig struct {
		ProjectId             string `json:"project-id"`
		GSpreadSheetId        string `json:"gspread-sheet-id"`
		GSpreadInputRange     string `json:"gspread-input-range"`
		GSpreadOutputColStart string `json:"gspread-output-col-start"`
		GSpreadOutputColEnd   string `json:"gspread-output-col-end"`
	} `json:"cloud-config"`

	ImplicitConversions []struct {
		Input  string `json:"input"`
		Output string `json:"output"`
	} `json:"implicit-conversions"`
}

func loadConfig() (*config, error) {
	f, err := os.Open("config.json")
	if err != nil {
		return nil, err
	}
	config := config{}
	err = json.NewDecoder(f).Decode(&config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func main() {
	config, err := loadConfig()
	if err != nil {
		panic(err)
	}

	dict := newStaticDict()

	// スプレッドシートから翻訳対象の日本語を取得
	data := GSheet.getData(config)

	candidates := make([][]string, 0)
	for _, d := range data {

		// 前処理
		preProcess := func(str string) string {
			// 特定の日本語は決められた英単語に変換する
			input := str
			for _, ic := range config.ImplicitConversions {
				input = strings.ReplaceAll(input, ic.Input, fmt.Sprintf("%s ", ic.Output))
			}

			// スペースが二つ以上続いている場合は１つにする
			re := regexp.MustCompile(`\s{2,}`)
			input = re.ReplaceAllString(input, " ")

			return input
		}

		// 前処理の実行
		input := preProcess(d)

		// 静的辞書で変数名を取得
		dictResults := dict.find(input)

		// TranslationAPIで文字列を翻訳する
		translated, err := Api.translateWithCloud(config.CloudConfig.ProjectId, input)
		if err != nil {
			panic(err)
		}

		// 後処理
		postProcess := func(str string) string {
			// 丸括弧を削除
			output := StrUtils.removeBrackets(str)

			// ハイフンをアンダーバーに変換
			output = strings.ReplaceAll(output, "-", "_")

			// 〜'sを空文字に変換
			output = strings.ReplaceAll(output, "'s", "")

			// スラッシュをアンダーバーに変換
			output = strings.ReplaceAll(output, "/", "_")

			// 「数値〜数値」を「数値 to 数値」に変換
			re := regexp.MustCompile(`(\d+)(?:~)(\d+)`)
			output = re.ReplaceAllString(output, "$1 to $2")

			// 文字列の両端に空白が存在した場合はトリムする
			output = strings.Trim(output, " ")

			return output
		}

		// 後処理を実行
		translated = postProcess(translated)

		// 翻訳後文字列を変数名に変換
		translatedVariableName := StrUtils.toSnakeCase(translated)
		fmt.Printf("%s => %s => %s => %s\n", d, input, translated, translatedVariableName)

		variableNameCandidates := append(append([]string{translatedVariableName}, dictResults...), genDisplayNameCandidates(translatedVariableName)...)
		candidates = append(candidates, Dedupe(variableNameCandidates))
	}

	// スプレッドシートに変数名を保存する
	GSheet.UpdateData(config, candidates)
}

// Dedupe は重複を排除した配列を返却します
func Dedupe(array []string) []string {
	results := make([]string, 0)
	for i, v := range array {
		sliced := array[i+1:]
		if slices.Contains(sliced, v) {
			continue
		}
		results = append(results, v)
	}
	return results
}

func genDisplayNameCandidates(translatedVariableName string) []string {
	return []string{}
}
