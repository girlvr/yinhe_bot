package model

import (
	"fmt"
	"github.com/go-xorm/xorm"
	"github.com/google/uuid"
	"github.com/pelletier/go-toml"
	"net/url"
	"reflect"
	"time"
)

var db *xorm.Engine
var syncTable map[string]interface{}

// RegisterTable ...
func RegisterTable(v interface{}) {
	tof := reflect.TypeOf(v).Name()
	if syncTable == nil {
		syncTable = map[string]interface{}{
			tof: v,
		}
	}
	syncTable[tof] = v
}

// DB ...
func DB() *xorm.Engine {
	if db == nil {
		if err := InitDB(); err != nil {
			panic(err)
		}
	}
	return db
}

// Database ...
type Database struct {
	ShowSQL  bool   `toml:"show_sql"`
	UseCache bool   `json:"use_cache"`
	Type     string `toml:"type"`
	Addr     string `toml:"addr"`
	Port     string `toml:"port"`
	Username string `toml:"username"`
	Password string `toml:"password"`
	Schema   string `toml:"schema"`
	Location string `toml:"location"`
	Charset  string `toml:"charset"`
	Prefix   string `toml:"prefix"`
}

// DefaultDB ...
func DefaultDB() *Database {
	return &Database{
		ShowSQL:  true,
		UseCache: true,
		Type:     "mysql",
		Addr:     "localhost",
		Port:     "3306",
		Username: "root",
		Password: "111111",
		Schema:   "yinhe",
		Location: url.QueryEscape("Asia/Shanghai"),
		Charset:  "utf8mb4",
		Prefix:   "",
	}
}

// Source ...
func (d *Database) Source() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?loc=%s&charset=%s&parseTime=true",
		d.Username, d.Password, d.Addr, d.Port, d.Schema, d.Location, d.Charset)
}

// InitDB ...
func InitDB(pathname string) (e error) {

	eng, e := xorm.NewEngine("mysql", LoadToml(pathname).Source())
	if e != nil {
		return e
	}
	eng.ShowSQL(true)
	eng.ShowExecTime(true)

	for _, val := range syncTable {
		e := eng.Sync2(val)
		if e != nil {
			return e
		}
	}

	db = eng
	return nil
}

// LoadToml ...
func LoadToml(path string) (db *Database) {
	db = DefaultDB()
	tree, err := toml.LoadFile(path)
	if err != nil {
		return db
	}
	err = tree.Unmarshal(db)
	if err != nil {
		return db
	}
	return db
}

// Model ...
type Model struct {
	ID        string     `json:"-" xorm:"id pk"`
	CreatedAt time.Time  `json:"-" xorm:"created_at created"`
	UpdatedAt time.Time  `json:"-" xorm:"updated_at updated"`
	DeletedAt *time.Time `json:"-" xorm:"deleted_at deleted"`
	//Version   int        `json:"-" xorm:"version"`
}

// BeforeInsert ...
func (m *Model) BeforeInsert() {
	if m.ID == "" {
		m.ID = uuid.Must(uuid.NewRandom()).String()
	}
}
