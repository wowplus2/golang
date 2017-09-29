package app

import (
	"bytes"
	"database/sql"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/json-iterator/go"
)

const (
	//DestinUrl = "https://happygagae.com/syncpos/"
	DestinUrl       = "http://192.168.0.50/syncpos/"
	PingAlarm       = "bridges/register"
	EpProduts       = "products/register"
	EpModifyPrds    = "products/modifier"
	EpCategory      = "categories/register"
	EpMembers       = "members/register"
	EpModifyMems    = "members/modifier"
	EpBargain       = "events/regiBargain"
	EpBgoods        = "events/regiProducts"
	EpModifyBgs     = "events/modiProducts"
	EpOrdered       = "receipt/regiReceipt"
	EpUpdateOrdered = "receipt/updateReceipt"
	Contype         = "application/json"
)

type PosCategogyData struct {
	App_idx string `json:"aidx"`
	Gb_top  int    `json:"top"`
	Gb_mid  int    `json:"mid"`
	Gb_bot  int    `json:"bot"`
	Gb_name string `json:"name"`
	Gb_code int    `json:"code"`
}

type PosMasterData struct {
	App_idx    string  `json:"aidx"`
	Cate_code  string  `json:"cate"`
	Gds_code   string  `json:"pcode"`
	Gds_bcode  string  `json:"bcode"`
	Gds_name   string  `json:"name"`
	Gds_bprice float32 `json:"bprice"`
	Gds_sprice int     `json:"sprice"`
}

type MartMemberData struct {
	App_idx    string `json:"aidx"`
	Mem_code   int    `json:"muid"`
	Mem_name   string `json:"name"`
	Mem_bcode  string `json:"bcode"`
	Mem_tel    string `json:"telno"`
	Mem_point  int    `json:"mpt"`
	Mem_tpoint int    `json:"spt"`
	Mem_rdate  string `json:"reg_dt"`
	Mem_odate  string `json:"ord_dt"`
}

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
	Bgg_amount  int    `json:"amount"`
	Bgg_price   int    `json:"price"`
}

type OrderedData struct {
	App_idx         string `json:"aidx"`
	Ord_head        int    `json:"dhead"`
	Mem_code        int    `json:"muid"`
	Mem_bcode       string `json:"mbcode"`
	Ord_code        string `json:"ordcd"`
	Pay_cash_price  int    `json:"paycash"`
	Pay_card_price  int    `json:"paycard"`
	Pay_point_price int    `json:"paypoint"`
	Gds_code        string `json:"pcode"`
	Gds_bcode       string `json:"bcode"`
	Pay_goods_price int    `json:"payprice"`
	Ord_date        string `json:"ord_dt"`
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
	var buf bytes.Buffer
	buf.WriteString(DestinUrl + PingAlarm + "?aidx=" + aidx + "&posit=" + m + "&stat=" + s + "&ver=" + v)

	resp, err := http.Get(buf.String())
	if err != nil {
		RecordErrorLog(err)
		MarkTransmitAlarm(aidx, "Alarm-http", "E", v)
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
		MarkTransmitAlarm(aidx, "Alarm-res", "E", v)
	}
}

// 투게더 포스 데이터베이스 접속 테스트
func DaoConnectionPing(db *sql.DB) (cnt int) {
	var buf bytes.Buffer
	buf.WriteString("SELECT COUNT(id) FROM sysobjects WHERE type = 'U'")

	err := db.QueryRow(buf.String()).Scan(&cnt)
	defer runtime.GC()

	if err != nil {
		RecordErrorLog(err)
	}
	return
}

// 회원구매 이력 데이터 Count(5주 이내 데이터)
func TrnOrderedRecordCount(aidx string, v string, db *sql.DB) (cnt int) {
	var buf bytes.Buffer
	buf.WriteString(`SELECT SUM(MAIN.rcnt) AS rcnt
		FROM (
			SELECT COUNT(R.ord_code) AS rcnt
			FROM Ord AS R
				INNER JOIN member AS M ON M.mem_code = R.ord_memcode
			WHERE M.mem_status = 1
				AND R.ord_memcode <> 0
				AND DATEDIFF(WEEK, R.ord_date, CURRENT_TIMESTAMP) BETWEEN 0 AND 5
			UNION ALL
			SELECT COUNT(R.ord_code) AS rcnt
			FROM Ordd AS P
				INNER JOIN Ord AS R ON P.ordd_code = R.ord_code
				INNER JOIN member AS M ON M.mem_code = R.ord_memcode
			WHERE M.mem_status = 1
				AND R.ord_memcode <> 0
				AND DATEDIFF(WEEK, R.ord_date, CURRENT_TIMESTAMP) BETWEEN 0 AND 5
		) AS MAIN`)

	err := db.QueryRow(buf.String()).Scan(&cnt)
	defer runtime.GC()

	if err != nil {
		RecordErrorLog(err)
		MarkTransmitAlarm(aidx, "OrderRCnt", "E", v)
		return 0
	}
	return
}

