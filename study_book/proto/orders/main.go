package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/json-iterator/go"
)

var Mode string
var pcs_done = make(chan bool, 1)

const (
	Aidx      = "833"   // 진로그린마트 선부점
	Ver       = "1.3.5" // Application version
	TermHours = 1       // Time interval terms(time.Duration)
	PrcsTerm  = 5       // POS Process observer interval terms(time.Duration)
	Ilimit    = 50      // Data record partition
	//DestinUrl = "https://happygagae.com/syncpos/"
	DestinUrl       = "http://192.168.0.50/syncpos/"
	PingAlarm       = "bridges/register"
	EpOrdered       = "posys/regiReceipt"
	EpUpdateOrdered = "posys/updateReceipt"
	Contype         = "application/json"
)

type OrderedData struct {
	App_idx    string  `json:"aidx"`
	Ord_head   int     `json:"dhead"`
	Ord_code   string  `json:"mcode"`
	Ord_key    string  `json:"ordkey"`
	Ord_type   int     `json:"ordtype"`
	Ord_num    int     `json:"ordnum"`
	Ord_info   string  `json:"ordinf"`
	Pay_code   string  `json:"paycd"`
	Pay_gprice float64 `json:"buyprice"`
	Pay_gqty   int     `json:"buyqty"`
	Pay_gtotal int     `json:"buysum"`
	Ord_date   string  `json:"ord_dt"`
}

type MyError struct {
	When time.Time
	What string
}

var db *sql.DB
var err error

func Dbconn() (db *sql.DB, err error) {
	//host := "192.168.10.201"
	host := "localhost"
	port := "3306"
	dbname := "myposys_db_833"
	uid := "mssaf"
	pass := "trian@akxmflej"

	// sql.DB 객체 생성
	db, err = sql.Open("mysql", uid+":"+pass+"@tcp("+host+":"+port+")/"+dbname)
	return
}

func (e MyError) Error() string {
	return fmt.Sprintf("%v: %v\n", e.When, e.What)
}

func RecordErrorLog(err error) {
	ymd := strings.Replace(time.Now().String()[:10], "-", "", -1)
	fp, _ := os.OpenFile("log-"+ymd+".dtm", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	defer fp.Close()
	// 표준로거를 파일로거로 변경
	log.SetOutput(fp)
	log.Printf("%v", err)
}

func MarkTransmitAlarm(aidx string, m string, s string, v string) {
	resp, err := http.Get(DestinUrl + PingAlarm + "?aidx=" + aidx + "&posit=" + m + "&stat=" + s + "&ver=" + v)
	if err != nil {
		RecordErrorLog(err)
		MarkTransmitAlarm(aidx, "MarkTransmitAlarm", "E", v)
	}

	defer func() {
		if resp.Body != nil {
			resp.Body.Close()
		}
	}()

	// Response 체크.
	respBody, err := ioutil.ReadAll(resp.Body)
	if err == nil {
		str := string(respBody)
		println(str)
	}
}

// POSYS 데이터베이스 접속 테스트
func DaoConnectionPing(db *sql.DB) (cnt int) {
	err := db.QueryRow("SELECT COUNT(*) AS rcnt FROM acust WHERE cust_use = 1").Scan(&cnt)
	if err != nil {
		RecordErrorLog(err)
	}
	return
}

// 회원구매 이력 데이터 Count(3일 이내 데이터)
func TrnOrderedRecordCount(aidx string, v string, db *sql.DB) (tcnt []int, darr []string) {
	cnt := 0
	date := ""
	for i := 0; i < 3; i++ {
		var buf bytes.Buffer
		buf.WriteString(`SELECT COUNT(*), CONVERT(DATE_FORMAT(DATE_SUB(DATE(NOW()), INTERVAL ? DAY), '%y-%m-%d'), CHAR(8))
			FROM zmbpan AS N
				INNER JOIN smember AS M ON M.member_code = N.mbpan_member_code
				INNER JOIN zpanme AS P ON P.panme_jum = N.mbpan_jum AND P.panme_key_date = N.mbpan_key_date AND P.panme_pos_no = N.mbpan_pos_no AND P.panme_junpo_no = N.mbpan_junpo_no
				WHERE P.panme_key_date = CONVERT(DATE_FORMAT(DATE_SUB(DATE(NOW()), INTERVAL ? DAY), '%y-%m-%d'), CHAR(8))
					AND P.panme_dam_name != '' AND P.panme_rec_type < 6`)

		err := db.QueryRow(buf.String(), i, i).Scan(&cnt, &date)
		defer runtime.GC()

		if err != nil {
			RecordErrorLog(err)
			MarkTransmitAlarm(aidx, "OrderRCnt", "E", v)
			return
		}
		tcnt = append(tcnt, cnt)
		darr = append(darr, date)
	}
	return
}

// 회원구매 이력 데이터 Transmit(3일 이내 데이터)
func TrnOrderedRecord(aidx string, v string, db *sql.DB, strDate string, istart int, ilimit int) int {

	var dHead, recType, seqNo, pmQty, pmTotal int
	var mbCode, ordCd, pmField, pmCode, ordDate string
	var pmPrice float64
	var res int = 0

	var buf bytes.Buffer
	buf.WriteString(`SELECT
		IF(P.panme_rec_type = 1 AND P.panme_seq_no = 1, 1, 10) AS header,
		N.mbpan_member_code AS mem_bcode,
		CONCAT(P.panme_key_date, '-', P.panme_pos_no, '-', LPAD(P.panme_junpo_no, 5, '0')) AS ord_code,
		P.panme_rec_type,
		P.panme_seq_no,
		IF(P.panme_dam_name = N.mbpan_member_code, '', P.panme_dam_name) AS  panme_dam_name,
		IF(P.panme_code = N.mbpan_member_code, '', P.panme_code) AS panme_code,
		P.panme_num1,
		P.panme_num2,
		P.panme_num3,
		IF(ASCII(SUBSTRING(P.panme_time, 1,2))>47 && ASCII(SUBSTRING(P.panme_time, 1,2))<58, CONCAT(DATE(P.panme_key_date), ' ', P.panme_time), '') AS ord_date
	FROM zmbpan AS N
		INNER JOIN smember AS M ON M.member_code = N.mbpan_member_code
		INNER JOIN zpanme AS P ON P.panme_jum = N.mbpan_jum AND P.panme_key_date = N.mbpan_key_date AND P.panme_pos_no = N.mbpan_pos_no AND P.panme_junpo_no = N.mbpan_junpo_no
	WHERE P.panme_key_date = ? AND P.panme_dam_name != '' AND P.panme_rec_type < 6 LIMIT ?, ?`)

	rows, err := db.Query(buf.String(), strDate, istart, ilimit)
	defer runtime.GC()

	if err != nil {
		RecordErrorLog(err)
		MarkTransmitAlarm(aidx, "OrderR-query", "E", v)
		return 0
	}
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()

	var OrderedDataArray []*OrderedData
	for rows.Next() {
		res = 0
		err := rows.Scan(&dHead, &mbCode, &ordCd, &recType, &seqNo, &pmField, &pmCode, &pmPrice, &pmQty, &pmTotal, &ordDate)
		if err != nil {
			RecordErrorLog(err)
			MarkTransmitAlarm(aidx, "OrderR-rows", "E", v)
			return 0
		}

		OrderedDataArray = append(OrderedDataArray, &OrderedData{
			App_idx:    aidx,
			Ord_head:   dHead,
			Ord_code:   mbCode,
			Ord_key:    ordCd,
			Ord_type:   recType,
			Ord_num:    seqNo,
			Ord_info:   pmField,
			Pay_code:   pmCode,
			Pay_gprice: pmPrice,
			Pay_gqty:   pmQty,
			Pay_gtotal: pmTotal,
			Ord_date:   ordDate,
		})
	}

	jsonBytes, err := jsoniter.Marshal(&OrderedDataArray)

	if err != nil {
		RecordErrorLog(err)
		MarkTransmitAlarm(aidx, "OrderR-1", "E", v)
		return 0
	}

	buff := bytes.NewBuffer(jsonBytes)
	resp, err := http.Post(DestinUrl+EpOrdered, Contype, buff)
	res = resp.StatusCode
	if err != nil {
		RecordErrorLog(err)
		MarkTransmitAlarm(aidx, "OrderR-2", "E", v)
		return 0
	}

	defer func() {
		if resp.Body != nil {
			resp.Body.Close()
		}
	}()
	// Response 체크.
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		RecordErrorLog(err)
		MarkTransmitAlarm(aidx, "OrderR-3", "E", v)
		return 0
	}
	return res
}

