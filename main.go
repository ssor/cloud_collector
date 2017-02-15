package main

import (
	"fmt"
	"os"
	"os/exec"

	"time"

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
	//netstat -apn | grep ESTABLISHED
	cmd := config_info.Get("cmd").(string)

	do_statistics := func(cmd string) {

		out, err := exec.Command(cmd).Output()
		if err != nil {
			fmt.Println("[ERR] Command err: ", err)
			return
		}
		fmt.Println("[OK] ", string(out)[:500])

		connections, err := parser.Parse(out)
		if err != nil {
			fmt.Println("[ERR] parse data err: ", err)
			return
		}

		statistics := parser.New_MongoConnectionTree().SortToTree(connections).ConnStatistics()
		fmt.Println("*********** result: *************")
		for key, count := range statistics {
			fmt.Println("conn: ", key, " -> ", count)
		}
		fmt.Println("*********************************")
	}
	do_statistics(cmd)
	// go RunTask(do_statistics, time.Second*30)
	// fmt.Printf("The date is %s\n", out)
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
