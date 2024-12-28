package cmd

import (
	"fmt"
	"github.com/oschwald/geoip2-golang"
	"github.com/spf13/cobra"
	"log"
	"net"
)

func init() {
	rootCmd.AddCommand(isoCmd)
}

var isoCmd = &cobra.Command{
	Use:   "iso [ip_address...]", // 支持多个IP地址
	Short: "根据IP地址查询ISO国家代码",
	Long:  "该命令接受一个或多个IP地址作为参数，并返回对应的ISO国家代码。",
	Args:  cobra.MinimumNArgs(1), // 最少接受一个参数
	Run: func(cmd *cobra.Command, args []string) {
		var invalidIps []string

		// 显示合法的IP的国家代码
		for _, ip := range args {
			if !IsValidIP(ip) {
				invalidIps = append(invalidIps, ip)
				continue // 如果IP无效，跳过
			}
			isoCode := GetISObyIP(ip)
			if isoCode == "" {
				// 如果ISO代码为空，输出提示
				fmt.Printf("IP地址 '%s' 查询不到ISO国家代码。\n", ip)
			} else {
				fmt.Printf("IP地址 '%s' 对应的ISO国家代码是: %s\n", ip, isoCode)
			}
		}

		// 列举不合法的IP
		if len(invalidIps) > 0 {
			fmt.Println("以下IP不合法(最多只展示前十条不合法IP):")
			for i := 0; i < 10 && i < len(invalidIps); i++ {
				fmt.Println(invalidIps[i])
			}
		}
	},
}

func GetISObyIP(ip string) string {
	db, err := geoip2.Open("GeoIP2-City.mmdb")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	parsedIp := net.ParseIP(ip)
	record, err := db.City(parsedIp)
	if err != nil {
		log.Fatal(err)
	}
	return record.Country.IsoCode
}

func IsValidIP(ip string) bool {
	// 可以根据需求做更严格的验证，例如判断IPv4或IPv6
	return net.ParseIP(ip) != nil
}
