package handle

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func WorkDay(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	writer.Header().Set("Access-Control-Allow-Headers", "*")

	query := request.URL.Query()
	dates := query["date"]
	if len(dates) < 1 {
		writer.Write(failWorkDay("参数错误"))
		return
	}
	var queryDates []time.Time
	yearsSet := make(map[int]bool)
	for _, d := range dates {
		t, err := time.Parse("2006-01-02", d)
		if err != nil {
			writer.Write(failWorkDay("错误参数：" + d))
			return
		}
		queryDates = append(queryDates, t)
		yearsSet[t.Year()] = true
	}
	var transfer []string
	var holiday []string
	for k := range yearsSet {
		var result *ResultWorkDay
		for i := 0; i < 5; i++ { // 重试 5 次
			var err error
			result, err = getFromUrlWorkDay(k)
			if err == nil {
				break
			}
		}

		if result == nil {
			writer.Write(failWorkDay("数据错误：" + strconv.Itoa(k) + " 年假期数据缺失！"))
			return
		}

		transfer = append(transfer, result.Transfer...)
		holiday = append(holiday, result.Holidays...)
	}

	transfersWorkDay = buildMapWorkDay(transfer)
	holidaysWorkDay = buildMapWorkDay(holiday)

	var result []bool

	for _, t := range queryDates {
		r := isTransferWorkDay(t) || (!isHolidayWorkDay(t) && !isWeekendWorkDay(t))
		result = append(result, r)
	}

	writer.Write(successWorkDay(result, "查询成功！"))

}

var transfersWorkDay map[string]bool
var holidaysWorkDay map[string]bool

type OutputWorkDay struct {
	Success bool   `json:"success"`
	Msg     string `json:"msg"`
	Days    []bool    `json:"days"`
}

type ResultWorkDay struct {
	Year     int      `json:"year"`
	Holidays []string `json:"holidays"`
	Transfer []string `json:"transfer"`
}

func failWorkDay(msg string) []byte {

	marshal, _ := json.Marshal(&OutputWorkDay{
		false, msg, nil,
	})

	return marshal
}

func successWorkDay(days []bool, msg string) []byte {
	marshal, _ := json.Marshal(&OutputWorkDay{
		true, msg, days,
	})
	return marshal
}

func getFromUrlWorkDay(year int) (*ResultWorkDay, error) {
	resp, err := http.Get(fmt.Sprintf("https://raw.githubusercontent.com/"+
		"yzsj98/api/main/holidayJsons/%d.json", year))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var result ResultWorkDay

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	respByte := buf.Bytes()

	if err = json.Unmarshal(respByte, &result); err != nil {
		return nil, err
	}
	return &result, nil

}

func buildMapWorkDay(in []string) map[string]bool {
	m := make(map[string]bool)
	for _, s := range in {
		m[s] = true
	}
	return m

}
func isWeekendWorkDay(t time.Time) bool {
	weekday := t.Weekday()
	return weekday == time.Saturday || weekday == time.Sunday
}

func isTransferWorkDay(t time.Time) bool {
	format := t.Format("2006-01-02")
	_, existed := transfersWorkDay[format]
	return existed
}

func isHolidayWorkDay(t time.Time) bool {
	format := t.Format("2006-01-02")
	_, existed := holidaysWorkDay[format]
	return existed
}