func main() {

	if len(os.Args) == 2 {
		Mode = os.Args[1]
	} else {
		Mode = "day"
	}

	fmt.Println("[", time.Now(), "] Start POS Record Transmitter.(Mode:"+Mode+")")

	switch Mode {
	case "ping":
		db, err := Dbconn()
		if err != nil {
			//RecordErrorLog(err)
			fmt.Println("[", time.Now(), "] fail to connect database.")
		} else {
			if rcnt := DaoConnectionPing(db); rcnt > 0 {
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
		//tsk_obs := time.NewTicker(time.Minute * PrcsTerm)
		tsk_obs := time.NewTicker(time.Minute * 1)
		d_ord := make(chan bool, 1)

		for {
			select {
			case <-tsk_obs.C:
				db, err := Dbconn()
				if err != nil {
					RecordErrorLog(err)
					MarkTransmitAlarm(Aidx, "Impossible", "E", Ver)
				}
				defer func() {
					if db != nil {
						db.Close()
					}
				}()

				if rcnt := DaoConnectionPing(db); rcnt > 0 {
					procs_done <- true
				} else {
					RecordErrorLog(MyError{time.Now(), "I can't access on database."})
				}
			case <-procs_done:
				tsk_obs.Stop()

				// Run main process
				db, err := Dbconn()
				if err != nil {
					RecordErrorLog(err)
					MarkTransmitAlarm(Aidx, "Impossible", "E", Ver)
				}
				defer func() {
					if db != nil {
						db.Close()
					}
				}()

				go func(done chan bool) {
					// 회원구매 이력 데이터 전송
					cnt, dates := TrnOrderedRecordCount(Aidx, Ver, db)
					for i := 0; i < len(cnt); i++ {
						if cnt[i] > 0 {
							block := (cnt[i] / Ilimit) + 1
							for j := 0; j <= block; j++ {
								istart := j * Ilimit
								ilimit := istart + Ilimit
								TrnOrderedRecord(Aidx, Ver, db, dates[i], istart, ilimit)
								time.Sleep(time.Second)
							}
						}
					}
					done <- true
				}(d_ord)
			}
		}
	default:
		RecordErrorLog(MyError{time.Now(), "Unknown worker position.(Mode:" + Mode + ")"})
		return
	}
}
