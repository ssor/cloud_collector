package parser

import (
	"fmt"
	"strings"
)

type StatisticsResult map[string]int

// Parse accept raw data and parse it to relative statistic results
func Parse(raw string) (StatisticsResult, error) {
	prefix, left := truncatePrefix(raw)
	switch prefix {
	case "netstat":
		return doStatisticsConnectToMongo(left)
	case "mongostat":
		return doStatisticsOfMongoConn(left)
	}
	return nil, fmt.Errorf("no protocol support")
}

func truncatePrefix(raw string) (string, string) {
	index := strings.Index(raw, ":::")
	if index < 0 {
		return "", ""
	}
	return raw[:index], raw[index+3:]
}

func doStatisticsConnectToMongo(raw string) (StatisticsResult, error) {
	connections, err := ParseConnections(raw)
	if err != nil {
		fmt.Println("[ERR] parse data err: ", err)
		return nil, err
	}

	return NewConnectionTree(IsConnectingToMongo).SortToTree(connections).ConnStatistics(), nil
}
