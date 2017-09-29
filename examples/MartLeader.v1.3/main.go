package main

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"martleader.v1.3/app"
	"martleader.v1.3/config"
	"martleader.v1.3/execs"
)

const (
	Aidx      = "833"   // 진로그린마트 선부점
	Ver       = "1.3.5" // Application version
	TermHours = 1       // Time interval terms(time.Duration)
	PrcsTerm  = 5       // POS Process observer interval terms(time.Duration)
	Ilimit    = 50      // Data record partition
)

var Mode string

type MyError struct {
	When time.Time
	What string
}

func (e MyError) Error() string {
	return fmt.Sprintf("%v: %v\n", e.When, e.What)
}

func main() {
	//모든 CPU를 사용하게 함.
	runtime.GOMAXPROCS(runtime.NumCPU())

	if len(os.Args) == 2 {
		Mode = os.Args[1]
	} else {
		Mode = "day"
	}

	fmt.Println("[", time.Now(), "] Start POS Record Transmitter.(Mode:"+Mode+")")

	switch Mode {
	case "ping":
		db, err := config.Dbconn()
		if err != nil {
			app.RecordErrorLog(err)
		} else {
			if rcnt := app.DaoConnectionPing(db); rcnt > 0 {
				fmt.Println("[", time.Now(), "] success to connect database.")
			} else {
				fmt.Println("[", time.Now(), "] fail to connect database.")
			}
		}
		defer func() {
			if db != nil {
				db.Close()
			}
		}()
	case "day":
		procs_done := make(chan bool, 1)
		tsk_obs := time.NewTicker(time.Minute * PrcsTerm)

		for {
			select {
			case <-tsk_obs.C:
				db, err := config.Dbconn()
				if err != nil {
					app.RecordErrorLog(err)
					app.MarkTransmitAlarm(Aidx, "Impossible", "E", Ver)
				}
				defer func() {
					if db != nil {
						db.Close()
					}
				}()

				if rcnt := app.DaoConnectionPing(db); rcnt > 0 {
					procs_done <- true
				} else {
					app.RecordErrorLog(MyError{time.Now(), "I can't access on database."})
				}
			case <-procs_done:
				tsk_obs.Stop()

				// Run main process
				db, err := config.Dbconn()
				if err != nil {
					app.RecordErrorLog(err)
					app.MarkTransmitAlarm(Aidx, "Impossible", "E", Ver)
				}
				defer func() {
					if db != nil {
						db.Close()
					}
				}()

				execs.PtsCronWorker(Aidx, Ver, Ilimit, TermHours)
			}
		}
	case "init":
		procs_done := make(chan bool, 1)
		tsk_obs := time.NewTicker(time.Minute * PrcsTerm)

		for {
			select {
			case <-tsk_obs.C:
				db, err := config.Dbconn()
				if err != nil {
					app.RecordErrorLog(err)
					app.MarkTransmitAlarm(Aidx, "Impossible", "E", Ver)
				}
				defer func() {
					if db != nil {
						db.Close()
					}
				}()

				if rcnt := app.DaoConnectionPing(db); rcnt > 0 {
					procs_done <- true
				} else {
					app.RecordErrorLog(MyError{time.Now(), "I can't access on database."})
				}
			case <-procs_done:
				tsk_obs.Stop()

				// Run main process
				db, err := config.Dbconn()
				if err != nil {
					app.RecordErrorLog(err)
					app.MarkTransmitAlarm(Aidx, "Impossible", "E", Ver)
				}
				defer func() {
					if db != nil {
						db.Close()
					}
				}()

				if <-execs.PtsMasterWorker(Aidx, Ver, Ilimit, db, err) {
					execs.PtsCronWorker(Aidx, Ver, Ilimit, TermHours)
				}
			}
		}
	case "base":
		db, err := config.Dbconn()
		if err != nil {
			app.RecordErrorLog(err)
			app.MarkTransmitAlarm(Aidx, "Impossible", "E", Ver)
		}
		defer func() {
			if db != nil {
				db.Close()
			}
		}()

		execs.PtsMasterWorker(Aidx, Ver, Ilimit, db, err)
	case "sub":
		execs.PtsCronWorker(Aidx, Ver, Ilimit, TermHours)
	default:
		app.RecordErrorLog(MyError{time.Now(), "Unknown worker position.(Mode:" + Mode + ")"})
		return
	}
}
