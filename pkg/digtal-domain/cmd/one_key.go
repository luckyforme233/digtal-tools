package cmd

import (
	"digtal/pkg/digtal-domain/config"
	"digtal/src/cloudflare"
	"digtal/src/digtal"
	ssh2 "digtal/src/ssh"
	"fmt"
	"github.com/spf13/cobra"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
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

		dns, _ := cloudflare.CreateDns(config.C.CLDomain, ip)
		domian := strings.ToLower(dns)
		fmt.Println("60 s 域名解析：", domian)
		time.Sleep(time.Minute * 1)
		var client *ssh2.SshClient
		fmt.Println("等待服务器准备完成，进行链接。。。。。。")

		for {
			client, err = ssh2.NewSshClient("root", ip, 22, config.C.PrvKeyPath, "")
			if err != nil {
				log.Println("链接 失败: ", err)
				time.Sleep(time.Second * 3)
			} else {
				fmt.Println("链接成功")
				break
			}
		}

		shellContent, err := ioutil.ReadFile("xray.sh")
		if err != nil {
			log.Fatalf("unable to read private key: %v", err)
			return
		}

		replace := strings.Replace(string(shellContent), "REPLACE_DOMAIN", domian, -1)
		localFileName := "./temp.sh"
		err = ioutil.WriteFile(localFileName, []byte(replace), fs.ModePerm)
		if err != nil {
			log.Fatalf("写入文件错误: %v", err)
			return
		}

		err = client.ScpCopy(localFileName, "/root/")
		if err != nil {
			log.Fatalf("拷贝文件错误: %v", err)
			return
		}

		output, err := client.RunCommand("chmod +x ~/temp.sh && bash ~/temp.sh")
		if err != nil {
			log.Println("安装失败：", err, "result: ", output)
			return
		}
		log.Println("安装结果：", output)
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		digtal.InitClient()
	},
}

// 最终方案-全兼容
func getCurrentAbPath() string {
	dir := getCurrentAbPathByExecutable()
	tmpDir, _ := filepath.EvalSymlinks(os.TempDir())
	if strings.Contains(dir, tmpDir) {
		return getCurrentAbPathByCaller()
	}
	return dir
}

// 获取当前执行文件绝对路径
func getCurrentAbPathByExecutable() string {
	exePath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	res, _ := filepath.EvalSymlinks(filepath.Dir(exePath))
	return res
}

// 获取当前执行文件绝对路径（go run）
func getCurrentAbPathByCaller() string {
	var abPath string
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		abPath = path.Dir(filename)
	}
	return abPath
}
