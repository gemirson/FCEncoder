package database

import (
	"encoder/domain"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/lib/pq"
)

type DataBase struct {
	Db             *gorm.DB
	Dsn            string
	DsnTest        string
	DbType         string
	DbTypeTest     string
	Debug          bool
	AutoMigratedDb bool
	Env            string
}

func NewDb() *DataBase {
	return &DataBase{}

}

func NewDbTest() *gorm.DB {

	dnInstace := NewDb()
	dnInstace.Env = "test"
	dnInstace.DbTypeTest = "sqlite3"
	dnInstace.DsnTest = ":memory:"

	dnInstace.AutoMigratedDb = true
	dnInstace.Debug = true

	connection, err := dnInstace.Connect()

	if err != nil {
		log.Fatalf("Test db error: %v", err)

	}

	return connection

}

func (d *DataBase) Connect() (*gorm.DB, error) {

	var err error

	if d.Env != "test" {
		d.Db, err = gorm.Open(d.DbType, d.Dsn)
	} else {
		d.Db, err = gorm.Open(d.DbTypeTest, d.DsnTest)
	}

	if err != nil {
		return nil, err
	}

	if d.Debug {
		d.Db.LogMode(true)
	}

	if d.AutoMigratedDb {
		d.Db.AutoMigrate(&domain.Video{}, &domain.Job{})
		d.Db.Model(domain.Job{}).AddForeignKey("video_id", "videos (id)", "CASCADE", "CASCADE")
		d.Db.Raw("PRAGMA foreign_keys=ON")
	}
	return d.Db, nil

}
