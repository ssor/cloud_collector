package conn_parser

import (
	"io/ioutil"
	"os/exec"
	"strings"
	"testing"
)

func TestShell(t *testing.T) {

	out, err := exec.Command("../shell/macos_netstat.sh").Output()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(out))
	raw := string(out)
	index := strings.Index(raw, ":::")
	if index < 0 {
		t.Fatal("should has prefix")
	}
	prefix := raw[:index]
	if prefix != "netstat" {
		t.Fatalf("prefix is %s,  != netstat", prefix)
	}
}

func TestParseLog(t *testing.T) {
	raw, err := ioutil.ReadFile("testdata/mongo_conn.log")
	if err != nil {
		t.Log("read file error: ", err)
		t.FailNow()
	}

	connections, err := Parse(raw)
	if err != nil {
		t.Fatal(err)
	}

	if len(connections) != 8 { // 10 lines in example.log
		t.Fatalf("there should be %d connections", 10)
	}

	isConnectingToMongo := func(port interface{}) bool {
		return port == "27017"
	}
	statistics := NewConnectionTree(isConnectingToMongo).SortToTree(connections).ConnStatistics()
	if len(statistics) != 3 {
		t.Fatalf("there should be %d programs", 3)
	}
	if count, ok := statistics["xsbexam"]; ok == false {
		t.Fatalf("there should be a xsbexam")
	} else {
		if count != 1 {
			t.Fatalf("there should be 1 connections for xsbexam_linux")
		}
	}

	if count, ok := statistics["xsbaccount"]; ok == false {
		t.Fatalf("there should be a xsbaccount_li")
	} else {
		if count != 6 {
			t.Fatalf("there should be 6 connections for xsbaccount")
		}
	}
}
