package days

import (
	"encoding/json"
	"net/http"
	"time"
)

func Days(writer http.ResponseWriter, request *http.Request) {
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
			writer.Write(failDays("参数错误：start 参数错误！"))
			return
		}
		e, err := time.Parse("2006-01-02", end)
		if err != nil {
			writer.Write(failDays("参数错误：end 参数错误！"))
			return
		}
		if days := calDays(s, e); days == -1 {
			writer.Write(failDays("参数错误：开始日期晚于结束日期！"))
			return
		} else {
			writer.Write(successDays(days, "计算成功！"))
			return
		}
	}
}

type OutputDays struct {
	Success bool   `json:"success"`
	Msg     string `json:"msg"`
	Days    int    `json:"days"`
}

func failDays(msg string) []byte {

	marshal, _ := json.Marshal(&OutputDays{
		false, msg, -1,
	})

	return marshal
}

func successDays(days int, msg string) []byte {
	marshal, _ := json.Marshal(&OutputDays{
		true, msg, days,
	})
	return marshal
}

func calDays(t time.Time, e time.Time) int {

	days := 0
	if e.Before(t) {
		return -1
	}

	for {

		days++

		if t.Equal(e) {
			return days
		}

		t = t.Add(time.Hour * 24)
	}
}