package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

var transfers map[string]bool
var holidays map[string]bool

type Output struct {
	Success bool   `json:"success"`
	Msg     string `json:"msg"`
	Days    int    `json:"days"`
}

type Result struct {
	Year     int      `json:"year"`
	Holidays []string `json:"holidays"`
	Transfer []string `json:"transfer"`
}

func Handler(writer http.ResponseWriter, request *http.Request) {

	writer.Header().Set("Content-Type", "application/json")
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	writer.Header().Set("Access-Control-Allow-Headers", "*")


	query := request.URL.Query()
	starts, ok1 := query["start"]
	ends, ok2 := query["end"]

	if ok1 && ok2 && len(starts) == 1 && len(ends) == 1 {
		start := starts[0]
		end := ends[0]
		s, err := time.Parse("2006-01-02", start)
		if err != nil {
			writer.Write(fail("参数错误：start 参数错误！"))
			return
		}
		e, err := time.Parse("2006-01-02", end)
		if err != nil {
			writer.Write(fail("参数错误：end 参数错误！"))
			return
		}

		years := make([]int, 0, 2)
		for y:=s.Year(); y<=e.Year(); y++ {
			years = append(years, y)
		}

		if len(years) == 0 {
			writer.Write(fail("参数错误：日期参数错误！"))
			return
		}

		var transfer []string
		var holiday []string

		for _, y := range years {

			var result *Result
			for i := 0; i < 5; i++ { // 重试 5 次
				result, err = getFromUrl(y)
				if err == nil {
					break
				}
			}

			if result == nil {
				writer.Write(fail("数据错误：" + strconv.Itoa(y) + " 年假期数据缺失！"))
				return
			}

			transfer = append(transfer, result.Transfer...)
			holiday = append(holiday, result.Holidays...)
		}



		if err != nil {
			writer.Write(fail("服务器错误"))
			return
		}
		transfers = buildMap(transfer)
		holidays = buildMap(holiday)

		if days := cal(s, e); days == -1 {
			writer.Write(fail("参数错误：开始日期晚于结束日期！"))
		} else {
			writer.Write(success(days, "计算成功！"))
			return
		}

	} else {
		writer.Write(fail("参数错误：缺少参数或者参数冗余。"))
		return
	}

}

func fail(msg string) []byte {

	marshal, _ := json.Marshal(&Output{
		false, msg, -1,
	})

	return marshal
}

func success(days int, msg string) []byte {
	marshal, _ := json.Marshal(&Output{
		true, msg, days,
	})
	return marshal
}

// func getResult(year int) (*Result, error) {
// 	// jsonFile, err := os.Open("holidayJsons/" + strconv.Itoa(year) + ".json")
// 	pwd, _ := os.Getwd()
// 	jsonFile, err := os.Open(path.Join(pwd, "holidayJsons", strconv.Itoa(year)+".json"))
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer jsonFile.Close()

// 	bytes, _ := ioutil.ReadAll(jsonFile)
// 	var result Result
// 	if err = json.Unmarshal(bytes, &result); err != nil {
// 		return nil, err
// 	}
// 	return &result, nil

// }

func getFromUrl(year int) (*Result, error) {
	resp, err := http.Get(fmt.Sprintf("https://raw.githubusercontent.com/"+
		"yzsj98/api/main/holidayJsons/%d.json", year))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var result Result

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	respByte := buf.Bytes()

	if err = json.Unmarshal(respByte, &result); err != nil {
		return nil, err
	}
	return &result, nil

}

func buildMap(in []string) map[string]bool {
	m := make(map[string]bool)
	for _, s := range in {
		m[s] = true
	}
	return m

}

func cal(t time.Time, e time.Time) int {

	days := 0
	if e.Before(t) {
		return -1
	}

	for {

		if isTransfer(t) || (!isHoliday(t) && !isWeekend(t)) {
			days++
		}

		if t.Equal(e) {
			return days
		}

		t = t.Add(time.Hour * 24)
	}
}

func isWeekend(t time.Time) bool {
	weekday := t.Weekday()
	return weekday == time.Saturday || weekday == time.Sunday
}

func isTransfer(t time.Time) bool {
	format := t.Format("2006-01-02")
	_, existed := transfers[format]
	return existed
}

func isHoliday(t time.Time) bool {
	format := t.Format("2006-01-02")
	_, existed := holidays[format]
	return existed
}
