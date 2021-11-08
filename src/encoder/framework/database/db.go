package database

import (
	"encoder/domain"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/lib/pq"
)

const ENV_TESTE = "test"
const DBTYPETEST = "sqlite3"
const DSNTEST = ":memory:"

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
	dnInstace.Env = ENV_TESTE
	dnInstace.DbTypeTest = DBTYPETEST
	dnInstace.DsnTest = DSNTEST

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

	if d.Env != ENV_TESTE {
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
	}
	return d.Db, nil

}
