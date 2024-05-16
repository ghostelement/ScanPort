package scan

import (
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

func checkPort(host string, port int, timeout int) bool {
	//fmt.Println(host, port, time.Duration(timeout)*time.Microsecond)
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), time.Duration(timeout)*time.Millisecond)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

var logMutex sync.Mutex // 互斥锁，用于保护日志写入操作

func scanHost(host string, portlist []int, concurrency int, timeout int) {
	//defer wg.Done()
	var localWg sync.WaitGroup
	//localWg.Add(len(portlist))

	pool := make(chan struct{}, concurrency)
	for _, port := range portlist {
		localWg.Add(1)
		pool <- struct{}{}
		go func(port int) {
			defer func() {
				<-pool
				localWg.Done()
			}()
			if available := checkPort(host, port, timeout); available {
				logMutex.Lock()
				fmt.Printf("%s: 端口 %d is open\n", host, port)
				log.Printf("%s: 端口 %d is open\n", host, port)
				logMutex.Unlock()
			}
		}(port)
	}

	localWg.Wait()
}

func (c *Job) Scan() {
	// 初始化端口1~65535
	portlist := make([]int, 65535)
	for i := 0; i < 65535; i++ {
		portlist[i] = i + 1
	}

	// 打开日志文件
	logFile, err := os.OpenFile("scanport.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	timeout := 1000
	if c.Timeout != 0 {
		timeout = c.Timeout
	}
	for _, host := range c.Hosts {
		fmt.Println("开始扫描主机：", host)
		scanHost(host, portlist, c.ParallelNum, timeout)
	}
}
