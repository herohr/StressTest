package main

import (
	"fmt"
	"net"
	"io/ioutil"
	"time"
)

type Result struct {
	Right bool
	TimeConsuming float64
}

var ResultsChan chan Result
var resultsSlice []Result

func ResultsReciver(receiver chan Result, t time.Duration){
	timer := time.NewTimer(t)
	for{
		select{
		case result := <- receiver:
			resultsSlice = append(resultsSlice, result)
		case <- timer.C:
			return

		}
	}
}

func calcBefore(startTime, stopTime time.Time) float64{
	s := stopTime.Second() - startTime.Second()
	ns := stopTime.Nanosecond() - startTime.Nanosecond()
	return float64(s) + float64(ns) / 10e9
}

func connect(raddr string) (*net.TCPConn, error){
	addr, _ := net.ResolveTCPAddr("tcp", raddr)
	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil{
		fmt.Println(err, "at connect ", raddr)
		return nil, nil
	}else{
		return conn, nil
	}
}

func test(raddr string){
	startTime := time.Now()
	conn, err := connect(raddr)
	defer conn.Close()
	_, err = conn.Write([]byte(a))
	if err == nil{
		data, err := ioutil.ReadAll(conn)
		if err == nil{
			stopTime := time.Now()
			resp := string(data)
			//fmt.Println(resp[:])
			if resp[:12] == OK{
				ResultsChan <- Result{true, calcBefore(startTime, stopTime)}
				return
			}else{
				fmt.Println(err)
			}
		}else{
			fmt.Println(err)
		}
	}else{
		fmt.Println(err)
	}
	stopTime := time.Now()
	ResultsChan <- Result{false, calcBefore(startTime, stopTime)}
	return
}

var a string = "GET /login HTTP/1.0\r\n\r\n"
var OK string = "HTTP/1.0 200"
var raddr = "192.168.111.129:8080"

func concurrent_test(c_num int){
	ResultsChan = make(chan Result, 1000)
	resultsSlice = make([]Result, 0)
	go ResultsReciver(ResultsChan, 1*time.Second)
	for i:=0; i< c_num;i++{
		go test(raddr)
	}
	timer := time.NewTimer(1 * time.Second)
	<-timer.C

	completeNum, successNum := 0, 0

	for _, r := range resultsSlice{
		fmt.Println(r)
		completeNum ++
		if r.Right{
			successNum ++
		}
	}
	fmt.Println(c_num, successNum)
}

func main(){
	for i:=100; i<3000;i+=100{
		concurrent_test(i)
	}
}
