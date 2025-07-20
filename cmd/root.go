package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"

	"github.com/abc1763613206/nabili/internal/constant"
	"github.com/abc1763613206/nabili/internal/db"
	"github.com/abc1763613206/nabili/pkg/common"
	"github.com/abc1763613206/nabili/pkg/entity"
)

var rootCmd = &cobra.Command{
	Use:   "nabili",
	Short: "An offline tool for querying IP geographic information",
	Long: `An offline tool for querying IP geographic information.

Find document on: https://github.com/abc1763613206/nabili

#1 Query a simple IP address

	$ nabili 1.2.3.4

  or use pipe

	$ echo IP 6.6.6.6 | nabili

#2 Query multiple IP addresses

	$ nabili 1.2.3.4 4.3.2.1 123.23.3.0

#3 Interactive query

	$ nabili
	123.23.23.23
	123.23.23.23 [越南 越南邮电集团公司]
	quit

#4 Use with dig

	$ dig nabili.zu1k.com +short | nabili

#5 Use with nslookup

	$ nslookup nabili.zu1k.com 8.8.8.8 | nabili

#6 Use with any other program

	bash abc.sh | nabili

#7 IPV6 support

#8 Specify database provider

	$ nabili -4 geoip 1.2.3.4
	$ nabili -6 geoip 2001:db8::1
	$ nabili --db4 ip2region 8.8.8.8
	$ nabili --db6 zxipv6wry 2001:db8::1
	$ nabili -4 bili 8.8.8.8
	$ nabili -6 bili 240e:b1:a810:2011::a1
	$ nabili -4 ipsb 8.8.8.8
	$ nabili -6 ipsb 240e:b1:a810:2011::a1
	$ nabili -4 iqiyi 8.8.8.8
	$ nabili -6 iqiyi 240e:b1:a810:2011::a1
	$ nabili -4 baidu 8.8.8.8
	$ nabili -6 baidu 240e:b1:a810:2011::a1
`,
	Version: constant.Version,
	Args:    cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		gbk, _ := cmd.Flags().GetBool("gbk")
		isJson, _ := cmd.Flags().GetBool("json")
		ipv4DB, _ := cmd.Flags().GetString("db4")
		ipv6DB, _ := cmd.Flags().GetString("db6")

		// Set command-line database selections
		if ipv4DB != "" {
			db.CmdIPv4DB = ipv4DB
		}
		if ipv6DB != "" {
			db.CmdIPv6DB = ipv6DB
		}

		if len(args) == 0 {
			stdin := bufio.NewScanner(os.Stdin)
			stdin.Split(common.ScanLines)
			for stdin.Scan() {
				line := stdin.Text()
				if gbk {
					line, _, _ = transform.String(simplifiedchinese.GBK.NewDecoder(), line)
				}
				if line := strings.TrimSpace(line); line == "quit" || line == "exit" {
					return
				}
				if isJson {
					_, _ = fmt.Fprintf(color.Output, "%s", entity.ParseLine(line).Json())
				} else {
					_, _ = fmt.Fprintf(color.Output, "%s", entity.ParseLine(line).ColorString())
				}
			}
		} else {
			if isJson {
				_, _ = fmt.Fprintf(color.Output, "%s", entity.ParseLine(strings.Join(args, " ")).Json())
			} else {
				for _, line := range args {
					_, _ = fmt.Fprintf(color.Output, "%s\n", entity.ParseLine(line).ColorString())
				}
			}
		}
	},
}

// Execute parse subcommand and run
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err.Error())
	}
}

func init() {
	rootCmd.Flags().Bool("gbk", false, "Use GBK decoder")
	rootCmd.Flags().BoolP("json", "j", false, "Output in JSON format")
	rootCmd.Flags().StringP("db4", "4", "", "IPv4 database provider (qqwry, geoip, ip2region, dbip, ipip, ip2location, bili, ipsb, iqiyi, baidu)")
	rootCmd.Flags().StringP("db6", "6", "", "IPv6 database provider (zxipv6wry, geoip, dbip, ipip, ip2location, bili, ipsb, iqiyi, baidu)"))
}
