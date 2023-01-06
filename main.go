package main

import (
	"crypto/md5"
	"encoding/json"
	"io"
	"log"
	"net/smtp"
	"os"
)

type Log struct {
	pathToFolder, fileName, hash string
}

func main() {

	var pathToFolder string = os.Args[1]
	var fileName string = os.Args[2]
	//var email string = os.Args[3]
	var pathToLoggerHash string = os.Args[4]

	if loggerHashChecker(pathToLoggerHash, pathToFolder, fileName) {
		//emailSend(email, fileName)
	}

}

func emailSend(email string, fileName string, pwdEmail string, fromEmail string) {
	from := fromEmail
	password := pwdEmail

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

	if _, err := os.Stat(pathToLoggerHash); os.IsNotExist(err) {
		err := os.Mkdir(pathToLoggerHash, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

	if _, err := os.Stat(pathToLoggerHash + filename + ".json"); os.IsNotExist(err) {
		file := pathToFolder + "/" + filename
		writeToJson(pathToFolder, filename, getHashSum(file), pathToLoggerHash)

		return false
	}

	return false
}

func writeToJson(pathToFolder string, fileName string, fileHash []byte, pathToLoggerHash string) {
	data := Log{
		pathToFolder: "pathToFolder",
		fileName:     "fileName",
		hash:         "string(fileHash)",
	}

	file, _ := json.MarshalIndent(data, "", "")
	path := pathToLoggerHash + "/" + fileName + ".json"
	_ = os.WriteFile(path, file, 0644)
}

func getHashSum(logFilePath string) []byte {
	file, err := os.Open(logFilePath)

	if err != nil {
		panic(err)
	}

	defer file.Close()

	hash := md5.New()
	_, err = io.Copy(hash, file)

	if err != nil {
		panic(err)
	}

	return hash.Sum(nil)
}
