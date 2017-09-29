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
	Aidx      = "833"    // 진로그린마트 선부점
	Ver       = "1.3.5a" // Application version
	TermHours = 1        // Time interval terms(time.Duration)
	PrcsTerm  = 5        // POS Process observer interval terms(time.Duration)
	Ilimit    = 50       // Data record partition
	//DestinUrl = "https://happygagae.com/syncpos/"
	DestinUrl   = "http://192.168.0.50/syncpos/"
	PingAlarm   = "bridges/register"
	EpBargain   = "events/regiBargain"
	EpBgoods    = "events/regiProducts"
	EpModifyBgs = "events/modiProducts"
	Contype     = "application/json"
)

type BargainMasterData struct {
	App_idx  string `json:"aidx"`
	Bg_code  int    `json:"code"`
	Bg_title string `json:"title"`
	Bg_rdate string `json:"rdate"`
	Bg_edate string `json:"edate"`
}

type BargainProductsData struct {
	App_idx     string `json:"aidx"`
	Bgg_code    int    `json:"code"`
	Bgg_bcode   string `json:"bcode"`
	Bgg_qty     int    `json:"qty"`
	Bgg_limit   int    `json:"limit"`
	Bgg_hour    int    `json:"hour"`
	Bgg_dccount int    `json:"dccount"`
	Bgg_amount  string `json:"amount"`
	Bgg_price   int    `json:"price"`
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

// 행사정보 Master데이터 Transmit
func TrnBargainMaster(aidx string, v string, db *sql.DB) []int {

	var bg_code int
	var bg_title, bg_rsdate, bg_edate string
	var bgarr []int

	var buf bytes.Buffer
	buf.WriteString(`SELECT
			CONVERT(REPLACE(salecode_code, '-', ''), UNSIGNED INT) AS bg_code,
			salecode_name AS bg_title,
			CONVERT(salecode_ss_date, DATETIME) AS bg_sdate,
			CONCAT(CONVERT(salecode_se_date, DATE), ' 23:59:59') AS bg_edate
		FROM ssalecode`)

	rows, err := db.Query(buf.String())

	if err != nil {
		RecordErrorLog(err)
		MarkTransmitAlarm(aidx, "BargainMasterR-query", "E", v)
	}
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()

	var BargainMasterArray []*BargainMasterData
	for rows.Next() {
		err := rows.Scan(&bg_code, &bg_title, &bg_rsdate, &bg_edate)
		if err != nil {
			RecordErrorLog(err)
			MarkTransmitAlarm(aidx, "BargainMasterR-rows", "E", v)
		}

		// 행사 idx array 저장
		bgarr = append(bgarr, bg_code)

		BargainMasterArray = append(BargainMasterArray, &BargainMasterData{
			App_idx:  aidx,
			Bg_code:  bg_code,
			Bg_title: bg_title,
			Bg_rdate: bg_rsdate,
			Bg_edate: bg_edate,
		})
	}

	jsonBytes, err := jsoniter.Marshal(&BargainMasterArray)

	if err != nil {
		RecordErrorLog(err)
		MarkTransmitAlarm(aidx, "BargainMasterR-1", "E", v)
	}

	buff := bytes.NewBuffer(jsonBytes)
	resp, err := http.Post(DestinUrl+EpBargain, Contype, buff)
	if err != nil {
		RecordErrorLog(err)
		MarkTransmitAlarm(aidx, "BargainMasterR-2", "E", v)
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
		MarkTransmitAlarm(aidx, "BargainMasterR-3", "E", v)
	}
	return bgarr
}

// 행사정보 상품 데이터 Count
func TrnBargainGoodsCount(aidx string, v string, bgidx int, db *sql.DB) (cnt int) {

	if bgidx < 1 {
		return 0
	}

	var buf bytes.Buffer
	buf.WriteString(`SELECT COUNT(S.salepum_sale_code)
			FROM ssalepum AS S
				INNER JOIN ssalecode AS M ON M.salecode_code = S.salepum_sale_code
				INNER JOIN spum AS P ON P.pum_code = S.salepum_pum_code
			WHERE P.pum_ipgo_yn IN ('0', '1')
				AND CONVERT(REPLACE(M.salecode_code, '-', ''), UNSIGNED INT) = ?`)

	err := db.QueryRow(buf.String(), &bgidx).Scan(&cnt)
	defer runtime.GC()

	if err != nil {
		RecordErrorLog(err)
		MarkTransmitAlarm(aidx, "BargainGoodsRCnt", "E", v)
		return 0
	}
	return
}

// 행사정보 상품 데이터 Transmit
func TrnBargainGoodsRecords(aidx string, v string, db *sql.DB, bgidx int, istart int, ilimit int) int {

	var bgg_goods, bgg_qty, bgg_limit, bgg_hour, bgg_dccount, bgg_price int
	var bgg_bcode, bgg_dcamount string
	var res int = 0

	var buf bytes.Buffer
	buf.WriteString(`SELECT
			CONVERT(REPLACE(S.salepum_sale_code, '-', ''), UNSIGNED INT) AS bgg_goods,
			S.salepum_pum_code AS goods_bcode,
			0 AS bgg_qty,
			0 AS bgg_limit,
			0 AS bgg_hour,
			0 AS bgg_dccount,
			S.salepum_sale_wonga AS bgg_dcamount,
			S.salepum_sale_danga AS bgg_sprice
		FROM ssalepum AS S
			INNER JOIN ssalecode AS M ON M.salecode_code = S.salepum_sale_code
			INNER JOIN spum AS P ON P.pum_code = S.salepum_pum_code
		WHERE P.pum_ipgo_yn IN ('0', '1')
			AND CONVERT(REPLACE(M.salecode_code, '-', ''), UNSIGNED INT) = ?
		ORDER BY S.salepum_pum_code LIMIT ?, ?`)

	rows, err := db.Query(buf.String(), bgidx, istart, ilimit)
	defer runtime.GC()

	if err != nil {
		RecordErrorLog(err)
		MarkTransmitAlarm(aidx, "BargainGoodsR-query", "E", v)
		return 0
	}
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()

	var BargainProductsArray []*BargainProductsData
	for rows.Next() {
		res = 0
		err := rows.Scan(&bgg_goods, &bgg_bcode, &bgg_qty, &bgg_limit, &bgg_hour, &bgg_dccount, &bgg_dcamount, &bgg_price)
		if err != nil {
			RecordErrorLog(err)
			MarkTransmitAlarm(aidx, "BargainGoodsR-rows", "E", v)
			return 0
		}

		BargainProductsArray = append(BargainProductsArray, &BargainProductsData{
			App_idx:     aidx,
			Bgg_code:    bgidx,
			Bgg_bcode:   bgg_bcode,
			Bgg_qty:     bgg_qty,
			Bgg_limit:   bgg_limit,
			Bgg_hour:    bgg_hour,
			Bgg_dccount: bgg_dccount,
			Bgg_amount:  bgg_dcamount,
			Bgg_price:   bgg_price,
		})
	}

	jsonBytes, err := jsoniter.Marshal(&BargainProductsArray)

	if err != nil {
		RecordErrorLog(err)
		MarkTransmitAlarm(aidx, "BargainGoodsR-1", "E", v)
		return 0
	}

	buff := bytes.NewBuffer(jsonBytes)
	resp, err := http.Post(DestinUrl+EpBgoods, Contype, buff)
	res = resp.StatusCode
	if err != nil {
		RecordErrorLog(err)
		MarkTransmitAlarm(aidx, "BargainGoodsR-2", "E", v)
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
		MarkTransmitAlarm(aidx, "BargainGoodsR-3", "E", v)
		return 0
	}
	return res
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
					// 행사정보 데이터 전송
					bgarr := TrnBargainMaster(Aidx, Ver, db)
					// 행사정보 상품 데이터 전송
					for i := 0; i < len(bgarr); i++ {
						cnt := TrnBargainGoodsCount(Aidx, Ver, bgarr[i], db)
						block := (cnt / Ilimit) + 1

						fmt.Println("bgarr: ", bgarr[i])
						fmt.Println("cnt: ", cnt)

						for j := 0; j < block; j++ {
							istart := j * Ilimit
							ilimit := istart + Ilimit
							TrnBargainGoodsRecords(Aidx, Ver, db, bgarr[i], istart, ilimit)
							time.Sleep(time.Second)
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
