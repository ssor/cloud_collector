package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"

	"time"

	"net/http"

	"encoding/json"

	"github.com/davecgh/go-spew/spew"
	"github.com/parnurzeal/gorequest"
	"github.com/ssor/cloud_collector/parser"
	"github.com/ssor/config"
)

var ()

func main() {
	config_info, err := config.LoadConfig("./conf/config.json")
	if err != nil {
		fmt.Println("[ERR] load config file err: ", err)
		return
	}
	cmd := config_info.Get("cmd").(string)

	f := func() {
		statistics := DoMongoConnStatistics(cmd)
		if statistics != nil {
			PushStatisticsToMonitor(statistics)
		}
	}
	go RunTask(f, time.Second*60)

	fmt.Println("[OK] start task")
	f() // do one time on start

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// Block until a signal is received.
	<-c
	fmt.Println("[OK] Quit")
}

func DoMongoConnStatistics(cmd string) map[string]int {
	fmt.Println("[OK] start statistics ...")

	out, err := exec.Command(cmd).Output()
	if err != nil {
		fmt.Println("[ERR] Command err: ", err)
		return nil
	}
	fmt.Println("[OK] ", string(out)[:500])

	connections, err := parser.Parse(out)
	if err != nil {
		fmt.Println("[ERR] parse data err: ", err)
		return nil
	}

	statistics := parser.New_MongoConnectionTree().SortToTree(connections).ConnStatistics()
	return statistics

}

func PushStatisticsToMonitor(statistics map[string]int) {
	fmt.Println("*********** result: *************")
	for key, count := range statistics {
		fmt.Println("conn: ", key, " -> ", count)

		msg := New_FalconMessage("www.exam", "conn_mongo_"+key, int(time.Now().Unix()), 60, count)
		json_str, err := json.Marshal(msg)
		if err != nil {
			fmt.Println("[ERR] marshal err: ", err)
			continue
		}

		request := gorequest.New()
		resp, _, errs := request.Post("http://127.0.0.1:1988/v1/push").Send(string(json_str)).End()
		if errs != nil {
			fmt.Println("[ERR] post data to monitor err: ", errs)
			fmt.Println("[TIP]", string(json_str))
			spew.Dump(msg)
			continue
		}

		if resp.StatusCode == http.StatusOK {
			fmt.Println("[OK] post ", key, " : ", count, " success")
		} else {
			fmt.Println("[ERR] post resp: ")
			fmt.Println("[TIP]", string(json_str))
			spew.Dump(resp)
		}
	}
	fmt.Println("*********************************")
}

func RunTask(f func(), duration time.Duration) {
	if f == nil {
		return
	}

	ticker := time.NewTicker(duration)
	for {
		<-ticker.C
		f()
	}
}

// exists returns whether the given file or directory exists or not
func IsFileExist(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}

/*
{
               "endpoint": 机器名（比如www等）,string
               "metric": 指标名称 ,string
               "timestamp": 时间戳, int
               "step": 60, int  (60s上传一次)
               "value": 指标, int
               "counterType": "GAUGE", string (计数器类型 增量/全量)
               "tags": "",(可默认留空字符串)
}

{"endpoint":"www.exam","metric":"conn_mongo_xsbexam_linux","timestamp":1487151349,"step":60,"value":549,"counterType":"GAUGE","tags":""}

*/

type FalconMessage struct {
	EndPoint    string `json:"endpoint"` // www.exam
	Metric      string `json:"metric"`
	Timestamp   int    `json:"timestamp"`
	Step        int    `json:"step"`
	Value       int    `json:"value"`
	CounterType string `json:"counterType"` // GAUGE
	Tags        string `json:"tags"`
}

func New_FalconMessage(endpoint, metric string, timestamp, step, value int) *FalconMessage {
	msg := &FalconMessage{
		EndPoint:    endpoint,
		Metric:      metric,
		Timestamp:   timestamp,
		Step:        step,
		Value:       value,
		CounterType: "GAUGE",
	}
	return msg
}
