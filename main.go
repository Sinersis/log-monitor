package main

import (
	"crypto/md5"
	"encoding/json"
	"github.com/joho/godotenv"
	"io"
	"log"
	"net/smtp"
	"os"
)

type Log struct {
	pathToFolder, fileName, hash string
}

func init() {
	// loads values from .env into the system
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

	toEmailAddress, _ := os.LookupEnv("EMAIL_TO")
	to := []string{toEmailAddress}

	host, _ := os.LookupEnv("EMAIL_HOST")
	port, _ := os.LookupEnv("EMAIL_PORT")
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
