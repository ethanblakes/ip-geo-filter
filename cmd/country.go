package cmd

import (
	"fmt"
	"log"
	"net"

	"github.com/oschwald/geoip2-golang"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(countryCmd)
	countryCmd.Flags().StringVarP(&lang, "lang", "l", "zh-CN", "指定国家名称的语言代码 (例如: zh-CN, en, ja)")
}

var (
	lang string
)

var countryCmd = &cobra.Command{
	Use:   "country [ip_address...]",
	Short: "根据IP地址查询国家名称",
	Long:  "该命令接受一个或多个IP地址作为参数，并返回对应的国家名称。\n示例: igf country -l en ip1 ip2 ip3",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var invalidIps []string

		for _, ip := range args {
			if !IsValidIP(ip) {
				invalidIps = append(invalidIps, ip)
				continue // 处理下一个IP地址
			}

			db, err := geoip2.Open("resources/GeoLite2-City.mmdb")
			if err != nil {
				log.Panic(err)
			}
			defer db.Close()

			validIp := net.ParseIP(ip)
			record, err := db.City(validIp)
			if err != nil {
				fmt.Printf("查询IP '%s' 信息时出错: %v\n", ip, err) // 打印更详细的错误信息
				continue                                     // 处理下一个IP地址
			}

			countryName, ok := record.Country.Names[lang]
			if !ok {
				fmt.Printf("IP地址 '%s' 没有找到%s语言的翻译\n", ip, lang)
				continue
			}
			fmt.Printf("IP地址 '%s' 对应的国家名称(%s)是: %s\n", ip, lang, countryName)
		}

		if len(invalidIps) > 0 {
			fmt.Println("以下IP不合法(最多只展示前十条不合法IP):")
			for i := 0; i < 10 && i < len(invalidIps); i++ {
				fmt.Println(invalidIps[i])
			}
		}
	},
}
