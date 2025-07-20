package main

import (
	"github.com/abc1763613206/nabili/internal/constant"

	"github.com/abc1763613206/nabili/cmd"
	"github.com/abc1763613206/nabili/internal/config"

	_ "github.com/abc1763613206/nabili/internal/migration"
)

func main() {
	config.ReadConfig(constant.ConfigDirPath)
	cmd.Execute()
}
