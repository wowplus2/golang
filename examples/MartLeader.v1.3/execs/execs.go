package execs

import (
	"database/sql"
	"fmt"
	"runtime"
	"time"

	"martleader.v1.3/app"
	"martleader.v1.3/config"
)

var pcs_done = make(chan bool, 1)

type MyError struct {
	When time.Time
	What string
}

func (e MyError) Error() string {
	return fmt.Sprintf("%v: %v\n", e.When, e.What)
}

func PtsSubObserverOrdered(aidx string, ver string, term int) {

	i_ord := make(chan bool, 1)

	db, err := config.Dbconn()
	if err != nil {
		app.RecordErrorLog(err)
	}
	defer func() {
		if db != nil {
			db.Close()
		}
	}()

	go func(done chan bool) {
		// 회원구매 이력 데이터 전송
		cnt := app.TrnUpdateOrderedRecordCount(aidx, ver, db)
		block := (cnt / term) + 1

		for i := 0; i <= block; i++ {
			istart := i * term
			ilimit := istart + term
			app.TrnUpdateOrderedRecord(aidx, ver, db, istart, ilimit)
			time.Sleep(time.Second)
		}
		done <- true
	}(i_ord)

	if <-i_ord {
		app.MarkTransmitAlarm(aidx, "Order", "M", ver)
	}

	runtime.GC()
}

func PtsMasterWorker(aidx string, ver string, term int, db *sql.DB, err error) chan bool {

	d_ct := make(chan bool, 1)
	d_pd := make(chan bool, 1)
	d_me := make(chan bool, 1)
	d_bg := make(chan bool, 1)
	d_ord := make(chan bool, 1)

	app.MarkTransmitAlarm(aidx, "Master", "S", ver)

	go func(done chan bool) {
		// 판매제품 카테고리 데이터 전송
		cnt := app.TrnCategoryCount(aidx, ver, db)
		block := (cnt / term) + 1

		for i := 0; i < block; i++ {
			istart := i * term
			ilimit := istart + term
			app.TrnCategoryRecord(aidx, ver, db, istart, ilimit)
			time.Sleep(time.Second)
		}
		done <- true
	}(d_ct)

	go func(done chan bool) {
		// 판매제품 마스터 데이터 전송
		cnt := app.TrnMasterRecordCount(aidx, ver, db)
		block := (cnt / term) + 1

		for i := 0; i <= block; i++ {
			istart := i * term
			ilimit := istart + term
			app.TrnMasterRecord(aidx, ver, db, istart, ilimit)
			time.Sleep(time.Second)
		}
		done <- true
	}(d_pd)

	go func(done chan bool) {
		// 회원일반 정보 데이터 전송
		cnt := app.TrnMemberRecordCount(aidx, ver, db)
		block := (cnt / term) + 1

		for i := 0; i <= block; i++ {
			istart := i * term
			ilimit := istart + term
			app.TrnMemberRecord(aidx, ver, db, istart, ilimit)
			time.Sleep(time.Second)
		}
		done <- true
	}(d_me)

	if <-d_ct && <-d_pd && <-d_me {
		app.MarkTransmitAlarm(aidx, "Master", "F", ver)

		go func(done chan bool) {
			// 행사정보 데이터 전송
			bgarr := app.TrnBargainMaster(aidx, ver, db)
			// 행사정보 상품 데이터 전송
			for i := 0; i < len(bgarr); i++ {
				cnt := app.TrnBargainGoodsCount(aidx, ver, bgarr[i], db)
				block := (cnt / term) + 1

				for j := 0; j < block; j++ {
					istart := j * term
					ilimit := istart + term
					app.TrnBargainGoodsRecords(aidx, ver, db, bgarr[i], istart, ilimit)
					time.Sleep(time.Second)
				}
			}
			done <- true
		}(d_bg)

		go func(done chan bool) {
			// 회원구매 이력 데이터 전송
			cnt, dates := app.TrnOrderedRecordCount(aidx, ver, db)
			for i := 0; i < len(cnt); i++ {
				if cnt[i] > 0 {
					block := (cnt[i] / term) + 1
					for j := 0; j <= block; j++ {
						istart := j * term
						ilimit := istart + term
						app.TrnOrderedRecord(aidx, ver, db, dates[i], istart, ilimit)
						time.Sleep(time.Second)
					}
				}
			}
			done <- true
		}(d_ord)

		pcs_done <- true
	}

	runtime.GC()
	return pcs_done
}

func PtsObserver(aidx string, ver string, term int) {
	i_pd := make(chan bool, 1)
	i_me := make(chan bool, 1)
	i_bg := make(chan bool, 1)

	db, err := config.Dbconn()
	if err != nil {
		app.RecordErrorLog(err)
	}
	defer func() {
		if db != nil {
			db.Close()
		}
	}()

	go func(done chan bool) {
		// 판매제품 갱신 데이터 전송
		cnt := app.TrnModifiedRecordCount(aidx, ver, db)
		block := (cnt / term) + 1

		for i := 0; i < block; i++ {
			istart := i * term
			ilimit := istart + term
			app.TrnModifiedRecord(aidx, ver, db, istart, ilimit)
			time.Sleep(time.Second)
		}
		done <- true
	}(i_pd)

	go func(done chan bool) {
		// 회원일반 정보 데이터 전송
		cnt := app.TrnMemberModifiedCount(aidx, ver, db)
		block := (cnt / term) + 1

		for i := 0; i <= block; i++ {
			istart := i * term
			ilimit := istart + term
			app.TrnMemberModifiedRecords(aidx, ver, db, istart, ilimit)
			time.Sleep(time.Second)
		}
		done <- true
	}(i_me)

	go func(done chan bool) {
		// 행사정보 데이터 전송
		app.TrnBargainMaster(aidx, ver, db)
		// 행사제품 갱신 데이터 전송
		cnt := app.TrnBargainModifiedCount(aidx, ver, db)
		block := (cnt / term) + 1

		for i := 0; i < block; i++ {
			istart := i * term
			ilimit := istart + term
			app.TrnBargainModifiedRecords(aidx, ver, db, istart, ilimit)
			time.Sleep(time.Second)
		}
		done <- true
	}(i_bg)

	go func(done chan bool) {
		// 행사정보 데이터 전송
		app.TrnBargainMaster(aidx, ver, db)
		// 행사제품 갱신 데이터 전송
		cnt := app.TrnBargainModifiedCount(aidx, ver, db)
		block := (cnt / term) + 1

		for i := 0; i < block; i++ {
			istart := i * term
			ilimit := istart + term
			app.TrnBargainModifiedRecords(aidx, ver, db, istart, ilimit)
			time.Sleep(time.Second)
		}
		done <- true
	}(i_bg)

	if <-i_pd {
		app.MarkTransmitAlarm(aidx, "Goods", "M", ver)
	}

	if <-i_me {
		app.MarkTransmitAlarm(aidx, "Member", "M", ver)
	}

	if <-i_bg {
		app.MarkTransmitAlarm(aidx, "Bargain", "M", ver)
	}

	PtsSubObserverOrdered(aidx, ver, term)

	runtime.GC()
}

func PtsCronWorker(aidx string, ver string, term int, htimer time.Duration) {
	d := time.Hour * htimer
	t := time.NewTicker(d)
	q := make(chan struct{})

	for {
		select {
		case <-t.C:
			PtsObserver(aidx, ver, term)
		case <-q:
			t.Stop()
			return
		}
	}
}
