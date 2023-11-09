package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"ticket-watcher/domain"
	"time"

	"github.com/erfanmomeniii/jalali"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger *logrus.Logger

func init() {
	lumberjackLogger := &lumberjack.Logger{
		Filename:   filepath.ToSlash("./logs/" + time.Now().Format("2006-01-02") + "_app.log"),
		MaxSize:    1, // MB
		MaxBackups: 10,
		MaxAge:     30,   // days
		Compress:   true, // disabled by default
	}

	multiWriter := io.MultiWriter(os.Stderr, lumberjackLogger)

	logFormatter := new(logrus.TextFormatter)
	logFormatter.TimestampFormat = time.RFC3339
	logFormatter.FullTimestamp = true

	Logger = logrus.New()
	Logger.SetFormatter(logFormatter)
	Logger.SetLevel(logrus.InfoLevel)
	Logger.SetOutput(multiWriter)
}

func GetJalaliDate(gregorianDate string) string {
	t, err := time.Parse("2006-01-02", gregorianDate)
	if err != nil {
		Logger.Error("Error:", err)
		return ""
	}

	// Convert to Jalali
	j := jalali.ConvertGregorianToJalali(t)

	jalaliDate := fmt.Sprintf(`%d-%d-%d`, j.Year(), j.Month(), j.Day())

	return jalaliDate
}

func GenerateUniqueID() string {
	timestamp := time.Now().UnixNano()
	randomNum := rand.Intn(100000)
	uniqueID := fmt.Sprintf("%d%d", timestamp, randomNum)
	return uniqueID
}

func ReadTravelsData() (travels []domain.Travel) {
	jsonData, err := ioutil.ReadFile("data.json")
	if err != nil {
		file, err := os.Create("data.json")
		if err != nil {
			Logger.Error("Error createing file:", err)
			return
		}
		defer file.Close()
		return
	}

	err = json.Unmarshal(jsonData, &travels)
	if err != nil {
		Logger.Error("Error decoding JSON:", err)
	}
	return
}

func StoreTravelsData(travels []domain.Travel) {
	fileContent, err := json.Marshal(travels)
	if err != nil {
		Logger.Error("Some error occurred while encoding travels: ", err)
		return
	}

	err = ioutil.WriteFile("data.json", fileContent, 0644)
	if err != nil {
		Logger.Error("Some error occurred while storing travels: ", err)
		return
	}
}
