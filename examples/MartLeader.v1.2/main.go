package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"time"

	"martleader.v1.2/app"
	"martleader.v1.2/config"
	"martleader.v1.2/execs"
)

const (
	Aidx      = "825"   // 광주식자재할인마트 경기광주점
	Ver       = "1.2.5" // Application version
	TermHours = 1       // Time interval terms(time.Duration)
	PrcsTerm  = 5       // POS Process observer interval terms(time.Duration)
	Ilimit    = 50      // Data record partition
	//PosRunner = "sqlservr.exe" // MS-SQL 2008R2 실행파일 명
	PosRunner = "iamMgr.exe" // 투게더POS 실행파일 명
)

var Mode string

func isProcRunning(names ...string) (bool, error) {
	if len(names) == 0 {
		return false, nil
	}

	cmd := exec.Command("TASKLIST.EXE", "/fo", "csv", "/nh")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, err := cmd.Output()
	if err != nil {
		return false, err
	}

	for _, name := range names {
		if bytes.Contains(out, []byte(name)) {
			return true, nil
		}
	}

	return false, nil
}

type MyError struct {
	When time.Time
	What string
}

func (e MyError) Error() string {
	return fmt.Sprintf("%v: %v\n", e.When, e.What)
}

func main() {
	if len(os.Args) == 2 {
		Mode = os.Args[1]
	} else {
		Mode = "all"
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

			chkRun, err := isProcRunning(PosRunner)
			if chkRun {
				fmt.Println("[", time.Now(), "] success.", PosRunner, "is running.")
			} else if err != nil {
				app.RecordErrorLog(err)
				fmt.Println("[", time.Now(), "]", err)
			} else {
				fmt.Println("[", time.Now(), "]", PosRunner, "is not working.")
			}
		}
		defer func() {
			if db != nil {
				db.Close()
			}
		}()
	case "init":
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
	case "day":
		procs_done := make(chan bool, 1)
		tsk_obs := time.NewTicker(time.Minute * PrcsTerm)

		for {
			select {
			case <-tsk_obs.C:
				isRun, err := isProcRunning(PosRunner)
				if isRun {
					procs_done <- true
				} else {
					app.RecordErrorLog(err)
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
	case "all":
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
