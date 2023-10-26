package main

import (
	"context"
	"fmt"
	"strings"

	translate "cloud.google.com/go/translate/apiv3"
	"cloud.google.com/go/translate/apiv3/translatepb"
)

type api struct{}

var Api api

func (api) translateWithCloud(projectID string, text string) (string, error) {
	// 入力文字列の前処理 (単語+数字の文字列を翻訳する場合、数字によって翻訳結果が異なるため事前に分離しておく)
	words, index := StrUtils.trimRightIndex(text)

	// クライアントの準備
	ctx := context.Background()
	client, err := translate.NewTranslationClient(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to NewTranslationClient: %w", err)
	}
	defer client.Close()

	// リクエストオブジェクトの準備
	req := &translatepb.TranslateTextRequest{
		Parent:             fmt.Sprintf("projects/%s/locations/global", projectID),
		SourceLanguageCode: "ja",
		TargetLanguageCode: "en-US",
		MimeType:           "text/plain",
		Contents:           []string{words},
	}

	// テキストの翻訳を行う
	resp, err := client.TranslateText(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to TranslateText: %w", err)
	}

	// 翻訳結果が複数行で返却されたとき、単一行に変換する
	lines := make([]string, 0)
	for _, t := range resp.GetTranslations() {
		lines = append(lines, t.GetTranslatedText())
	}

	// 事前にトリムしておいたインデックスを文字列右側に付与して終了
	return strings.Join(lines, " ") + index, nil
}
