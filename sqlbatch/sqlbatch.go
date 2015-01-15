package main

import (
	"flag"
	"fmt"
	"github.com/udoyu/utils/simini"
	"github.com/udoyu/utils/simsql"
	"strconv"
	"strings"
)

type DbConfInfo struct {
	SplitFlag int
	Db        *simsql.SimSql
	SqlStr    string
	DbPrefix  string
	Infos     []*DbSplitInfo
	Tables    []*TableConfInfo
}

type DbSplitInfo struct {
	Db         *simsql.SimSql
	StartIndex int
	EndIndex   int
}

type TableConfInfo struct {
	SplitFlag   int
	TablePrefix string
	SqlStr      string
	StartIndex  int
	EndIndex    int
}

const (
	DBNAME    = "$DbName"
	TABLENAME = "$TableName"
)

type SqlDbMap map[string]*DbConfInfo

var g_dbmap SqlDbMap

func main() {
	flag.Parse()
	if 1 != flag.NArg() {

		fmt.Printf("Usage:pp_server conf_file\n", flag.Arg(0))
		return
	}
	ini := new(simini.SimIni)
	if 0 != ini.LoadFileExtern(flag.Arg(0)) {
		fmt.Println("errmsg=", ini.ErrMsg())
	}
	DbLoadConf(ini, &g_dbmap)
	TableLoadConf(ini, &g_dbmap)
	ExecDb(&g_dbmap)
}

func ParseHoststrs(Db *simsql.SimSql, hoststr string) {
	s := strings.Split(hoststr, "|")
	Db.User = s[2]
	Db.Passwd = s[3]
	Db.Host = s[0] + ":" + s[1]
	Db.Dbname = s[4]
	Db.Charset = s[5]
}

func DbLoadConf(ini *simini.SimIni, dbmap *SqlDbMap) {
	*dbmap = make(map[string]*DbConfInfo)
	sessmap := ini.GetAllSession()
	for k, v := range sessmap {
		fmt.Println("k=", k)
		if k[0:3] != "db_" {
			continue
		}
		ci := new(DbConfInfo)
		ci.DbPrefix = ini.GetStringVal(k, "db_prefix")
		ci.SqlStr = ini.GetStringVal(k, "sql_str")
		ci.SplitFlag, _ = ini.GetIntVal(k, "split_flag")
		fmt.Println("k=", k, "|split_flag=", ci.SplitFlag)
		if 0 == ci.SplitFlag {
			ci.Db = new(simsql.SimSql)
			hoststr := ini.GetStringVal(k, "host")
			if hoststr == "" {
				continue
			}
			ParseHoststrs(ci.Db, hoststr)
		} else {
			for vk, vv := range v {
				fmt.Println("vk=", vk, "|vv=", vv)
				s := strings.Split(vk, "-")
				if 2 != len(s) {
					continue
				}
				info := new(DbSplitInfo)
				info.StartIndex, _ = strconv.Atoi(s[0])
				info.EndIndex, _ = strconv.Atoi(s[1])
				info.Db = new(simsql.SimSql)
				ParseHoststrs(info.Db, vv)
				ci.Infos = append(ci.Infos, info)
			}
		}
		(*dbmap)[k] = ci
	}
}

func TableLoadConf(ini *simini.SimIni, dbmap *SqlDbMap) {
	sessmap := ini.GetAllSession()
	for k, _ := range sessmap {
		if k[0:6] != "table_" {
			continue
		}
		dbkey := ini.GetStringVal(k, "db")
		db, ok := (*dbmap)[dbkey]
		if !ok {
			continue
		}
		info := new(TableConfInfo)
		db.Tables = append(db.Tables, info)
		info.SplitFlag, _ = ini.GetIntValWithDefault(k, "split_flag", 0)
		info.TablePrefix = ini.GetStringVal(k, "table_prefix")
		info.SqlStr = ini.GetStringVal(k, "sql_str")
		if 0 != info.SplitFlag {
			info.StartIndex, _ = ini.GetIntVal(k, "start_i")
			info.EndIndex, _ = ini.GetIntVal(k, "end_i")
		}
	}
}

func ExecDb(dbmap *SqlDbMap) {
	for _, v := range *dbmap {
		if 0 == v.SplitFlag {
			func() {
				if nil != v.Db.Init() {
					return
				}
				defer v.Db.Close()
				ExecSql(v.Db, v.SqlStr, v.DbPrefix, "")
				for _, tbinfo := range v.Tables {
					ExecTable(v.Db, tbinfo, v.DbPrefix)
				}
			}()
		} else {
			for _, info := range v.Infos {
				func() {
					if nil != info.Db.Init() {
						return
					}
					defer info.Db.Close()
					for i := info.StartIndex; i <= info.EndIndex; i++ {
						dbname := v.DbPrefix + fmt.Sprintf("%d", i)
						ExecSql(info.Db, v.SqlStr, dbname, "")
						for _, tbinfo := range v.Tables {
							ExecTable(info.Db, tbinfo, dbname)
						}
					}
				}()
			}
		}
	}
}

func ExecTable(s *simsql.SimSql, info *TableConfInfo, dbname string) {
	s.Exec("USE " + dbname + ";")
	s.Exec("START TRANSACTION;")
	if 0 == info.SplitFlag {
		ExecSql(s, info.SqlStr, dbname, info.TablePrefix)
	} else {
		for i := info.StartIndex; i <= info.EndIndex; i++ {
			tablename := info.TablePrefix + fmt.Sprintf("%d", i)
			ExecSql(s, info.SqlStr, dbname, tablename)
		}
	}
	s.Exec("COMMIT;")
}

func ExecSql(s *simsql.SimSql, sqlstr, dbname, tablename string) {
	if "" == sqlstr {
		return
	}
	sqlstr = strings.Replace(sqlstr, DBNAME, dbname, -1)
	sqlstr = strings.Replace(sqlstr, TABLENAME, tablename, -1)
	retstr := s.Exec(sqlstr)
	Info(retstr)
}

func Info(str string) {
	fmt.Println(str)
}
