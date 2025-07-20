package zxipv6wry

import (
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/saracen/go7z"
	"github.com/abc1763613206/nabili/internal/constant"
	"github.com/abc1763613206/nabili/pkg/common"
)

func Download(filePath ...string) (data []byte, err error) {
	data, err = getData()
	if err != nil {
		log.Printf("❌ ZX IPv6数据库下载失败！\n")
		log.Printf("📁 请手动下载并保存到: %s\n", filepath.Join(constant.DataDirPath, "zxipv6wry.db"))
		log.Printf("🔗 下载地址: https://ip.zxinc.org/ip.7z\n")
		log.Printf("💡 操作步骤:\n")
		log.Printf("   1. 从上述链接下载 ip.7z 文件\n")
		log.Printf("   2. 解压文件，找到 zxipv6wry.db\n")
		log.Printf("   3. 将 zxipv6wry.db 复制到数据目录: %s\n", constant.DataDirPath)
		log.Printf("   4. 重新运行 nabili\n")
		return
	}

	if !CheckFile(data) {
		log.Printf("❌ ZX IPv6数据库下载出错！\n")
		log.Printf("📁 请重新下载并保存到: %s\n", filepath.Join(constant.DataDirPath, "zxipv6wry.db"))
		log.Printf("🔗 下载地址: https://ip.zxinc.org/ip.7z\n")
		log.Printf("💡 操作步骤:\n")
		log.Printf("   1. 从上述链接下载 ip.7z 文件\n")
		log.Printf("   2. 解压文件，找到 zxipv6wry.db\n")
		log.Printf("   3. 将 zxipv6wry.db 复制到数据目录: %s\n", constant.DataDirPath)
		log.Printf("   4. 重新运行 nabili\n")
		return nil, errors.New("数据库下载内容出错")
	}

	if len(filePath) == 1 {
		if err := common.SaveFile(filePath[0], data); err == nil {
			log.Println("已将最新的 ZX IPv6数据库 保存到本地:", filePath)
		}
	}
	return
}

const (
	zx = "https://ip.zxinc.org/ip.7z"
)

func getData() (data []byte, err error) {
	data, err = common.GetHttpClient().Get(zx)
	if err != nil {
		return nil, err
	}

	file7z, err := os.CreateTemp("", "*")
	if err != nil {
		return nil, err
	}
	defer os.Remove(file7z.Name())
	if err := os.WriteFile(file7z.Name(), data, 0644); err == nil {
		return Un7z(file7z.Name())
	}
	return
}

func Un7z(filePath string) (data []byte, err error) {
	sz, err := go7z.OpenReader(filePath)
	if err != nil {
		return nil, err
	}
	defer sz.Close()

	fileNoNeed, err := os.CreateTemp("", "*")
	if err != nil {
		return nil, err
	}
	fileNeed, err := os.CreateTemp("", "*")
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}
	for {
		hdr, err := sz.Next()
		if err == io.EOF {
			break // IdxEnd of archive
		}
		if err != nil {
			return nil, err
		}

		if hdr.Name == "ipv6wry.db" {
			if _, err := io.Copy(fileNeed, sz); err != nil {
				log.Fatalln("ZX ipv6数据库解压出错：", err.Error())
			}
		} else {
			if _, err := io.Copy(fileNoNeed, sz); err != nil {
				log.Fatalln("ZX ipv6数据库解压出错：", err.Error())
			}
		}
	}
	err = fileNoNeed.Close()
	if err != nil {
		return nil, err
	}
	defer os.Remove(fileNoNeed.Name())
	defer os.Remove(fileNeed.Name())
	return os.ReadFile(fileNeed.Name())
}
