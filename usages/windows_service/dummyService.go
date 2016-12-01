package main

import (
	"golang.org/x/sys/windows/svc"
	"time"
	"io/ioutil"
	"log"
	"golang.org/x/sys/windows/svc/debug"
)


// 서비스 Type
type dummySvc struct {
}


func (srv *dummySvc) Execute(args []string, req <-chan svc.ChangeRequest, stat chan <- svc.Status) (svcSpecificEC bool, exitCode uint32) {
	stat <- svc.Status{State: svc.StartPending}

	// 실제서비스 내용
	stopChan := make(chan bool, 1)
	go runBody(stopChan)

	stat <- svc.Status{State: svc.Running, Accepts: svc.AcceptStop | svc.AcceptShutdown}

	LOOP:
	for {
		// 서비스 변경 요청에 대해 핸들링
		switch r := <- req; r.Cmd {
		case svc.Stop, svc.Shutdown:
			stopChan <- true
			break LOOP
		case svc.Interrogate:
			stat <- r.CurrentStatus
			time.Sleep(100 * time.Millisecond)
			stat <- r.CurrentStatus
		}
	}

	stat <- svc.Status{State: svc.StopPending}
	return
}

/*** 서비스에서 실제 하는 일 ***/
func runBody( stopChan chan bool) {
	for {
		select {
		case <- stopChan:
			return
		default:
			// 10초마다 현재시간 갱신
		time.Sleep(10 * time.Second)
			ioutil.WriteFile("C:\\Temp\\svc_log.txt", []byte(time.Now().String()), 0)
		}
	}
}


func main() {
	//err := svc.Run("DummySvc", &dummySvc{})
	err := debug.Run("DummySvc", &dummySvc{})	// 디버깅 시 Console 출력
	if err != nil {
		log.Println(err)
	}
}