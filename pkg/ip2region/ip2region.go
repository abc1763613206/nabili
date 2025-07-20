package ip2region

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"path/filepath"

	"github.com/abc1763613206/nabili/internal/constant"
	"github.com/abc1763613206/nabili/pkg/download"
	"github.com/abc1763613206/nabili/pkg/wry"

	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
)

var DownloadUrls = []string{
	"https://cdn.jsdelivr.net/gh/lionsoul2014/ip2region/data/ip2region.xdb",
	"https://raw.githubusercontent.com/lionsoul2014/ip2region/master/data/ip2region.xdb",
}

type Ip2Region struct {
	seacher *xdb.Searcher
}

func NewIp2Region(filePath string) (*Ip2Region, error) {
	_, err := os.Stat(filePath)
	if err != nil && os.IsNotExist(err) {
		log.Println("文件不存在，尝试从网络获取最新 ip2region 库")
		_, err = download.Download(filePath, DownloadUrls...)
		if err != nil {
			log.Printf("❌ ip2region 数据库下载失败！\n")
			log.Printf("📁 请手动下载并保存到: %s\n", filepath.Join(constant.DataDirPath, "ip2region.xdb"))
			log.Printf("🔗 下载地址: %v\n", DownloadUrls)
			log.Printf("💡 操作步骤:\n")
			log.Printf("   1. 从上述链接下载 ip2region.xdb 文件\n")
			log.Printf("   2. 将下载的文件重命名为: ip2region.xdb\n")
			log.Printf("   3. 复制到数据目录: %s\n", constant.DataDirPath)
			log.Printf("   4. 重新运行 nabili\n")
			return nil, err
		}
	}

	f, err := os.OpenFile(filePath, os.O_RDONLY, 0400)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	searcher, err := xdb.NewWithBuffer(data)
	if err != nil {
		fmt.Printf("无法解析 ip2region xdb 数据库: %s\n", err)
		return nil, err
	}
	return &Ip2Region{
		seacher: searcher,
	}, nil
}

func (db Ip2Region) Find(query string, params ...string) (result fmt.Stringer, err error) {
	if db.seacher != nil {
		res, err := db.seacher.SearchByStr(query)
		if err != nil {
			return nil, err
		} else {
			return wry.Result{
				Country: strings.ReplaceAll(res, "|0", ""),
			}, nil
		}
	}

	return nil, errors.New("ip2region 未初始化")
}

func (db Ip2Region) Name() string {
	return "ip2region"
}
