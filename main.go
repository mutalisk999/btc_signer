package main

import (
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/kataras/iris/v12"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"syscall"
)

var SecurityPassStr string = "xxkz1&rlje\x00\x00\x00\x00\x00\x00"

func readSecurityPass() ([]byte, error) {
	var fd int
	if terminal.IsTerminal(int(syscall.Stdin)) {
		fd = int(syscall.Stdin)
	} else {
		tty, err := os.Open("/dev/tty")
		if err != nil {
			return nil, errors.New("error allocating terminal")
		}
		defer tty.Close()
		fd = int(tty.Fd())
	}

	SecurityPass, err := terminal.ReadPassword(fd)
	if err != nil {
		return nil, err
	}
	return SecurityPass, nil
}

func LoadConf() error {
	// init config
	jsonParser := new(JsonStruct)
	err := jsonParser.Load("config.json", &GlobalConfig)
	if err != nil {
		fmt.Println("Load config.json", err)
		return err
	}
	return nil
}

func main() {
	//fmt.Printf("Enter Security Password: ")
	//secPassBytes, err := readSecurityPass()
	//if err != nil {
	//	fmt.Println("Read Security Password error: ", err.Error())
	//	os.Exit(-1)
	//}
	//fmt.Println("")
	//
	//fmt.Printf("Enter Security Password (Verify): ")
	//secPassBytes2, err := readSecurityPass()
	//if err != nil {
	//	fmt.Println("Read Security Password (Verify) error: ", err.Error())
	//	os.Exit(-1)
	//}
	//fmt.Println("")
	//
	//if bytes.Compare(secPassBytes, secPassBytes2) != 0 {
	//	fmt.Println("Enter Different Passwords")
	//	os.Exit(-1)
	//}
	//
	//if len(secPassBytes) <= 16 {
	//	paddingCount := 16 - len(secPassBytes)
	//	for i := 0; i < paddingCount; i++ {
	//		secPassBytes = append(secPassBytes, 0x0)
	//	}
	//} else {
	//	secPassBytes = secPassBytes[0:16]
	//}
	//SecurityPassStr = string(secPassBytes)

	iLogFile := "info.log"
	eLogFile := "error.log"
	InitLog(iLogFile, eLogFile, DEBUG)

	err := LoadConf()
	if err != nil {
		Error.Println("LoadConf fail")
		os.Exit(-1)
	}

	err = InitDB(GlobalConfig.DbConfig.DbType, GlobalConfig.DbConfig.DbSource)
	if err != nil {
		Error.Println("InitDB fail")
		os.Exit(-1)
	}

	app = iris.New()
	app.Use(func(ctx iris.Context) {
		ctx.Application().Logger().Infof("Begin request for path: %s", ctx.Path())
		ctx.Next()
	})
	app.Post("/api/wallet/BTC", Controller)
	app.Run(iris.Addr("0.0.0.0:15060"), iris.WithCharset("UTF-8"))
}
