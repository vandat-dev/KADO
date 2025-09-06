package initialize

import (
	"base_go_be/global"
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func checkErrPanic(err error, errStr string) {
	if err != nil {
		global.Logger.Error(errStr, zap.Error(err))
		println(err)
		panic(err)
	}
}

func Mysql() {
	/*
		refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	*/
	global.Logger.Info("Start connecting to mysql")
	m := global.Config.Mysql
	dsn := "%s:%s@tcp(%s:%v)/%s?charset=utf8mb4&parseTime=True&loc=Local"
	var s = fmt.Sprintf(dsn, m.UserName, m.Password, m.Host, m.Port, m.DBName)
	db, err := gorm.Open(mysql.Open(s), &gorm.Config{
		SkipDefaultTransaction: true, // turn on transaction
	})
	checkErrPanic(err, "Initialize MySQL database failed")
	global.Mysql = db
	global.Logger.Info("Mysql connect successfully!")
	setPool()
	// Migration is now handled manually, not automatically
	// migrateTables()
	//RemoveEmailColumn()
}

func setPool() {
	m := global.Config.Mysql
	sqlDB, err := global.Mysql.DB()
	checkErrPanic(err, "Set Pool MySQL database failed")
	sqlDB.SetMaxIdleConns(m.MaxIdleConn)
	sqlDB.SetMaxOpenConns(m.MaxOpenConn)
	sqlDB.SetConnMaxLifetime(time.Duration(m.ConnMaxLifeTime))
}

// Migration functions removed - migrations are now handled manually
//func migrateTables() {
//	err := global.Mysql.AutoMigrate(&model.User{}, &model.Product{})
//	checkErrPanic(err, "AutoMigrate MySQL database failed")
//}

//func RemoveEmailColumn() {
//	checkErrPanic(global.Mysql.Migrator().DropColumn(&po.User{}, "email"), "Remove email column")
//}
