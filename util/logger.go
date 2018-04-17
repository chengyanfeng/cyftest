package utils

import (
	"bufio"
	"fmt"
	"os"

	"github.com/astaxie/beego"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

func init() {
	Logger = logrus.New()
	Logger.SetLevel(logrus.DebugLevel)
	logger_root := beego.AppConfig.DefaultString("logger_path", "logs/")
	pathMap := lfshook.PathMap{
		logrus.InfoLevel:  logger_root + "info.log",
		logrus.ErrorLevel: logger_root + "error.log",
		logrus.DebugLevel: logger_root + "debug.log",
	}
	hook := lfshook.NewHook(pathMap, &logrus.JSONFormatter{})
	Logger.AddHook(hook)
}

func GetLogger(subject string) *logrus.Logger {
	logger := logrus.New()
	src, err := os.OpenFile(os.DevNull, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Println("err", err)
	}
	logger.Out = bufio.NewWriter(src)
	logger.SetLevel(logrus.DebugLevel)
	logger_root := beego.AppConfig.DefaultString("logger_path", "logs/")
	pathMap := lfshook.PathMap{
		logrus.InfoLevel:  logger_root + subject + ".log",
		logrus.ErrorLevel: logger_root + subject + ".log",
		logrus.DebugLevel: logger_root + subject + ".log",
	}
	hook := lfshook.NewHook(pathMap, &logrus.JSONFormatter{})
	logger.AddHook(hook)
	return logger
}

func Info(v ...interface{}) {
	Logger.Info(v)
}

func Debug(v ...interface{}) {
	Logger.Debug(v)
}

func Error(v ...interface{}) {
	Logger.Error(v)
}

func Subject(subject string, v ...interface{}) {
	Logger.WithFields(logrus.Fields{
		"subject": subject,
	}).Info(v)
}
