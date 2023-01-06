package main

import (
	"fmt"
	"net/smtp"
	"os"
)

type Log struct {
	pathToFolder, fileName, hash string
}

func main() {

	var pathToFolder string = os.Args[1]
	var fileName string = os.Args[2]
	var email string = os.Args[3]
	var pathToLoggerHash string = os.Args[4]

	fmt.Println(email)
	fmt.Println(pathToFolder)
	fmt.Println(fileName)
	fmt.Println(pathToLoggerHash)

	if loggerHashChecker(pathToLoggerHash, pathToFolder, fileName) {
		emailSend(email, fileName)
	}

}

func emailSend(email string, fileName string) {
	from := "melentev.av@gmail.com"
	password := "rssjsfqhfeglosql"

	toEmailAddress := email
	to := []string{toEmailAddress}

	host := "smtp.gmail.com"
	port := "587"
	address := host + ":" + port

	subject := "Что-то пошло не так: Лог " + fileName + " поменял свое содержимое \n"
	body := "Вы получили это письмо так как что-то случилось с сервисом" + fileName
	message := []byte(subject + body)

	auth := smtp.PlainAuth("", from, password, host)

	err := smtp.SendMail(address, auth, from, to, message)
	if err != nil {
		panic(err)
	}
}

func loggerHashChecker(pathToLoggerHash string, pathToFolder string, filename string) bool {
	return true
}

func writeToJson() {

}
