package main

import (
	"context"
	"fmt"

	"google.golang.org/api/sheets/v4"
)

type gsheet struct{}

var GSheet gsheet

func (gsheet) getData(config *config) []string {
	srv, err := sheets.NewService(context.Background())
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

func (gsheet) UpdateData(config *config, data [][]string) {
	srv, err := sheets.NewService(context.Background())
	if err != nil {
		panic(err)
	}

	values := ConvType(data, func(v []string) []interface{} {
		result := make([]interface{}, 0)
		for _, vv := range v {
			result = append(result, vv)
		}
		return result
	})

	vr := &sheets.ValueRange{
		Values: values,
	}

	outputRange := fmt.Sprintf("%s:%s", config.CloudConfig.GSpreadOutputColStart, config.CloudConfig.GSpreadOutputColEnd)

	_, err = srv.Spreadsheets.Values.Update(config.CloudConfig.GSpreadSheetId,
		outputRange, vr).ValueInputOption("RAW").Do()
	if err != nil {
		panic(err)
	}
}
