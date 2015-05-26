package simsql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type SimSql struct {
	DB      *sql.DB
	User    string
	Passwd  string
	Host    string
	Dbname  string
	Charset string
}

func (p *SimSql) Exec(exestr string) string {
	if nil == p.DB {
		return ""
	}
	_, err := p.DB.Exec(exestr)
	if err != nil {
		fmt.Println("g_DB.Exec|", err.Error())
		return ""
	}
	return exestr
}

func (p *SimSql) Init() error {
	mysqlUrl := ""
	if "" != p.Dbname {
		mysqlUrl = p.User + ":" + p.Passwd + "@tcp(" + p.Host + ")/" + p.Dbname + "?charset=" + p.Charset
	} else {
		mysqlUrl = p.User + ":" + p.Passwd + "@tcp(" + p.Host + ")/" + "information_schema" + "?charset" + p.Charset
	}
	var err error
	p.DB, err = sql.Open("mysql", mysqlUrl)
	if err != nil {
		fmt.Printf("sql.Open|", err.Error())
		return err
	}
	return nil
}

func (p *SimSql) Close() {
	if nil != p.DB {
		p.DB.Close()
	}
}

func (p *SimSql) Copy() *SimSql {
	s := &SimSql{}
	s.User = p.User
	s.Passwd = p.Passwd
	s.Host = p.Host
	s.Dbname = p.Dbname
	s.Charset = p.Charset
	return s
}
