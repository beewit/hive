package global

import (
	"encoding/json"
	"fmt"

	"github.com/beewit/beekit/conf"
	"github.com/beewit/beekit/log"
	"github.com/beewit/beekit/mysql"
	"github.com/beewit/beekit/redis"
	"github.com/beewit/beekit/utils/convert"
	"strings"
)

var (
	CFG  = conf.New("config.json")
	Log  = log.Logger
	DB   = mysql.DB
	RD   = redis.Cache
	IP   = CFG.Get("server.ip")
	Port = CFG.Get("server.port")
	Host = fmt.Sprintf("http://%v:%v", IP, Port)

	FileConf = &fileConf{
		BasePath: convert.ToString(CFG.Get("files.basePath")),
		Path:     convert.ToString(CFG.Get("files.path")),
		DoMain:   convert.ToString(CFG.Get("files.doMain")),
	}
)

const (
	PAGE_SIZE = 10
)

type fileConf struct {
	BasePath string
	Path     string
	DoMain   string
}

type Account struct {
	ID       int64  `json:"id"`
	Nickname string `json:"nickname"`
	Photo    string `json:"photo"`
	Mobile   string `json:"mobile"`
	Status   string `json:"status"`
	OrgId    int64  `json:"org_id"`
}

//组织
type Org struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	ParentId  int64  `json:"parent_id"`
	Relation  string `json:"relation"`
	AccountId int64  `json:"account_id"`
	Type      string `json:"type"`
}

func ToByteAccount(b []byte) *Account {
	var rp = new(Account)
	err := json.Unmarshal(b[:], &rp)
	if err != nil {
		Log.Error(err.Error())
		return nil
	}
	return rp
}

func ToMapAccount(m map[string]interface{}) *Account {
	b := convert.ToMapByte(m)
	if b == nil {
		return nil
	}
	return ToByteAccount(b)
}

func ToInterfaceAccount(m interface{}) *Account {
	b := convert.ToInterfaceByte(m)
	if b == nil {
		return nil
	}
	return ToByteAccount(b)
}


func ToByteOrg(b []byte) *Org {
	var rp = new(Org)
	err := json.Unmarshal(b[:], &rp)
	if err != nil {
		Log.Error(err.Error())
		return nil
	}
	return rp
}

func ToMapOrg(m map[string]interface{}) *Org {
	b := convert.ToMapByte(m)
	if b == nil {
		return nil
	}
	return ToByteOrg(b)
}

func ToInterfaceOrg(m interface{}) *Org {
	b := convert.ToInterfaceByte(m)
	if b == nil {
		return nil
	}
	return ToByteOrg(b)
}

func GetSavePath(path string) string {
	return strings.Replace(path, FileConf.BasePath, "", -1)
}
