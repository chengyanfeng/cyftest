package utils

import (
	"fmt"
	"time"

	"github.com/astaxie/beego"
	"github.com/sirupsen/logrus"
)

var JDBCLogger *logrus.Logger
var Jdbc_proxy_url string

func init() {
	JDBCLogger = GetLogger("jdbc")
	Jdbc_proxy_url = beego.AppConfig.DefaultString("jdbc_url", "http://127.0.0.1:4567/sql")
}

func JDBC(sql string, db P) (result string, err error) {
	db_config := JsonEncode(db)
	logger := JDBCLogger.WithFields(logrus.Fields{
		"sql": sql,
		"db":  db_config,
	})
	begin := time.Now()
	logger.Info("begin")
	result, err = HttpPost(Jdbc_proxy_url, nil, &P{"sql": sql, "db": db_config})
	if err != nil {
		logger.Error(err)
	}
	finish := time.Now()
	nanoseconds := finish.Sub(begin).Nanoseconds()
	milliseconds := fmt.Sprintf("%d.%d", nanoseconds/1e6, nanoseconds%1e6)
	logger.WithField("consume", milliseconds).Info("finish")
	return
}

func JDBCToCSV(sql string, db P) (result string, err error) {
	db_config := JsonEncode(db)
	logger := JDBCLogger.WithFields(logrus.Fields{
		"sql": sql,
		"db":  db_config,
	})
	begin := time.Now()
	logger.Info("begin")
	result, err = HttpPost(Jdbc_proxy_url, nil, &P{"sql": sql, "db": db_config, "o": "csv"})
	if err != nil {
		logger.Error(err)
	}
	finish := time.Now()
	nanoseconds := finish.Sub(begin).Nanoseconds()
	milliseconds := fmt.Sprintf("%d.%d", nanoseconds/1e6, nanoseconds%1e6)
	logger.WithField("consume", milliseconds).Info("finish")
	return
}
