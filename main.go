package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var (
	logPath  string
	interval int
	url      string
	name     string
	version  string
	test     bool
)

func init() {
	flag.IntVar(&interval, "i", 120, "blockheight unincrease time second: default: 120")
	flag.BoolVar(&test, "t", false, "restart docker,just for test: default: false")
	flag.StringVar(&name, "name", "trust-bsc", "restart docker name,default: trust-bsc")
	flag.StringVar(&logPath, "path", "./restart_bsc.log", "log path,default: ./restart_bsc.log")
	flag.StringVar(&url, "url", "http://127.0.0.1:8545", "bsc rpc url,default: http://127.0.0.1:8545")
}

func main() {
	flag.Parse()
	InitLog("info", logPath)
	Logger.Infof("version: %s", version)
	if test {
		RestartBsc(name)
		return
	}
	MonitorBlockIncrease(url, name, interval)
}

/*
url: chain rpc
name: docker name
interval: chain blockheight unincrease time
*/
func MonitorBlockIncrease(url, name string, interval int) {
	var oldH int64
	for {
		newH, err := GetChainHeight(url)
		if err == nil {
			if !(newH > oldH) {
				RestartBsc(name)
			} else {
				oldH = newH
			}
		} else {
			Logger.Warn(fmt.Sprintf("get height err: %s", err.Error()))
		}
		time.Sleep(time.Duration(interval) * time.Second)
	}
}

func GetChainHeight(url string) (int64, error) {
	type Result struct {
		Jsonrpc string `json:"jsonrpc"`
		Id      int    `json:"id"`
		Result  string `json:"result"`
	}
	var res Result
	// url := "http://127.0.0.1:8545"

	payload := strings.NewReader(`{"jsonrpc":"2.0","method":"eth_blockNumber", "id":1}`)

	b, err := post(url, payload)
	if err != nil {
		Logger.Error(err.Error())
		return 0, err
	}
	err = json.Unmarshal(b, &res)
	if err != nil {
		return 0, err
	}

	res.Result = strings.TrimPrefix(res.Result, "0x")

	if res.Result == "" {
		res.Result = "0"
	}

	h, err := strconv.ParseInt(res.Result, 16, 64)
	if err != nil {
		return 0, err
	}
	return h, nil

}

func post(url string, payload *strings.Reader) ([]byte, error) {

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func RestartBsc(name string) {
	Logger.Info("RestartBsc!")
	out, err := exe("docker", "restart", name)
	if err != nil {
		Logger.Warn(err.Error())
	}
	Logger.Infof("RestartBsc Result: %s", out)

}

func exe(name string, arg ...string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	cmd := exec.CommandContext(ctx, name, arg...)
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	err := cmd.Run()
	if err != nil {
		return outb.Bytes(), err
	}
	if len(errb.Bytes()) != 0 {
		return outb.Bytes(), fmt.Errorf(errb.String())
	}
	return outb.Bytes(), nil
}
