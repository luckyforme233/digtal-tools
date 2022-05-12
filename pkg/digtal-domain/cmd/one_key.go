package cmd

import (
	"digtal/pkg/digtal-domain/config"
	"digtal/src/digtal"
	ssh2 "digtal/src/ssh"
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"time"
)

var one_key = &cobra.Command{
	Use:   "nb",
	Short: "一键部署",
	Long:  `一键部署`,
	Run: func(cmd *cobra.Command, args []string) {
		droplet, _, err := digtal.CreateDroplet()
		if err != nil {
			log.Println("创建vps 失败: ", err)
			return
		}
		var ip string

		for {
			info, _, _ := digtal.GetDropLetInfo(droplet.ID)
			if len(info.Networks.V4) > 0 {
				ip = info.Networks.V4[0].IPAddress
				break
			} else {
				time.Sleep(time.Second * 3)
			}
		}
		fmt.Println("ip:", ip)

		var client *ssh2.SshClient

		for {
			client, err = ssh2.NewSshClient("root", ip, 22, config.C.PrvKeyPath, "")
			if err != nil {
				log.Println("链接 失败: ", err)
				time.Sleep(time.Second * 3)
			} else {
				break
			}

			fmt.Println(client)
			fmt.Println(err)
		}

		fmt.Println("服务器已创建，等待服务器链接...")
		time.Sleep(time.Minute * 3)

		output, err := client.RunCommand("bash <(curl -L https://raw.githubusercontent.com/v2fly/fhs-install-v2ray/master/install-release.sh)")
		if err != nil {
			log.Println("安装失败：", err)
			return
		}
		log.Println("安装结果：", output)

		// detect
		testV2ray, err := client.RunCommand("/usr/local/bin/v2ray -test -config /usr/local/etc/v2ray/config.json")
		if err != nil {
			log.Println("检测失败：", err)
			return
		}

		log.Println("检测结果", testV2ray)

		jsonPlainText, err := ioutil.ReadFile("tt.json")
		if err != nil {
			log.Fatalf("unable to read private key: %v", err)
			return
		}

		runCmd := fmt.Sprintf(`echo '%s' >  /usr/local/etc/v2ray/config.json`, jsonPlainText)
		_, err = client.RunCommand(runCmd)
		if err != nil {
			log.Println("检测失败：", err)
			return
		}

		add2, err := client.RunCommand("cat /usr/local/etc/v2ray/config.json")
		if err != nil {
			log.Println("检测失败：", err)
			return
		}

		log.Println("配置文件", add2)

		client.RunCommand("service v2ray restart")

	},
	PreRun: func(cmd *cobra.Command, args []string) {
		digtal.InitClient()
	},
}
