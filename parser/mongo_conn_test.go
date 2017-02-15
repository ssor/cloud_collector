package parser

import (
	"io/ioutil"
	"testing"
)

func TestParseLog(t *testing.T) {
	raw, err := ioutil.ReadFile("../data/mongo_conn.log")
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

	statistics := New_MongoConnectionTree().SortToTree(connections).ConnStatistics()
	if len(statistics) != 3 {
		t.Fatalf("there should be %d programs", 3)
	}
	if count, ok := statistics["xsbexam_linux"]; ok == false {
		t.Fatalf("there should be a xsbexam_linux")
	} else {
		if count != 1 {
			t.Fatalf("there should be 1 connections for xsbexam_linux")
		}
	}

	if count, ok := statistics["xsbaccount_li"]; ok == false {
		t.Fatalf("there should be a xsbaccount_li")
	} else {
		if count != 6 {
			t.Fatalf("there should be 6 connections for xsbaccount_li")
		}
	}
}
