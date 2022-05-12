package digtal

import (
	"digtal/pkg/digtal-domain/config"
	"github.com/digitalocean/godo"
	"log"
)

var client *godo.Client

func InitClient() {
	if config.C.Token == "" {
		log.Fatalln("token 未设置")
	}
	client = godo.NewFromToken(config.C.Token)
}
