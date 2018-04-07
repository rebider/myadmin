package utils

import (
	"github.com/astaxie/beego"
	"os/exec"
	"bytes"
	"os"
	"github.com/astaxie/beego/logs"
)

func CheckError(err error, msg... string) {
	if err != nil {
		//_, file, line, _ := runtime.Caller(1)
		//fileBaseName := filepath.Base(file)
		logs.GetBeeLogger().Error("%s %v", msg, err)
		//fmt.Printf("[ERROR]%s:%d %s %v", fileBaseName, line, msg, err)
	}
}

func Nodetool(arg ... string) (string, error) {
	centerNode := beego.AppConfig.String("node::center")
	cookie := beego.AppConfig.String("node::cookie")
	commandArgs := []string{
		"nodetool",
		"-name",
		centerNode,
		"-setcookie",
		cookie,
		"rpc",
	}
	for _, v := range arg {
		commandArgs = append(commandArgs, v)
	}
	out, err := Cmd("escript", commandArgs)
	return out, err
}


func Cmd(commandName string, params []string) (string, error) {
	cmd := exec.Command(commandName, params...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		panic(err)
		//log.Fatal(err)
	}

	err = cmd.Wait()
	return out.String(), err
}

func RemoveDuplicateArray(s [] interface{}) [] interface{} {
	maps := make(map[interface{}]interface{}, len(s))
	r := make([] interface{}, 0)
	for _,v := range s {
		if _, ok := maps[v]; ok {
			continue
		}
		maps[v] = true
		r = append(r, v)
	}
	return r
}
