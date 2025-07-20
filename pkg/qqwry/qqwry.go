package qqwry

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"path/filepath"

	"github.com/abc1763613206/nabili/internal/constant"
	"github.com/abc1763613206/nabili/pkg/download"
	"github.com/abc1763613206/nabili/pkg/wry"
)

var DownloadUrls = []string{
	"https://github.com/metowolf/qqwry.dat/releases/latest/download/qqwry.dat",
	// Other repo:
	// https://github.com/HMBSbige/qqwry // This repository has been archived since Jun 27, 2024.
	// https://github.com/FW27623/qqwry // This repository's dat format will not be maintained after October 2024.
	// https://github.com/metowolf/qqwry.dat
}

type QQwry struct {
	wry.IPDB[uint32]
}

// NewQQwry new database from path
func NewQQwry(filePath string) (*QQwry, error) {
	var fileData []byte

	_, err := os.Stat(filePath)
	if err != nil && os.IsNotExist(err) {
		log.Println("文件不存在，尝试从网络获取最新纯真 IP 库")
		fileData, err = download.Download(filePath, DownloadUrls...)
		if err != nil {
			log.Printf("❌ 纯真 IP 库下载失败！\n")
			log.Printf("📁 请手动下载并保存到: %s\n", filepath.Join(constant.DataDirPath, "qqwry.dat"))
			log.Printf("🔗 下载地址: %v\n", DownloadUrls)
			log.Printf("💡 操作步骤:\n")
			log.Printf("   1. 从上述链接下载 qqwry.dat 文件\n")
			log.Printf("   2. 将下载的文件重命名为: qqwry.dat\n")
			log.Printf("   3. 复制到数据目录: %s\n", constant.DataDirPath)
			log.Printf("   4. 重新运行 nabili\n")
			return nil, err
		}
	} else {
		fileBase, err := os.OpenFile(filePath, os.O_RDONLY, 0400)
		if err != nil {
			return nil, err
		}
		defer fileBase.Close()

		fileData, err = io.ReadAll(fileBase)
		if err != nil {
			return nil, err
		}
	}

	if !CheckFile(fileData) {
		log.Fatalln("纯真 IP 库存在错误，请重新下载")
	}

	header := fileData[0:8]
	start := binary.LittleEndian.Uint32(header[:4])
	end := binary.LittleEndian.Uint32(header[4:])

	return &QQwry{
		IPDB: wry.IPDB[uint32]{
			Data: fileData,

			OffLen:   3,
			IPLen:    4,
			IPCnt:    (end-start)/7 + 1,
			IdxStart: start,
			IdxEnd:   end,
		},
	}, nil
}

func (db QQwry) Find(query string, params ...string) (result fmt.Stringer, err error) {
	ip := net.ParseIP(query)
	if ip == nil {
		return nil, errors.New("query should be IPv4")
	}
	ip4 := ip.To4()
	if ip4 == nil {
		return nil, errors.New("query should be IPv4")
	}
	ip4uint := binary.BigEndian.Uint32(ip4)

	offset := db.SearchIndexV4(ip4uint)
	if offset <= 0 {
		return nil, errors.New("query not valid")
	}

	reader := wry.NewReader(db.Data)
	reader.Parse(offset + 4)
	return reader.Result.DecodeGBK(), nil
}

func (db QQwry) Name() string {
	return "qqwry"
}

func CheckFile(data []byte) bool {
	if len(data) < 8 {
		return false
	}

	header := data[0:8]
	start := binary.LittleEndian.Uint32(header[:4])
	end := binary.LittleEndian.Uint32(header[4:])

	if start >= end || uint32(len(data)) < end+7 {
		return false
	}

	return true
}
