package app

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
	EpOrdered       = "posys/regiReceipt"
	EpUpdateOrdered = "posys/updateReceipt"
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
	App_idx     string  `json:"aidx"`
	Cate_code   string  `json:"cate"`
	Gds_code    string  `json:"pcode"`
	Gds_bcode   string  `json:"bcode"`
	Gds_name    string  `json:"name"`
	Gds_bprice  float32 `json:"bprice"`
	Gds_sprice  int     `json:"sprice"`
	Gds_issales string  `json:"issales"`
}

type MartMemberData struct {
	App_idx    string `json:"aidx"`
	Mem_code   int    `json:"muid"`
	Mem_name   string `json:"name"`
	Mem_bcode  string `json:"bcode"`
	Mem_tel    string `json:"telno"`
	Mem_point  int    `json:"mpt"`
	Mem_tpoint string `json:"spt"`
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
	Bgg_amount  string `json:"amount"`
	Bgg_price   int    `json:"price"`
}

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

// 회원구매 이력 데이터 Count(2일 이내 데이터)
func TrnOrderedRecordCount(aidx string, v string, db *sql.DB) (tcnt []int, darr []string) {

	cnt := 0
	date := ""
	for i := 0; i < 2; i++ {
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

// 회원구매 이력 데이터 Transmit
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
	WHERE N.mbpan_key_date = ? AND P.panme_dam_name != '' AND P.panme_rec_type < 6 LIMIT ?, ?`)

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

// 판매제품 카테고리 Record Count Transmit
func TrnCategoryCount(aidx string, v string, db *sql.DB) (cnt int) {

	var buf bytes.Buffer
	buf.WriteString(`SELECT COUNT(*)
		FROM sdae AS D
			LEFT JOIN sjung AS J ON J.jung_dae = D.dae_code
			LEFT JOIN sso AS S ON S.so_dae = D.dae_code AND S.so_jung = J.jung_code
		WHERE D.dae_flag = 'F' AND J.jung_flag = 'F' AND S.so_flag = 'F'`)

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

	var gb_top, gb_mid, gb_bot, gb_code int
	var gb_name string
	var res int = 0

	var buf bytes.Buffer
	buf.WriteString(`SELECT
			D.dae_code AS gb_top,
			J.jung_code AS gb_mid,
			S.so_code AS gb_bot,
			CONCAT(D.dae_name, '>', J.jung_name, '>', S.so_name) AS gb_name,
			CONCAT(LPAD(D.dae_code, 3, '0'), LPAD(J.jung_code, 3, '0'), LPAD(S.so_code, 3, '0')) AS gb_code
		FROM sdae AS D
			LEFT JOIN sjung AS J ON J.jung_dae = D.dae_code
			LEFT JOIN sso AS S ON S.so_dae = D.dae_code AND S.so_jung = J.jung_code
		WHERE D.dae_flag = 'F' AND J.jung_flag = 'F' AND S.so_flag = 'F' ORDER BY D.dae_code ASC LIMIT ?, ?`)

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
		err := rows.Scan(&gb_top, &gb_mid, &gb_bot, &gb_name, &gb_code)
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
	buf.WriteString(`SELECT COUNT(*) FROM spum WHERE pum_ipgo_yn IN ('0', '1') AND pum_danga > 0`)

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

	var gds_sprice int
	var gds_bprice float32
	var cate_code, gds_code, gds_bcode, gds_name, gds_issales string
	var res int = 0

	var buf bytes.Buffer
	buf.WriteString(`SELECT
			CONCAT(LPAD(P.pum_dae, 3, '0'), LPAD(P.pum_jung, 3, '0'), LPAD(P.pum_so, 3, '0')) AS gb_code,
			(@rnum := @rnum + 1) AS goods_code,
			P.pum_code AS gds_bcode,
			CONCAT(P.pum_name, ' ', IF(LENGTH(P.pum_size) > 0, P.pum_size, P.pum_unit)) AS gds_name,
			P.pum_wonga AS goods_bprice,
			P.pum_danga AS goods_sprice,
			'Y' AS issale
		FROM spum AS P
			LEFT JOIN (SELECT @rnum := 0) AS T ON 1 = 1
		WHERE P.pum_ipgo_yn IN ('0', '1') AND P.pum_danga > 0 LIMIT ?, ?`)

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
		err := rows.Scan(&cate_code, &gds_code, &gds_bcode, &gds_name, &gds_bprice, &gds_sprice, &gds_issales)
		if err != nil {
			RecordErrorLog(err)
			MarkTransmitAlarm(aidx, "GoodsR-rows", "E", v)
			return 0
		}

		PosMasterArray = append(PosMasterArray, &PosMasterData{
			App_idx:     aidx,
			Cate_code:   cate_code,
			Gds_code:    gds_code,
			Gds_bcode:   gds_bcode,
			Gds_name:    gds_name,
			Gds_bprice:  gds_bprice,
			Gds_sprice:  gds_sprice,
			Gds_issales: gds_issales,
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
	buf.WriteString(`SELECT COUNT(*)
		FROM (
			SELECT M.member_code
			FROM smember AS M
				INNER JOIN smempot AS P ON P.mempot_jum = M.member_jum AND P.mempot_member_code = M.member_code
			WHERE LENGTH(M.member_hd_phone) > 0
				AND M.member_jum = '001' AND M.member_last_date != '' GROUP BY M.member_code
		) AS M`)

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

	var mem_code, mem_point int
	var mem_name, mem_bcode, mem_tel, mem_rdate, mem_ord_date, mem_pointsum string
	var res int = 0

	var buf bytes.Buffer
	buf.WriteString(`SELECT
			(@rnum := @rnum + 1) AS mem_code,
			M.member_name AS mem_name,
			M.member_code AS mem_bcode,
			M.member_hd_phone AS mem_cel,
			0 AS mem_point,
			SUM(P.mempot_amt) AS mem_pointsum,
			CONCAT(DATE(M.member_input_date), ' ', M.member_pwd) AS mem_rdate,
			CONCAT(DATE(M.member_last_date), ' ', M.member_last_time) AS mem_hdate
		FROM smember AS M
			LEFT JOIN (SELECT @rnum := 0) AS T ON 1 = 1
			INNER JOIN smempot AS P ON P.mempot_jum = M.member_jum AND P.mempot_member_code = M.member_code
		WHERE LENGTH(M.member_hd_phone) > 0 AND M.member_jum = '001' AND M.member_last_date != ''
		GROUP BY M.member_code LIMIT ?, ?`)

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
		err := rows.Scan(&mem_code, &mem_name, &mem_bcode, &mem_tel, &mem_point, &mem_pointsum, &mem_rdate, &mem_ord_date)
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
			CONVERT(REPLACE(salecode_code, '-', ''), UNSIGNED INT) AS bg_code,
			salecode_name AS bg_title,
			CONVERT(salecode_ss_date, DATETIME) AS bg_sdate,
			CONCAT(CONVERT(salecode_se_date, DATE), ' 23:59:59') AS bg_edate
		FROM ssalecode`)

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
	buf.WriteString(`SELECT COUNT(S.salepum_sale_code)
			FROM ssalepum AS S
				INNER JOIN ssalecode AS M ON M.salecode_code = S.salepum_sale_code
				INNER JOIN spum AS P ON P.pum_code = S.salepum_pum_code
			WHERE P.pum_ipgo_yn IN ('0', '1')
				AND CONVERT(REPLACE(M.salecode_code, '-', ''), UNSIGNED INT) = ?`)

	err := db.QueryRow(buf.String(), &bgidx).Scan(&cnt)

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

// 회원구매 이력 데이터 Count(70분 이내 신규발생 데이터)
func TrnUpdateOrderedRecordCount(aidx string, v string, db *sql.DB) (cnt int) {

	var buf bytes.Buffer
	buf.WriteString(`SELECT COUNT(*)
		FROM zmbpan AS N
			INNER JOIN smember AS M ON M.member_code = N.mbpan_member_code
			INNER JOIN zpanme AS P ON P.panme_jum = N.mbpan_jum AND P.panme_key_date = N.mbpan_key_date AND P.panme_pos_no = N.mbpan_pos_no AND P.panme_junpo_no = N.mbpan_junpo_no
		WHERE P.panme_key_date = CONVERT(DATE_FORMAT(DATE(NOW()), '%y-%m-%d'), CHAR(8))
			AND P.panme_dam_name != ''
			AND P.panme_rec_type < 6
			AND TIMESTAMPDIFF(MINUTE, CONCAT(DATE(P.panme_date), ' ', P.panme_time), NOW()) BETWEEN 0 AND 70`)

	err := db.QueryRow(buf.String()).Scan(&cnt)
	defer runtime.GC()

	if err != nil {
		RecordErrorLog(err)
		MarkTransmitAlarm(aidx, "OrderMCnt", "E", v)
		return 0
	}
	return
}

// 회원구매 이력 데이터 Transmit(금일 70분 이내 신규발생 데이터)
func TrnUpdateOrderedRecord(aidx string, v string, db *sql.DB, istart int, ilimit int) int {

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
			IF(ASCII(SUBSTRING(P.panme_time, 1,2))>47 && ASCII(SUBSTRING(P.panme_time, 1,2))<58, CONCAT(DATE(P.panme_date), ' ', P.panme_time), '') AS ord_date
		FROM zmbpan AS N
			INNER JOIN smember AS M ON M.member_code = N.mbpan_member_code
			INNER JOIN zpanme AS P ON P.panme_jum = N.mbpan_jum AND P.panme_key_date = N.mbpan_key_date AND P.panme_pos_no = N.mbpan_pos_no AND P.panme_junpo_no = N.mbpan_junpo_no
		WHERE P.panme_key_date = CONVERT(DATE_FORMAT(DATE(NOW()), '%y-%m-%d'), CHAR(8))
			AND P.panme_dam_name != ''
			AND P.panme_rec_type < 6
			AND TIMESTAMPDIFF(MINUTE, CONCAT(DATE(P.panme_date), ' ', P.panme_time), NOW()) BETWEEN 0 AND 70
		LIMIT ?, ?`)

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
		err := rows.Scan(&dHead, &mbCode, &ordCd, &recType, &seqNo, &pmField, &pmCode, &pmPrice, &pmQty, &pmTotal, &ordDate)
		if err != nil {
			RecordErrorLog(err)
			MarkTransmitAlarm(aidx, "OrderM-rows", "E", v)
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
	buf.WriteString(`SELECT COUNT(*) FROM spum
			WHERE pum_ipgo_yn IN ('0', '1')
				AND pum_danga > 0
				AND TIMESTAMPDIFF(MINUTE, CONCAT(DATE(pum_update_date), ' ', pum_update_time), NOW()) BETWEEN 0 AND 70`)

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

	var gds_sprice int
	var gds_bprice float32
	var cate_code, gds_code, gds_bcode, gds_name string
	var res int = 0

	var buf bytes.Buffer
	buf.WriteString(`SELECT
			CONCAT(LPAD(P.pum_dae, 3, '0'), LPAD(P.pum_jung, 3, '0'), LPAD(P.pum_so, 3, '0')) AS gb_code,
			(@rnum := @rnum + 1) AS goods_code,
			P.pum_code AS gds_bcode,
			CONCAT(P.pum_name, ' ', IF(LENGTH(P.pum_size) > 0, P.pum_size, P.pum_unit)) AS gds_name,
			P.pum_wonga AS goods_bprice,
			P.pum_danga AS goods_sprice,
			'Y' AS issale
		FROM spum AS P
			LEFT JOIN (SELECT @rnum := 0) AS T ON 1 = 1
		WHERE P.pum_ipgo_yn IN ('0', '1') AND P.pum_danga > 0
			AND TIMESTAMPDIFF(MINUTE, CONCAT(DATE(P.pum_update_date), ' ', P.pum_update_time), NOW()) BETWEEN 0 AND 70
		LIMIT ?, ?`)

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
		err := rows.Scan(&cate_code, &gds_code, &gds_bcode, &gds_name, &gds_bprice, &gds_sprice)
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
	buf.WriteString(`SELECT COUNT(*) FROM (
			SELECT M.member_code
			FROM smember AS M
				INNER JOIN smempot AS P ON P.mempot_jum = M.member_jum AND P.mempot_member_code = M.member_code
			WHERE LENGTH(M.member_hd_phone) > 0
				AND M.member_jum = '001'
				AND M.member_last_date != ''
				AND TIMESTAMPDIFF(MINUTE, CONCAT(DATE(M.member_last_date), ' ', M.member_last_time), NOW()) BETWEEN 0 AND 70
        	GROUP BY M.member_code
        ) AS M`)

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

	var mem_code, mem_point int
	var mem_name, mem_bcode, mem_tel, mem_rdate, mem_ord_date, mem_pointsum string
	var res int = 0

	var buf bytes.Buffer
	buf.WriteString(`SELECT
			(@rnum := @rnum + 1) AS mem_code,
			M.member_name AS mem_name,
			M.member_code AS mem_bcode,
			M.member_hd_phone AS mem_cel,
			0 AS mem_point,
			SUM(P.mempot_amt) AS mem_pointsum,
			CONCAT(DATE(M.member_input_date), ' ', M.member_pwd) AS mem_rdate,
			CONCAT(DATE(M.member_last_date), ' ', M.member_last_time) AS mem_hdate
		FROM smember AS M
			LEFT JOIN (SELECT @rnum := 0) AS T ON 1 = 1
			INNER JOIN smempot AS P ON P.mempot_jum = M.member_jum AND P.mempot_member_code = M.member_code
		WHERE LENGTH(M.member_hd_phone) > 0
			AND M.member_jum = '001'
			AND M.member_last_date != ''
			AND TIMESTAMPDIFF(MINUTE, CONCAT(DATE(M.member_last_date), ' ', M.member_last_time), NOW()) BETWEEN 0 AND 70
		GROUP BY M.member_code LIMIT ?, ?`)

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
		err := rows.Scan(&mem_code, &mem_name, &mem_bcode, &mem_tel, &mem_point, &mem_pointsum, &mem_rdate, &mem_ord_date)
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

// 행사상품 갱신데이터 Count
func TrnBargainModifiedCount(aidx string, v string, db *sql.DB) (cnt int) {

	var buf bytes.Buffer
	buf.WriteString(`SELECT COUNT(S.salepum_sale_code)
			FROM ssalepum AS S
				INNER JOIN ssalecode AS M ON M.salecode_code = S.salepum_sale_code
				INNER JOIN spum AS P ON P.pum_code = S.salepum_pum_code
			WHERE P.pum_ipgo_yn IN ('0', '1')
				AND DATE(NOW()) BETWEEN DATE(M.salecode_ss_date) AND DATE(M.salecode_se_date)`)

	err := db.QueryRow(buf.String()).Scan(&cnt)
	defer runtime.GC()

	if err != nil {
		RecordErrorLog(err)
		MarkTransmitAlarm(aidx, "BargainMCnt", "E", v)
		return 0
	}
	return
}

// 행사상품 갱신데이터 Transmit
func TrnBargainModifiedRecords(aidx string, v string, db *sql.DB, istart int, ilimit int) int {

	var bgg_idx, bgg_qty, bgg_limit, bgg_hour, bgg_dccount, bgg_price int
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
			AND DATE(NOW()) BETWEEN DATE(M.salecode_ss_date) AND DATE(M.salecode_se_date)
		ORDER BY S.salepum_pum_code LIMIT ?, ?`)

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
		err := rows.Scan(&bgg_idx, &bgg_bcode, &bgg_qty, &bgg_limit, &bgg_hour, &bgg_dccount, &bgg_dcamount, &bgg_price)
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
