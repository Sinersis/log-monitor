package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
	"io"
	"log"
	"os"
)

type LogStruct struct {
	PathToFolder, FileName, Hash string
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {

	pathToFolder, _ := os.LookupEnv("PATH_TO_LOG_FOLDER")
	fileName, _ := os.LookupEnv("FILE_NAME")
	pathToLoggerHash, _ := os.LookupEnv("PATH_TO_LOGGER_HASH_FOLDER")

	if loggerHashChecker(pathToLoggerHash, pathToFolder, fileName) {
		emailSend(fileName)
	}

}

func emailSend(fileName string) {

	from, _ := os.LookupEnv("EMAIL_SENDER")
	password, _ := os.LookupEnv("EMAIL_APP_PWD")
	processId, _ := os.LookupEnv("PROCESS_ID")

	toEmailAddress, _ := os.LookupEnv("EMAIL_TO")

	subject := "Что-то пошло не так: Лог процесса" + processId + " поменял свое содержимое \n"
	body := "Вы получили это письмо так как что-то случилось с сервисом: " + processId

	msg := gomail.NewMessage()
	msg.SetHeader("From", from)
	msg.SetHeader("To", toEmailAddress)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", body)

	n := gomail.NewDialer("smtp.gmail.com", 587, from, password)

	// Send the email
	if err := n.DialAndSend(msg); err != nil {
		panic(err)
	}
}

func loggerHashChecker(pathToLoggerHash string, pathToFolder string, filename string) bool {
	file := pathToFolder + "/" + filename
	if _, err := os.Stat(pathToLoggerHash); os.IsNotExist(err) {
		err := os.Mkdir(pathToLoggerHash, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

	if _, err := os.Stat(pathToLoggerHash + "/" + filename + ".json"); os.IsNotExist(err) {
		writeToJson(pathToFolder, filename, getHashSum(file), pathToLoggerHash)
		return false
	}

	if getHashSum(file) != readFromJson(filename, pathToLoggerHash) {
		writeToJson(pathToFolder, filename, getHashSum(file), pathToLoggerHash)
		return true
	} else {
		return false
	}
}

func readFromJson(fileName string, pathToLoggerHash string) string {

	var logger LogStruct
	path := pathToLoggerHash + "/" + fileName + ".json"

	jsonFile, err := os.Open(path)

	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	err = json.Unmarshal(byteValue, &logger)

	if err != nil {
		return "Halt"
	}

	return logger.Hash
}

func writeToJson(pathToFolder string, fileName string, fileHash string, pathToLoggerHash string) {
	data := LogStruct{
		PathToFolder: pathToFolder,
		FileName:     fileName,
		Hash:         fileHash,
	}

	file, err := json.Marshal(data)

	if err != nil {
		panic(err)
	}

	path := pathToLoggerHash + "/" + fileName + ".json"
	_ = os.WriteFile(path, file, os.ModePerm)
}

func getHashSum(logFilePath string) string {
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

	return hex.EncodeToString(hash.Sum(nil))
}
