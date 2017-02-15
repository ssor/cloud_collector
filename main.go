package main

import (
	"flag"
	"os"

	"strings"

	"fmt"

	"gopkg.in/gomail.v2"
)

var (
	attach_file = flag.String("a", "", "attached file, only one")
	to_address  = flag.String("to", "zhang", "Address to mail to. If no host is set, @mfer.com.cn will be appended. eg. zhangsan will be zhangsan@mfer.com.cn")
	body        = flag.String("body", "", "mail body, default is empty")
)

func main() {
	flag.Parse()
	if flag.Parsed() == false {
		flag.PrintDefaults()
		return
	}

	if len(*to_address) > 0 {
		if strings.Contains(*to_address, "@") == false {
			*to_address = *to_address + "@mfer.com.cn"
			fmt.Println("[OK] target address will be ", *to_address)
		}
	}

	m := gomail.NewMessage()
	m.SetHeader("From", "zhoubao@mfer.com.cn")
	m.SetHeader("To", *to_address)
	// m.SetAddressHeader("Cc", "dan@example.com", "Dan")
	m.SetHeader("Subject", "New File Arrived")
	m.SetBody("text/plain", *body)
	if len(*attach_file) > 0 {
		if IsFileExist(*attach_file) == false {
			fmt.Println("[ERR] attached file no found")
			return
		}
		m.Attach(*attach_file)
	}

	d := gomail.NewDialer("smtp.mxhichina.com", 587, "zhoubao@mfer.com.cn", "Xsb123456")

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		fmt.Println("[ERR] send email err: ", err)
		return
	}
	fmt.Println("[OK] Send email OK")
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
