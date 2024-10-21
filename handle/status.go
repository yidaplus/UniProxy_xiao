package handle

import (
	"os"
	"path"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/wyx2685/UniProxy/common/balance"
	"github.com/wyx2685/UniProxy/conf"
	"github.com/wyx2685/UniProxy/proxy"
	"github.com/wyx2685/UniProxy/v2b"
	"github.com/wyx2685/UniProxy/common/encrypt"
)

type initParamsRequest struct {
	MixedPort int    `json:"mixed_port"`
	AppName   string `json:"app_name"`
	Url       string `json:"url"`
	Token     string `json:"token"`
	License   string `json:"license"`
	UserPath  string `json:"user_path"`
}

var inited bool

func InitParams(c *gin.Context) {
	p := initParamsRequest{}
	err := c.ShouldBindJSON(&p)
	if err != nil {
		c.JSON(400, &Rsp{Success: false, Message: err.Error()})
		return
	}
	if encrypt.Sha([]byte(encrypt.Sha([]byte(p.Url+p.Url[5:]))+"大傻逼")) != p.License {
		c.JSON(400, &Rsp{Success: false})
		return
	}
	f, err := os.OpenFile(path.Join(p.UserPath, "uniproxy.log"), os.O_TRUNC|os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		c.JSON(400, &Rsp{Success: false, Message: err.Error()})
		return
	}
	log.SetOutput(f)
	if len(conf.C.Api.Baseurl) == 0 {
		conf.C.Api.Baseurl = []string{p.Url}
	}
	urlBalance = balance.New[string](conf.C.Api.Balance, conf.C.Api.Baseurl)
	v2b.Init(conf.C.Api.Balance, conf.C.Api.Baseurl, p.Token)
	proxy.InPort = p.MixedPort
	proxy.DataPath = p.UserPath
	servers = make(map[string]*v2b.ServerInfo)
	inited = true
	c.JSON(200, &Rsp{Success: true})
}

func GetStatus(c *gin.Context) {
	c.JSON(200, &Rsp{
		Success: true,
		Data: StatusData{
			Inited:      inited,
			Running:     proxy.Running,
			GlobalMode:  proxy.GlobalMode,
			SystemProxy: proxy.SystemProxy,
		},
	})
}
