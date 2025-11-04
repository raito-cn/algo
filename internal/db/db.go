package db

import (
	"algo/internal/model"
	"algo/internal/util"
	"algo/pkg/config"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
	"sync"
	"time"
)

var db *gorm.DB
var once sync.Once

func initDB(debug bool) {
	var log = util.GetLog()
	var err error
	level := logger.Silent
	if debug {
		level = logger.Info
	}
	datasource := config.GetConfig().Dir.Datasource
	if err = os.MkdirAll(datasource, 0755); err != nil {
		util.GetLog().Error("create code dir error", zap.Error(err))
		panic(err)
	}

	db, err = gorm.Open(sqlite.Open("file:"+datasource+"/algo.db"), &gorm.Config{
		Logger: logger.Default.LogMode(level),
	})
	if err != nil {
		log.Error("failed to connect database", zap.Error(err))
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Error("failed to get sql db", zap.Error(err))
	}

	sqlDB.SetMaxOpenConns(1)
	sqlDB.SetMaxIdleConns(1)
	sqlDB.SetConnMaxLifetime(time.Hour)
	sqlDB.SetConnMaxIdleTime(time.Minute * 30)

	err = db.AutoMigrate(&model.Problem{}, &model.Tag{}, &model.Contest{})
	if err != nil {
		log.Error("failed to auto migrate", zap.Error(err))
	}
}

func GetDB(debug bool) *gorm.DB {
	once.Do(func() {
		initDB(debug)
	})
	return db
}
