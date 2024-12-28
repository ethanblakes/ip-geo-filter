package cmd

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"github.com/oschwald/geoip2-golang"
	"github.com/spf13/cobra"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
)

func init() {
	rootCmd.AddCommand(countryCmd)
	countryCmd.Flags().StringVarP(&lang, "lang", "l", "zh-CN", "指定国家名称的语言代码 (例如: zh-CN, en, ja)")
	countryCmd.Flags().StringVarP(&url, "url", "u", "", "指定批量转换使用的文件路径")
	countryCmd.Flags().StringVarP(&skip, "skip", "s", "true", "是否跳过异常和空值(true:跳过,false:保留)")
	countryCmd.Flags().StringVarP(&out, "output", "o", "", "指定输出为csv文件的路径(尽量避免中文路径已减少未知错误的发生)")
}

var (
	lang       string //记录指定语言 默认为中文
	url        string //记录文件路径 默认为空
	skip       string
	invalidIps []string
	results    map[string]string
	out        string
	sum        int
	vaildIps   []string
)

var countryCmd = &cobra.Command{
	Use:   "country [ip_address...]",
	Short: "根据IP地址查询国家名称",
	Long:  "该命令接受一个或多个IP地址作为参数，并返回对应的国家名称。\n示例: igf country -l en ip1 ip2 ip3",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		if url != "" {
			log.Printf("---读取路径[%s]下的文件开始---", url)
			//读取对应路径的文件中的ip作为开始,每次成功读取到ip则sum加一
			//判断路径是否合法
			if isPathValid(url) == false {
				return
			}

			//逐行读取文件内容
			file, file_open_err := os.Open(url)

			if file_open_err != nil {
				fmt.Println("文件读取异常")
				return
			}

			defer file.Close()

			log.Println("文件打开成功，开始逐行读取内容...")

			// 创建逐行读取器
			scanner := bufio.NewScanner(file)
			sum := 0 // 记录成功读取的行数

			for scanner.Scan() {
				line := scanner.Text()         // 获取当前行内容
				line = strings.TrimSpace(line) // 去除前后空格

				if line == "" {
					continue // 跳过空行
				}

				vaildIps = append(vaildIps, line) // 假设每行是一个有效 IP
				sum++                             // 成功读取 IP，计数加一
			}

			// 检查读取是否有错误
			if err := scanner.Err(); err != nil {
				log.Println("读取文件时发生错误:", err)
			} else {
				log.Printf("文件读取完成，共读取到 %d 个 IP。", sum)
			}

		}
		if len(args) > 0 {
			//读取参数中的ip
			for _, ip := range args {
				//验证ip合法性
				if IsValidIP(ip) == false {
					log.Printf("非法的IP:[%s]", ip)
					continue
				}
				vaildIps = append(vaildIps, ip)
			}
		}

		results = make(map[string]string)

		//调用函数对有效ip进行处理
		for _, ip := range vaildIps {
			convert(ip)
		}

		//将csv文件输出的指定目录
		if out != "" {
			//首先判断路径下是否有文件存在
			if strings.HasSuffix(out, "\\") {
				out = out + "data.csv" // 如果已经有 "/", 直接拼接文件名
			} else if strings.HasSuffix(out, ".csv") {
				//啥也不用干
			} else {
				out = out + "\\data.csv" // 如果没有 "/", 添加 "/" 再拼接文件名
			}
			nfile, file_open_error := os.Open(out)

			chose := ""

			if file_open_error == nil {
				//该目录下有文件存在，询问是否覆盖
				defer nfile.Close()
				log.Println("该路径下已经有文件存在，是否覆盖(y/n)")
				fmt.Scan(&chose)

				if chose == "n" || chose == "no" {
					return
				}
			}

			//创建csv文件
			file, file_create_error := os.Create(out)

			if file_create_error != nil {
				log.Printf("创建文件失败，{error:%v}", file_create_error)
				return
			}
			defer func(file *os.File) {
				err := file.Close()
				if err != nil {
					fmt.Println(err)
				}
			}(file)

			// 写入 BOM 头，确保 UTF-8 文件可被正确识别
			file.WriteString("\uFEFF") // 写入 BOM

			writer := csv.NewWriter(file)
			defer writer.Flush()
			//写入csv的列头
			err := writer.Write([]string{"ip", "country"})
			if err != nil {
				fmt.Println(err)
				return
			}

			//获取map中的键值对
			for ip, country := range results {
				err := writer.Write([]string{ip, country})
				if err != nil {
					fmt.Println(err)
					return
				}
			}

		}

	},
}

func convert(ip string) {

	if !IsValidIP(ip) {
		if skip == "false" {
			results[ip] = ""
		}
		log.Printf("非法的ip:[%s]", ip)
		invalidIps = append(invalidIps, ip)
		return // 处理下一个IP地址
	}

	db, err := geoip2.Open("resources/GeoLite2-City.mmdb")
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	validIp := net.ParseIP(ip)
	record, err := db.City(validIp)
	if err != nil {
		if skip == "false" {
			results[ip] = ""
		}
		log.Printf("异常-->%v\n", err)
		return // 处理下一个IP地址
	}

	countryName, ok := record.Country.Names[lang]
	if !ok {
		if skip == "false" {
			results[ip] = ""
			log.Printf("未获取到%s对应的国家名称\n", ip)
		}
		return
	}

	//成功获取
	results[ip] = countryName
	log.Printf("IP地址 '%s' 对应的国家名称(%s)是: %s\n", ip, lang, countryName)
}

func isPathValid(path string) bool {
	// 获取绝对路径，确保路径格式正确
	absolutePath, err := filepath.Abs(path)
	if err != nil {
		log.Println("无法解析的路径")
		return false // 无法解析路径
	}

	// 检查路径是否存在（可以注释掉此行如果只需要验证格式）
	if _, err := os.Stat(absolutePath); os.IsNotExist(err) {
		log.Println("路径下文件不存在")
		return false // 文件或目录不存在
	}

	// 返回合法性
	return true
}
