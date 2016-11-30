package main

import (
	"net/smtp"
	"log"
)

func main() {
	// 메일서버 로그인 정보 설정
	auth := smtp.PlainAuth("", "wowadd@naver.com", "wlsk@0314", "smtp.naver.com")

	from := "wowadd@naver.com"
	to := []string{"wowplus2@gmail.com", "wowplus@revolution.co.kr"}
	//to := []string{"wowplus@revolution.co.kr", "wowadd@naver.com"}	// 복수 수신자 가능

	// 메세지 작성
	headerSubject := "Subject: Golang 이메일 발송 테스트\r\n"
	headerBlank := "\r\n"
	body := "Golang net/smtp 이메일 발송 테스트입니다.\r\n"
	msg := []byte(headerSubject + headerBlank + body)

	// 이메일 발송
	err := smtp.SendMail("smtp.naver.com:465", auth, from, to, msg)
	if err != nil {
		//panic(err)
		log.Fatal("ERROR => ", err)
	}
}
