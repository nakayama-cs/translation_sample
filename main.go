package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type config struct {
	CloudConfig struct {
		ProjectId          string `json:"project-id"`
		GSpreadSheetId     string `json:"gspread-sheet-id"`
		GSpreadInputRange  string `json:"gspread-input-range"`
		GSpreadOutputRange string `json:"gspread-output-range"`
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

	data := GSheet.getData(config)

	variableNames := make([]string, 0)
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

		// TranslationAPIで文字列を翻訳する
		translated, err := Api.translationWithCloud(config.CloudConfig.ProjectId, input)
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
		variableName := StrUtils.toSnakeCase(translated)
		fmt.Printf("%s => %s => %s => %s\n", d, input, translated, variableName)

		variableNames = append(variableNames, variableName)
	}

	// スプレッドシートに変数名を保存する
	GSheet.UpdateData(config, variableNames)
}