// 회원구매 이력 데이터 Transmit(5주 이내 데이터)
func TrnOrderedRecord(aidx string, v string, db *sql.DB, istart int, ilimit int) int {
	var nrows, dhead, muid, paycash, paycard, paypoint, gprice int
	var mbcode, ordcd, gcode, bcode, ord_date string
	var res int = 0

	var buf bytes.Buffer
	buf.WriteString(`SELECT MAIN.*
		FROM (
			SELECT
				ROW_NUMBER() OVER (ORDER BY SUB.ord_date DESC) AS nrows, SUB.*
			FROM(
				SELECT
					1 AS header,
					M.mem_code,
					RTRIM(M.mem_bcode) AS mem_bcode,
					R.ord_code,
					R.ord_paycash,
					R.ord_paycredit,
					R.ord_paypoint,
					0 AS ordd_goods,
					'' AS goods_bcode,
					0 AS ordd_price,
					R.ord_date
				FROM Ord AS R
					INNER JOIN member AS M ON M.mem_code = R.ord_memcode
				WHERE M.mem_status = 1
					AND R.ord_memcode <> 0
					AND DATEDIFF(WEEK, R.ord_date, CURRENT_TIMESTAMP) BETWEEN 0 AND 5
				UNION ALL
				SELECT
					10 AS header,
					M.mem_code,
					RTRIM(M.mem_bcode) AS mem_bcode,
					R.ord_code,
					0 AS ord_paycash,
					0 AS ord_paycredit,
					0 AS ord_paypoint,
					P.ordd_goods,
					(SELECT TOP 1 goods_bcode FROM goods WHERE goods_code = P.ordd_goods) AS goods_bcode,
					P.ordd_price,
					R.ord_date
				FROM Ordd AS P
					INNER JOIN Ord AS R ON P.ordd_code = R.ord_code
					INNER JOIN member AS M ON M.mem_code = R.ord_memcode
				WHERE
					M.mem_status = 1
					AND R.ord_memcode <> 0
					AND DATEDIFF(WEEK, R.ord_date, CURRENT_TIMESTAMP) BETWEEN 0 AND 5
			) AS SUB
		) AS MAIN
		WHERE MAIN.nrows BETWEEN ? AND ?`)

	rows, err := db.Query(buf.String(), istart, ilimit)
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
		err := rows.Scan(&nrows, &dhead, &muid, &mbcode, &ordcd, &paycash, &paycard, &paypoint, &gcode, &bcode, &gprice, &ord_date)
		if err != nil {
			RecordErrorLog(err)
			MarkTransmitAlarm(aidx, "OrderR-rows", "E", v)
			return 0
		}

		OrderedDataArray = append(OrderedDataArray, &OrderedData{
			App_idx:         aidx,
			Ord_head:        dhead,
			Mem_code:        muid,
			Mem_bcode:       mbcode,
			Ord_code:        ordcd,
			Pay_cash_price:  paycash,
			Pay_card_price:  paycard,
			Pay_point_price: paypoint,
			Gds_code:        gcode,
			Gds_bcode:       bcode,
			Pay_goods_price: gprice,
			Ord_date:        ord_date,
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

// 판매제품 카테고리 Record Count Transmit
func TrnCategoryCount(aidx string, v string, db *sql.DB) (cnt int) {
	var buf bytes.Buffer
	buf.WriteString(`SELECT COUNT(gb_top) FROM dbo.gbook WHERE gb_status = 1`)

	err := db.QueryRow(buf.String()).Scan(&cnt)
	defer runtime.GC()

	if err != nil {
		RecordErrorLog(err)
		MarkTransmitAlarm(aidx, "CategoryRCnt", "E", v)
		return 0
	}
	return
}

// 판매제품 카테고리 데이터 Transmit
func TrnCategoryRecord(aidx string, v string, db *sql.DB, istart int, ilimit int) int {
	var nrows, gb_top, gb_mid, gb_bot, gb_code int
	var gb_name string
	var res int = 0

	var buf bytes.Buffer
	buf.WriteString(`SELECT MAIN.*
		FROM(
			SELECT
				ROW_NUMBER() OVER (ORDER BY gb_top ASC) AS nrows,
				gb_top,
				gb_mid,
				gb_bot,
				gb_name,
				gb_code
			FROM dbo.gbook WHERE gb_status = 1
		) AS MAIN
		WHERE MAIN.nrows BETWEEN ? AND ?`)

	rows, err := db.Query(buf.String(), istart, ilimit)
	defer runtime.GC()

	if err != nil {
		RecordErrorLog(err)
		MarkTransmitAlarm(aidx, "CategoryR-query", "E", v)
		return 0
	}
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()

	var PosCategoryArray []*PosCategogyData
	for rows.Next() {
		res = 0
		err := rows.Scan(&nrows, &gb_top, &gb_mid, &gb_bot, &gb_name, &gb_code)
		if err != nil {
			RecordErrorLog(err)
			MarkTransmitAlarm(aidx, "CategoryR-rows", "E", v)
			return 0
		}

		PosCategoryArray = append(PosCategoryArray, &PosCategogyData{
			App_idx: aidx,
			Gb_top:  gb_top,
			Gb_mid:  gb_mid,
			Gb_bot:  gb_bot,
			Gb_name: gb_name,
			Gb_code: gb_code,
		})
	}

	jsonBytes, err := jsoniter.Marshal(&PosCategoryArray)

	if err != nil {
		RecordErrorLog(err)
		MarkTransmitAlarm(aidx, "CategoryR-1", "E", v)
		return 0
	}

	buff := bytes.NewBuffer(jsonBytes)
	resp, err := http.Post(DestinUrl+EpCategory, Contype, buff)
	res = resp.StatusCode
	if err != nil {
		RecordErrorLog(err)
		MarkTransmitAlarm(aidx, "CategoryR-2", "E", v)
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
		MarkTransmitAlarm(aidx, "CategoryR-3", "E", v)
		return 0
	}
	return res
}

// 판매제품 Master데이터 Count
func TrnMasterRecordCount(aidx string, v string, db *sql.DB) (cnt int) {
	var buf bytes.Buffer
	buf.WriteString(`SELECT COUNT(C.gb_code)
		FROM goods AS G
			INNER JOIN gbook AS C ON C.gb_top = G.goods_bktop AND C.gb_mid = G.goods_bkmid AND C.gb_bot = G.goods_bkbot
		WHERE G.goods_status = 1`)

	err := db.QueryRow(buf.String()).Scan(&cnt)
	defer runtime.GC()

	if err != nil {
		RecordErrorLog(err)
		MarkTransmitAlarm(aidx, "GoodsRCnt", "E", v)
		return 0
	}
	return
}

// 판매제품 Master데이터 Transmit
func TrnMasterRecord(aidx string, v string, db *sql.DB, istart int, ilimit int) int {

	var nrows, gds_sprice int
	var gds_bprice float32
	var cate_code, gds_code, gds_bcode, gds_name string
	var res int = 0

	var buf bytes.Buffer
	buf.WriteString(`SELECT *
		FROM (
			SELECT
				ROW_NUMBER() OVER (ORDER BY G.goods_bcode ASC) AS nrows,
				C.gb_code,
				G.goods_code,
				RTRIM(G.goods_bcode) AS gds_bcode,
				(LTRIM(RTRIM(G.goods_name)) + ' ' + RTRIM(G.goods_sspec)) AS gds_name,
				G.goods_bprice,
				G.goods_sprice
			FROM dbo.goods AS G
				INNER JOIN dbo.gbook AS C ON C.gb_top = G.goods_bktop AND C.gb_mid = G.goods_bkmid AND C.gb_bot = G.goods_bkbot
			WHERE G.goods_status = 1
		) AS MAIN
		WHERE MAIN.nrows BETWEEN ? AND ?`)

	rows, err := db.Query(buf.String(), istart, ilimit)
	defer runtime.GC()

	if err != nil {
		RecordErrorLog(err)
		MarkTransmitAlarm(aidx, "GoodsR-query", "E", v)
		return 0
	}
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()

	var PosMasterArray []*PosMasterData
	for rows.Next() {
		res = 0
		err := rows.Scan(&nrows, &cate_code, &gds_code, &gds_bcode, &gds_name, &gds_bprice, &gds_sprice)
		if err != nil {
			RecordErrorLog(err)
			MarkTransmitAlarm(aidx, "GoodsR-rows", "E", v)
			return 0
		}

		PosMasterArray = append(PosMasterArray, &PosMasterData{
			App_idx:    aidx,
			Cate_code:  cate_code,
			Gds_code:   gds_code,
			Gds_bcode:  gds_bcode,
			Gds_name:   gds_name,
			Gds_bprice: gds_bprice,
			Gds_sprice: gds_sprice,
		})
	}

	jsonBytes, err := jsoniter.Marshal(&PosMasterArray)

	if err != nil {
		RecordErrorLog(err)
		MarkTransmitAlarm(aidx, "GoodsR-1", "E", v)
		return 0
	}

	buff := bytes.NewBuffer(jsonBytes)
	resp, err := http.Post(DestinUrl+EpProduts, Contype, buff)
	res = resp.StatusCode
	if err != nil {
		RecordErrorLog(err)
		MarkTransmitAlarm(aidx, "GoodsR-2", "E", v)
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
		MarkTransmitAlarm(aidx, "GoodsR-3", "E", v)
		return 0
	}
	return res
}

// 마트회원 데이터 Count
func TrnMemberRecordCount(aidx string, v string, db *sql.DB) (cnt int) {
	var buf bytes.Buffer
	buf.WriteString(`SELECT COUNT(mem_code) FROM member WHERE mem_status = 1`)

	err := db.QueryRow(buf.String()).Scan(&cnt)
	defer runtime.GC()

	if err != nil {
		RecordErrorLog(err)
		MarkTransmitAlarm(aidx, "MemberRCnt", "E", v)
		return 0
	}
	return
}

// 마트회원 데이터 Transmit
func TrnMemberRecord(aidx string, v string, db *sql.DB, istart int, ilimit int) int {
	var nrows, mem_code, mem_point, mem_pointsum int
	var mem_name, mem_bcode, mem_tel, mem_rdate, mem_ord_date string
	var res int = 0

	var buf bytes.Buffer
	buf.WriteString(`SELECT M.*
		FROM (
			SELECT
				ROW_NUMBER() OVER (ORDER BY mem_code ASC) AS nrows,
				mem_code,
				mem_name,
				mem_bcode,
				mem_cel,
				mem_point,
				mem_pointsum,
				CONVERT(CHAR(19), mem_rdate, 120) AS mem_rdate,
				ISNULL(CONVERT(CHAR(19), mem_hdate, 120), '0000-00-00 00:00:00') AS mem_hdate
			FROM member
			WHERE mem_status = 1
		) AS M
		WHERE M.nrows BETWEEN ? AND ?`)

	rows, err := db.Query(buf.String(), istart, ilimit)
	defer runtime.GC()

	if err != nil {
		RecordErrorLog(err)
		MarkTransmitAlarm(aidx, "MemberR-query", "E", v)
		return 0
	}
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()

	var MartMemberArray []*MartMemberData
	for rows.Next() {
		res = 0
		err := rows.Scan(&nrows, &mem_code, &mem_name, &mem_bcode, &mem_tel, &mem_point, &mem_pointsum, &mem_rdate, &mem_ord_date)
		if err != nil {
			RecordErrorLog(err)
			MarkTransmitAlarm(aidx, "MemberR-rows", "E", v)
			return 0
		}

		MartMemberArray = append(MartMemberArray, &MartMemberData{
			App_idx:    aidx,
			Mem_code:   mem_code,
			Mem_name:   mem_name,
			Mem_bcode:  mem_bcode,
			Mem_tel:    mem_tel,
			Mem_point:  mem_point,
			Mem_tpoint: mem_pointsum,
			Mem_rdate:  mem_rdate,
			Mem_odate:  mem_ord_date,
		})
	}

	jsonBytes, err := jsoniter.Marshal(&MartMemberArray)

	if err != nil {
		RecordErrorLog(err)
		MarkTransmitAlarm(aidx, "MemberR-1", "E", v)
		return 0
	}

	buff := bytes.NewBuffer(jsonBytes)
	resp, err := http.Post(DestinUrl+EpMembers, Contype, buff)
	res = resp.StatusCode
	if err != nil {
		RecordErrorLog(err)
		MarkTransmitAlarm(aidx, "MemberR-2", "E", v)
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
		MarkTransmitAlarm(aidx, "MemberR-3", "E", v)
		return 0
	}
	return res
}

// 행사정보 Master데이터 Transmit
func TrnBargainMaster(aidx string, v string, db *sql.DB) []int {
	var bg_code int
	var bg_title, bg_rsdate, bg_edate string
	var bgarr []int

	var buf bytes.Buffer
	buf.WriteString(`SELECT
			bg_code,
			bg_title,
			CONVERT(CHAR(10), CONVERT(DATETIME, CONVERT(CHAR(8), bg_start)), 120) + ' 00:00:00' AS sdate,
			CONVERT(CHAR(10), CONVERT(DATETIME, CONVERT(CHAR(8), bg_end)), 120) + ' 23:59:59' AS edate
		FROM bargain`)

	rows, err := db.Query(buf.String())
	defer runtime.GC()

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
	buf.WriteString(`SELECT COUNT(S.bgg_goods)
		FROM bargaingoods AS S
			INNER JOIN goods AS G ON G.goods_code = S.bgg_goods
		WHERE bgg_bargain = ? AND bgg_status = 1`)

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
	var nrows, bgg_goods, bgg_qty, bgg_limit, bgg_hour, bgg_dccount, bgg_dcamount, bgg_price int
	var bgg_bcode string
	var res int = 0

	var buf bytes.Buffer
	buf.WriteString(`SELECT MAIN.*
		FROM (
			SELECT
				ROW_NUMBER() OVER (ORDER BY G.goods_bcode ASC) AS nrows,
				S.bgg_goods,
				G.goods_bcode,
				S.bgg_qty,
				S.bgg_limit,
				S.bgg_hour,
				S.bgg_dccount,
				S.bgg_dcamount,
				S.bgg_sprice
			FROM bargaingoods AS S
				INNER JOIN goods AS G ON G.goods_code = S.bgg_goods
			WHERE S.bgg_bargain = ? AND S.bgg_status = 1
		) AS MAIN
		WHERE MAIN.nrows BETWEEN ? AND ?`)

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
		err := rows.Scan(&nrows, &bgg_goods, &bgg_bcode, &bgg_qty, &bgg_limit, &bgg_hour, &bgg_dccount, &bgg_dcamount, &bgg_price)
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

// 회원구매 이력 데이터 Count(70분 이내 신규발생 데이터)
func TrnUpdateOrderedRecordCount(aidx string, v string, db *sql.DB) (cnt int) {
	var buf bytes.Buffer
	buf.WriteString(`SELECT SUM(MAIN.rcnt) AS rcnt
		FROM (
			SELECT COUNT(R.ord_code) AS rcnt
			FROM Ord AS R
				INNER JOIN member AS M ON M.mem_code = R.ord_memcode
			WHERE M.mem_status = 1
				AND R.ord_memcode <> 0
				AND DATEDIFF(MINUTE, R.ord_date, CURRENT_TIMESTAMP) BETWEEN 0 AND 70
			UNION ALL
			SELECT COUNT(R.ord_code) AS rcnt
			FROM Ordd AS P
				INNER JOIN Ord AS R ON P.ordd_code = R.ord_code
				INNER JOIN member AS M ON M.mem_code = R.ord_memcode
			WHERE M.mem_status = 1
				AND R.ord_memcode <> 0
				AND DATEDIFF(MINUTE, R.ord_date, CURRENT_TIMESTAMP) BETWEEN 0 AND 70
		) AS MAIN`)

	err := db.QueryRow(buf.String()).Scan(&cnt)
	defer runtime.GC()

	if err != nil {
		RecordErrorLog(err)
		MarkTransmitAlarm(aidx, "OrderMCnt", "E", v)
		return 0
	}
	return
}

// 회원구매 이력 데이터 Transmit(70분 이내 신규발생 데이터)
func TrnUpdateOrderedRecord(aidx string, v string, db *sql.DB, istart int, ilimit int) int {
	var nrows, dhead, muid, paycash, paycard, paypoint, gprice int
	var mbcode, ordcd, gcode, bcode, ord_date string
	var res int = 0

	var buf bytes.Buffer
	buf.WriteString(`SELECT MAIN.*
		FROM (
			SELECT
				ROW_NUMBER() OVER (ORDER BY SUB.ord_date DESC) AS nrows, SUB.*
			FROM(
				SELECT
					1 AS header,
					M.mem_code,
					RTRIM(M.mem_bcode) AS mem_bcode,
					R.ord_code,
					R.ord_paycash,
					R.ord_paycredit,
					R.ord_paypoint,
					0 AS ordd_goods,
					'' AS goods_bcode,
					0 AS ordd_price,
					R.ord_date
				FROM Ord AS R
					INNER JOIN member AS M ON M.mem_code = R.ord_memcode
				WHERE M.mem_status = 1
					AND R.ord_memcode <> 0
					AND DATEDIFF(MINUTE, R.ord_date, CURRENT_TIMESTAMP) BETWEEN 0 AND 70
				UNION ALL
				SELECT
					10 AS header,
					M.mem_code,
					RTRIM(M.mem_bcode) AS mem_bcode,
					R.ord_code,
					0 AS ord_paycash,
					0 AS ord_paycredit,
					0 AS ord_paypoint,
					P.ordd_goods,
					(SELECT TOP 1 goods_bcode FROM goods WHERE goods_code = P.ordd_goods) AS goods_bcode,
					P.ordd_price,
					R.ord_date
				FROM Ordd AS P
					INNER JOIN Ord AS R ON P.ordd_code = R.ord_code
					INNER JOIN member AS M ON M.mem_code = R.ord_memcode
				WHERE M.mem_status = 1
					AND R.ord_memcode <> 0
					AND DATEDIFF(MINUTE, R.ord_date, CURRENT_TIMESTAMP) BETWEEN 0 AND 70
			) AS SUB
		) AS MAIN
		WHERE MAIN.nrows BETWEEN ? AND ?`)

	rows, err := db.Query(buf.String(), istart, ilimit)
	defer runtime.GC()

	if err != nil {
		RecordErrorLog(err)
		MarkTransmitAlarm(aidx, "OrderM-query", "E", v)
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
		err := rows.Scan(&nrows, &dhead, &muid, &mbcode, &ordcd, &paycash, &paycard, &paypoint, &gcode, &bcode, &gprice, &ord_date)
		if err != nil {
			RecordErrorLog(err)
			MarkTransmitAlarm(aidx, "OrderM-rows", "E", v)
			return 0
		}

		OrderedDataArray = append(OrderedDataArray, &OrderedData{
			App_idx:         aidx,
			Ord_head:        dhead,
			Mem_code:        muid,
			Mem_bcode:       mbcode,
			Ord_code:        ordcd,
			Pay_cash_price:  paycash,
			Pay_card_price:  paycard,
			Pay_point_price: paypoint,
			Gds_code:        gcode,
			Gds_bcode:       bcode,
			Pay_goods_price: gprice,
			Ord_date:        ord_date,
		})
	}

	jsonBytes, err := jsoniter.Marshal(&OrderedDataArray)

	if err != nil {
		RecordErrorLog(err)
		MarkTransmitAlarm(aidx, "OrderM-1", "E", v)
		return 0
	}

	buff := bytes.NewBuffer(jsonBytes)
	resp, err := http.Post(DestinUrl+EpUpdateOrdered, Contype, buff)
	res = resp.StatusCode
	if err != nil {
		RecordErrorLog(err)
		MarkTransmitAlarm(aidx, "OrderM-2", "E", v)
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
		MarkTransmitAlarm(aidx, "OrderM-3", "E", v)
		return 0
	}
	return res
}

// 판매제품 갱신데이터 Count
func TrnModifiedRecordCount(aidx string, v string, db *sql.DB) (cnt int) {
	var buf bytes.Buffer
	buf.WriteString(`SELECT COUNT(gl_goods) FROM goodslog
		WHERE DATEDIFF(MINUTE, gl_date, CURRENT_TIMESTAMP) BETWEEN 0 AND 70`)

	err := db.QueryRow(buf.String()).Scan(&cnt)
	defer runtime.GC()

	if err != nil {
		RecordErrorLog(err)
		MarkTransmitAlarm(aidx, "GoodsMCnt", "E", v)
		return 0
	}
	return
}

// 판매제품 갱신데이터 Transmit
func TrnModifiedRecord(aidx string, v string, db *sql.DB, istart int, ilimit int) int {

	var nrows, gds_sprice int
	var gds_bprice float32
	var cate_code, gds_code, gds_bcode, gds_name string
	var res int = 0

	var buf bytes.Buffer
	buf.WriteString(`SELECT M.*
		FROM (
			SELECT
				ROW_NUMBER() OVER (ORDER BY l.gl_date DESC) AS nrows,
				C.gb_code,
				G.goods_code,
				RTRIM(G.goods_bcode) AS gds_bcode,
				(LTRIM(RTRIM(G.goods_name)) + ' ' + RTRIM(G.goods_sspec)) AS gds_name,
				G.goods_bprice,
				G.goods_sprice
			FROM goods AS G
				INNER JOIN gbook AS C ON C.gb_top = G.goods_bktop AND C.gb_mid = G.goods_bkmid AND C.gb_bot = G.goods_bkbot
				INNER JOIN goodslog AS L ON L.gl_goods = G.goods_code
			WHERE G.goods_status = 1
				AND DATEDIFF(MINUTE, L.gl_date, CURRENT_TIMESTAMP) BETWEEN 0 AND 70
		) AS M
		WHERE M.nrows BETWEEN ? AND ?`)

	rows, err := db.Query(buf.String(), istart, ilimit)
	defer runtime.GC()

	if err != nil {
		RecordErrorLog(err)
		MarkTransmitAlarm(aidx, "GoodsM-query", "E", v)
		return 0
	}
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()

	var PosMasterArray []*PosMasterData
	for rows.Next() {
		res = 0
		err := rows.Scan(&nrows, &cate_code, &gds_code, &gds_bcode, &gds_name, &gds_bprice, &gds_sprice)
		if err != nil {
			RecordErrorLog(err)
			MarkTransmitAlarm(aidx, "GoodsM-rows", "E", v)
			return 0
		}

		PosMasterArray = append(PosMasterArray, &PosMasterData{
			App_idx:    aidx,
			Cate_code:  cate_code,
			Gds_code:   gds_code,
			Gds_bcode:  gds_bcode,
			Gds_name:   gds_name,
			Gds_bprice: gds_bprice,
			Gds_sprice: gds_sprice,
		})
	}

	jsonBytes, err := jsoniter.Marshal(&PosMasterArray)

	if err != nil {
		RecordErrorLog(err)
		MarkTransmitAlarm(aidx, "GoodsM-1", "E", v)
		return 0
	}

	buff := bytes.NewBuffer(jsonBytes)
	resp, err := http.Post(DestinUrl+EpModifyPrds, Contype, buff)
	res = resp.StatusCode
	if err != nil {
		RecordErrorLog(err)
		MarkTransmitAlarm(aidx, "GoodsM-2", "E", v)
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
		MarkTransmitAlarm(aidx, "GoodsM-3", "E", v)
		return res
	}
	return res
}

// 회원정보 갱신데이터 Count
func TrnMemberModifiedCount(aidx string, v string, db *sql.DB) (cnt int) {
	var buf bytes.Buffer
	buf.WriteString(`SELECT COUNT(mem_code) FROM member
		WHERE mem_status = 1 AND DATEDIFF(MINUTE, mem_edate, CURRENT_TIMESTAMP) BETWEEN 0 AND 70`)

	err := db.QueryRow(buf.String()).Scan(&cnt)
	defer runtime.GC()

	if err != nil {
		RecordErrorLog(err)
		MarkTransmitAlarm(aidx, "MemberMCnt", "E", v)
		return 0
	}
	return
}

// 회원정보 갱신데이터 Transmit
func TrnMemberModifiedRecords(aidx string, v string, db *sql.DB, istart int, ilimit int) int {
	var nrows, mem_code, mem_point, mem_pointsum int
	var mem_name, mem_bcode, mem_tel, mem_rdate, mem_ord_date string
	var res int = 0

	var buf bytes.Buffer
	buf.WriteString(`SELECT M.*
		FROM (
			SELECT
				ROW_NUMBER() OVER (ORDER BY mem_code ASC) AS nrows,
				mem_code,
				mem_name,
				mem_bcode,
				mem_cel,
				mem_point,
				mem_pointsum,
				CONVERT(CHAR(19), mem_rdate, 120) AS mem_rdate,
				ISNULL(CONVERT(CHAR(19), mem_hdate, 120), '0000-00-00 00:00:00') AS mem_hdate
			FROM member
			WHERE mem_status = 1
				AND DATEDIFF(MINUTE, mem_edate, CURRENT_TIMESTAMP) BETWEEN 0 AND 70
		) AS M
		WHERE M.nrows BETWEEN ? AND ?`)

	rows, err := db.Query(buf.String(), istart, ilimit)
	defer runtime.GC()

	if err != nil {
		RecordErrorLog(err)
		MarkTransmitAlarm(aidx, "MemberM-query", "E", v)
		return res
	}
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()

	var MartMemberArray []*MartMemberData
	for rows.Next() {
		res = 0
		err := rows.Scan(&nrows, &mem_code, &mem_name, &mem_bcode, &mem_tel, &mem_point, &mem_pointsum, &mem_rdate, &mem_ord_date)
		if err != nil {
			RecordErrorLog(err)
			MarkTransmitAlarm(aidx, "MemberM-rows", "E", v)
			return res
		}

		MartMemberArray = append(MartMemberArray, &MartMemberData{
			App_idx:    aidx,
			Mem_code:   mem_code,
			Mem_name:   mem_name,
			Mem_bcode:  mem_bcode,
			Mem_tel:    mem_tel,
			Mem_point:  mem_point,
			Mem_tpoint: mem_pointsum,
			Mem_rdate:  mem_rdate,
			Mem_odate:  mem_ord_date,
		})
	}

	jsonBytes, err := jsoniter.Marshal(&MartMemberArray)

	if err != nil {
		RecordErrorLog(err)
		MarkTransmitAlarm(aidx, "MemberM-1", "E", v)
		return res
	}

	buff := bytes.NewBuffer(jsonBytes)
	resp, err := http.Post(DestinUrl+EpModifyMems, Contype, buff)
	res = resp.StatusCode
	if err != nil {
		RecordErrorLog(err)
		MarkTransmitAlarm(aidx, "MemberM-2", "E", v)
		return res
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
		MarkTransmitAlarm(aidx, "MemberM-3", "E", v)
		return res
	}
	return res
}

// 행사정보 갱신데이터 Count
func TrnBargainModifiedCount(aidx string, v string, db *sql.DB) (cnt int) {
	var buf bytes.Buffer
	buf.WriteString(`SELECT COUNT(G.bgg_bargain)
		FROM bargaingoods AS G
			INNER JOIN bargainlog AS L ON L.bl_goods = G.bgg_goods
			LEFT JOIN goods AS P ON P.goods_code = G.bgg_goods
		WHERE G.bgg_status = 1
			AND DATEDIFF(MINUTE, L.bl_date, CURRENT_TIMESTAMP) BETWEEN 0 AND 70`)

	err := db.QueryRow(buf.String()).Scan(&cnt)
	defer runtime.GC()

	if err != nil {
		RecordErrorLog(err)
		MarkTransmitAlarm(aidx, "BargainMCnt", "E", v)
		return 0
	}
	return
}

// 행사정보 갱신데이터 Transmit
func TrnBargainModifiedRecords(aidx string, v string, db *sql.DB, istart int, ilimit int) int {
	var nrows, bgg_idx, bgg_qty, bgg_limit, bgg_hour, bgg_dccount, bgg_dcamount, bgg_price int
	var bgg_bcode string
	var res int = 0

	var buf bytes.Buffer
	buf.WriteString(`SELECT M.*
		FROM (
			SELECT
				ROW_NUMBER() OVER (ORDER BY L.bl_date DESC) AS nrows,
				G.bgg_bargain,
				RTRIM(P.goods_bcode) AS gds_bcode,
				G.bgg_qty,
				G.bgg_limit,
				G.bgg_hour,
				G.bgg_dccount,
				G.bgg_dcamount,
				G.bgg_sprice
			FROM bargaingoods AS G
				INNER JOIN bargainlog AS L ON L.bl_goods = G.bgg_goods
				LEFT JOIN goods AS P ON P.goods_code = G.bgg_goods
			WHERE G.bgg_status = 1
				AND DATEDIFF(MINUTE, L.bl_date, CURRENT_TIMESTAMP) BETWEEN 0 AND 70
		) AS M
		WHERE M.nrows BETWEEN ? AND ?`)

	rows, err := db.Query(buf.String(), istart, ilimit)
	defer runtime.GC()

	if err != nil {
		RecordErrorLog(err)
		MarkTransmitAlarm(aidx, "BargainM-query", "E", v)
		return res
	}
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()

	var BargainProductsArray []*BargainProductsData
	for rows.Next() {
		res = 0
		err := rows.Scan(&nrows, &bgg_idx, &bgg_bcode, &bgg_qty, &bgg_limit, &bgg_hour, &bgg_dccount, &bgg_dcamount, &bgg_price)
		if err != nil {
			RecordErrorLog(err)
			MarkTransmitAlarm(aidx, "BargainM-rows", "E", v)
			return res
		}

		BargainProductsArray = append(BargainProductsArray, &BargainProductsData{
			App_idx:     aidx,
			Bgg_code:    bgg_idx,
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
		MarkTransmitAlarm(aidx, "BargainM-1", "E", v)
		return res
	}

	buff := bytes.NewBuffer(jsonBytes)
	resp, err := http.Post(DestinUrl+EpModifyBgs, Contype, buff)
	res = resp.StatusCode
	if err != nil {
		RecordErrorLog(err)
		MarkTransmitAlarm(aidx, "BargainM-2", "E", v)
		return res
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
		MarkTransmitAlarm(aidx, "BargainM-3", "E", v)
	}
	return res
}
