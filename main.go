package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/signal"

	"time"

	"net/http"

	"strings"

	"strconv"

	"github.com/davecgh/go-spew/spew"
	"github.com/ssor/cloud_collector/conn_parser"
	"github.com/ssor/config"
)

var ()

func main() {
	config_info, err := config.LoadConfig("./conf/config.json")
	if err != nil {
		fmt.Println("[ERR] load config file err: ", err)
		return
	}
	cmds := config_info.Get("cmds")
	if cmds == nil {
		fmt.Println("[ERR] need cmds set")
		return
	}
	metrics := config_info.Get("metrics")
	if metrics == nil {
		fmt.Println("[ERR] need metrics set")
		return
	}
	if len(metrics.([]string)) != len(cmds.([]string)) {
		fmt.Println("[ERR] cmds and metrics not in pairs")
		return
	}
	endPoint := ""
	endPointRaw := config_info.Get("endpoint")
	if (endPointRaw) == nil {
		fmt.Println("[ERR] need endPoint ")
		return
	} else {
		endPoint = endPointRaw.(string)
	}

	taskInterval := 60
	interval := config_info.Get("interval")
	if interval == nil {
		fmt.Println("[ERR] interval not set, will use default 60 seconde")
	} else {
		taskInterval = interval.(int)
	}

	f := func(cmd, metric string) {
		statistics := DoConnStatistics(cmd)
		if statistics != nil {
			PushStatisticsToMonitor(statistics, endPoint, metric)
			// showStatistics(statistics, endPoint, metric)
		}
	}

	metricsList := metrics.([]string)
	for index, cmd := range cmds.([]string) {
		currentCmd := cmd
		currentMetric := metricsList[index]
		go RunTask(func() {
			f(currentCmd, currentMetric)
		}, time.Second*time.Duration(taskInterval))
	}

	fmt.Println("[OK] start task")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// Block until a signal is received.
	<-c
	fmt.Println("[OK] Quit")
}

func truncatePrefix(raw string) (string, string) {
	index := strings.Index(raw, ":::")
	if index < 0 {
		return "", ""
	}
	return raw[:index], raw[index+3:]
}

func DoConnStatistics(cmd string) map[string]int {
	fmt.Println("[OK] start statistics for ", cmd)
	var statistics map[string]int

	out, err := exec.Command(cmd).Output()
	if err != nil {
		fmt.Println("[ERR] Command err: ", err)
		return nil
	}
	// if out != nil {
	// 	if len(out) > 500 {
	// 		fmt.Println("[OK] ", string(out)[:500])
	// 	} else {
	// 		fmt.Println("[OK] ", string(out))
	// 	}
	// }

	raw := string(out)
	prefix, left := truncatePrefix(raw)
	switch prefix {
	case "netstat":
		statistics = doStatisticsConnectToMongo([]byte(left))
	case "mongostat":
		statistics = doStatisticsOfMongoConn(left)
	}
	return statistics

}

func doStatisticsOfMongoConn(raw string) map[string]int {
	hostCounts := strings.Split(raw, "|")
	// spew.Dump(hostCounts)
	if len(hostCounts) <= 0 {
		spew.Dump(raw)
		return nil
	}
	statistics := make(map[string]int)
	for _, hostCount := range hostCounts {
		host, count, err := splitHostAndCount(hostCount)
		if err != nil {
			fmt.Println("[ERR] ", err)
			continue
		}
		statistics[host] = count
	}
	return statistics
}

func splitHostAndCount(raw string) (string, int, error) {
	list := strings.Split(raw, "->")
	if len(list) < 2 {
		spew.Dump(raw)
		return "", -1, errors.New("data format error")
	}

	count, err := strconv.Atoi(list[1])
	if err != nil {
		fmt.Println("[Tip] not number for ", list[1])
		spew.Dump(raw)
		return "", -1, errors.New("data format error")
	}
	return list[0], count, nil
}

func doStatisticsConnectToMongo(raw []byte) map[string]int {
	connections, err := conn_parser.Parse(raw)
	if err != nil {
		fmt.Println("[ERR] parse data err: ", err)
		return nil
	}

	isConnectingToMongo := func(port interface{}) bool {
		if port == nil {
			return false
		}
		return port == "27017"
	}
	return conn_parser.NewConnectionTree(isConnectingToMongo).SortToTree(connections).ConnStatistics()
}

func showStatistics(statistics map[string]int, endPoint, metricPrefix string) {
	if statistics == nil {
		fmt.Println("[Tip] no statistics to push")
		return
	}

	fmt.Println("endpoint: ", endPoint)

	for key, count := range statistics {
		fmt.Println("metric : ", metricPrefix+key, " -> ", count)
	}
}
func PushStatisticsToMonitor(statistics map[string]int, endPoint, metricPrefix string) {
	if statistics == nil {
		fmt.Println("[Tip] no statistics to push")
		return
	}
	now := time.Now()
	fmt.Println("*********** result (", now.Format(time.RFC3339), "): *************")
	messages := []*FalconMessage{}
	timestamp := int(now.Unix())
	for key, count := range statistics {
		fmt.Println(": ", key, " -> ", count)

		msg := New_FalconMessage(endPoint, metricPrefix+key, timestamp, 60, count)
		messages = append(messages, msg)
	}

	json_bs, err := json.Marshal(messages)
	if err != nil {
		fmt.Println("[ERR] marshal err: ", err)
		spew.Dump(messages)
		return
	}

	contentReader := bytes.NewReader(json_bs)
	req, _ := http.NewRequest("POST", "http://127.0.0.1:1988/v1/push", contentReader)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("[ERR] Post data err: ", err)
		fmt.Println(string(json_bs))
		return
	}
	if resp.StatusCode == http.StatusOK {
		fmt.Println("[OK] post  success")
	} else {
		fmt.Println("[ERR] post resp: ")
		fmt.Println(string(json_bs))
		spew.Dump(resp)
	}

	fmt.Println("******************************************************************")
}

func RunTask(f func(), duration time.Duration) {
	if f == nil {
		return
	}

	ticker := time.NewTicker(duration)
	for {
		select {
		case <-ticker.C:
			f()
		}
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
