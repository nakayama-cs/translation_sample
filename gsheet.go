package main

import (
	"context"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type gsheet struct{}

var GSheet gsheet

func (gsheet) getData(config *config) []string {
	credential := option.WithCredentialsFile(config.CloudConfig.GSpreadKeyJson)

	srv, err := sheets.NewService(context.Background(), credential)
	if err != nil {
		panic(err)
	}

	resp, err := srv.Spreadsheets.Values.Get(config.CloudConfig.GSpreadSheetId, config.CloudConfig.GSpreadInputRange).Do()
	if err != nil {
		panic(err)
	}

	if len(resp.Values) == 0 {
		panic("data not found")
	}

	rows := make([]string, 0)
	for _, row := range resp.Values {
		if row, ok := row[0].(string); ok {
			rows = append(rows, row)
		}
	}
	return rows
}

func (gsheet) UpdateData(config *config, data []string) {
	credential := option.WithCredentialsFile(config.CloudConfig.GSpreadKeyJson)

	srv, err := sheets.NewService(context.Background(), credential)
	if err != nil {
		panic(err)
	}

	values := ConvType(data, func(v string) []interface{} {
		return []interface{}{v}
	})

	vr := &sheets.ValueRange{
		Values: values,
	}
	_, err = srv.Spreadsheets.Values.Update(config.CloudConfig.GSpreadSheetId,
		config.CloudConfig.GSpreadOutputRange, vr).ValueInputOption("RAW").Do()
	if err != nil {
		panic(err)
	}
}
