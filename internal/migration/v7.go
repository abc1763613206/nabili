package migration

import (
	"log"
	"strings"

	"github.com/spf13/viper"
	"github.com/abc1763613206/nabili/internal/constant"
	"github.com/abc1763613206/nabili/internal/db"
	"github.com/abc1763613206/nabili/pkg/qqwry"
)

func migration2v7() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(constant.ConfigDirPath)

	err := viper.ReadInConfig()
	if err != nil {
		return
	}

	dbList := db.List{}
	err = viper.UnmarshalKey("databases", &dbList)
	if err != nil {
		log.Fatalln("Config invalid:", err)
	}

	needOverwrite := false
	for _, adb := range dbList {
		if adb.Name == "qqwry" {
			if len(adb.DownloadUrls) == 0 ||
				adb.DownloadUrls[0] == "https://99wry.cf/qqwry.dat" ||
				strings.Contains(adb.DownloadUrls[0], "sspanel-uim") {
				needOverwrite = true
				adb.DownloadUrls = qqwry.DownloadUrls
			}
		}
	}

	if needOverwrite {
		viper.Set("databases", dbList)
		err = viper.WriteConfig()
		if err != nil {
			log.Println(err)
		}
	}
}
