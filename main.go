package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type config struct {
	CloudConfig struct {
		ProjectId          string `json:"project-id"`
		GSpreadKeyJson     string `json:"gspread-key-json"`
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

		// 特定の日本語は決められた英単語に変換する
		input := d
		for _, ic := range config.ImplicitConversions {
			input = strings.ReplaceAll(input, ic.Input, ic.Output)
		}

		// TranslationAPIで文字列を翻訳する
		translated, err := Api.translationWithCloud(config.CloudConfig.ProjectId, input)
		if err != nil {
			panic(err)
		}

		// 翻訳後文字列を変数名に変換
		variableName := StrUtils.toUpperCamelCase(translated)
		fmt.Printf("%s => %s\n", d, variableName)

		variableNames = append(variableNames, variableName)
	}

	// スプレッドシートに変数名を保存する
	GSheet.UpdateData(config, variableNames)
}
