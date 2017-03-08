package conn_parser

import (
	"fmt"
	"strings"

	"github.com/davecgh/go-spew/spew"
)

type ActiveInternetConnection struct {
	LocalHost   string
	LocalPort   string
	ForeignHost string
	ForeignPort string
	ProgramName string
}

type ActiveInternetConnectionArray []*ActiveInternetConnection

func Parse(raw []byte) (ActiveInternetConnectionArray, error) {
	connections := ActiveInternetConnectionArray{}

	raw_string := string(raw)
	lines := strings.Split(raw_string, "\n")
	print_debug(func() {
		fmt.Println("lines: ", len(lines))
	})

	for index, line := range lines {
		print_debug(func() {
			fmt.Println(index, " : ", line)
		})
		conn := parseConnection(line)
		if conn != nil {
			connections = append(connections, conn)
		}
	}

	return connections, nil
}

func parseConnection(raw string) *ActiveInternetConnection {
	if len(raw) <= 0 {
		return nil
	}
	items := strings.Split(raw, " ")

	items_no_space := []string{}
	for _, item := range items {
		item_no_space := strings.TrimSpace(item)
		if len(item_no_space) > 0 {
			items_no_space = append(items_no_space, item_no_space)
		}
	}

	if len(items_no_space) < 7 {
		fmt.Println("data format error: ", raw)
		return nil
	}

	print_debug(func() {
		for index, item := range items_no_space {
			fmt.Println(index, " -> ", item)
		}
	})

	if len(items_no_space) < 6 {
		return nil
	}

	localhost_and_port := strings.SplitN(items_no_space[3], ":", 2)
	foreign_host_and_port := strings.SplitN(items_no_space[4], ":", 2)
	program_name := strings.Replace(items_no_space[6], "-", "", 1)
	if len(program_name) > 0 {
		pid_split_index := strings.Index(program_name, "/") // trim 17039/xsbaccount_li to xsbaccount_li
		if pid_split_index > 0 {
			program_name = program_name[pid_split_index+1:]
			underscord_index := strings.Index(program_name, "_")
			if underscord_index > 0 {
				program_name = program_name[:underscord_index]
			}
		}
	}

	conn := &ActiveInternetConnection{
		LocalHost:   localhost_and_port[0],
		LocalPort:   localhost_and_port[1],
		ForeignHost: foreign_host_and_port[0],
		ForeignPort: foreign_host_and_port[1],
		ProgramName: program_name,
	}
	print_debug(func() {
		spew.Dump(conn)
	})
	return conn
}
