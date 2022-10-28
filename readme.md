此程序用于修正bsc因分叉而导致的不出块问题的临时解决方案,具体流程为重启本机bsc docker程序.


## 编译
需要`go,make`环境支持
`make`

## 使用说明
```bash
Usage of ./bsc_balance:
  -i int
        blockheight unincrease time second: default: 120 (default 120)
  -name string
        restart docker name,default: trust-bsc (default "trust-bsc")
  -path string
        log path,default: ./restart_bsc.log (default "./restart_bsc.log")
  -url string
        bsc rpc url,default: http://127.0.0.1:8545 (default "http://127.0.0.1:8545")
```