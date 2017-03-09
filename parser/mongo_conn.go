package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/davecgh/go-spew/spew"
)

func doStatisticsOfMongoConn(raw string) (StatisticsResult, error) {
	hostCounts := strings.Split(raw, "|")
	if len(hostCounts) <= 0 {
		spew.Dump(raw)
		return nil, fmt.Errorf("data error: %s", raw)
	}
	statistics := make(StatisticsResult)
	for _, hostCount := range hostCounts {
		host, count, err := splitHostAndCount(hostCount)
		if err != nil {
			fmt.Println("[ERR] ", err)
			continue
		}
		statistics[host] = count
	}
	return statistics, nil
}

func splitHostAndCount(raw string) (string, int, error) {
	list := strings.Split(raw, "->")
	if len(list) < 2 {
		return "", -1, fmt.Errorf("data format error: %s", raw)
	}

	count, err := strconv.Atoi(list[1])
	if err != nil {
		return "", -1, fmt.Errorf("not number for %s", list[1])
	}
	return list[0], count, nil
}
