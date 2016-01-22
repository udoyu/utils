package simini

import (
	"os"
	"io/ioutil"
	"testing"
)

var (
	_testIni SimIni
	inidata = `
	[string]
	key=val
	key1 = val1
		key2=	val2	
	#key3=val3
	[int]
	key=123
	key1=0
	`
)

func testGetString(t *testing.T, ini SimIni, sess, key, val string) {
	v := ini.GetStringVal(sess, key)
	if v != val {
		t.Errorf("sess=%s,key=%s,val=%s,v=%s", 
			sess, key, val, v)
	} 
}

func testNotes(t *testing.T, ini SimIni, sess, key string) {
	if ini.GetStringVal(sess, key) != "" {
		t.Error("testNotes failed")
	}
}

func testGetInt(t *testing.T, ini SimIni, sess, key string, val int) {
	v,e := ini.GetIntVal(sess, key)
	if e != nil || v != val {
		t.Errorf("sess=%s,key=%s,val=%d,v=%d,err=%v", 
			sess, key, val, v, e)
	}
}

func testGetStringWithDefault(t *testing.T, 
							  ini SimIni, sess, key, defaultStr, val string) {
	v := ini.GetStringValWithDefault(sess, key, defaultStr)
	if v != val {
		t.Errorf("sess=%s,key=%s,val=%s,defaultStr=%s,v=%s",
			sess, key, val, defaultStr, v)
	}
}

func testGetIntWithDefault(t *testing.T, 
							  ini SimIni, sess, key string, defaultInt, val int) {
	v, e := ini.GetIntValWithDefault(sess, key, defaultInt)
	if e != nil || v != val {
		t.Errorf("sess=%s,key=%s,val=%d,defaultInt=%d,v=%d",
			sess, key, val, defaultInt, v, e)
	}
}

func TestLoadFile(t *testing.T) {
	filename := "test.ini"
	if err := ioutil.WriteFile(filename, []byte(inidata), os.ModePerm); err != nil {
		t.Fatal(err)
	}
	defer os.Remove(filename)
	if 0 != _testIni.LoadFile(filename) {
		t.Fatal(_testIni.ErrMsg())
	}
}

func TestGetString(t *testing.T) {
	testGetString(t, _testIni, "string", "key", "val")
	testGetString(t, _testIni, "string", "key1", "val1")
	testGetString(t, _testIni, "string", "key2", "val2")
}

func TestNotes(t *testing.T) {
	testNotes(t, _testIni, "string", "key3")
}

func TestGetInt(t *testing.T) {
	testGetInt(t, _testIni, "int", "key", 123)
}

func TestGetSTringWithDefault(t *testing.T) {
	testGetStringWithDefault(t, _testIni, "string", "key", "default", "val")
	testGetStringWithDefault(t, _testIni, "string", "kkk", "default", "default")
}

func TestGetIntWithDefault(t *testing.T) {
	testGetIntWithDefault(t, _testIni, "int", "kkk", 111, 111)
	testGetIntWithDefault(t, _testIni, "int", "key", 123, 123)
	testGetIntWithDefault(t, _testIni, "int", "key1", 2, 0)
}