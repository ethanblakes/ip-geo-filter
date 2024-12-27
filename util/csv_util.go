package util

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

var countryMap map[string]string

func init() {
	countryMap = make(map[string]string)
	err := readCountryCSV()
	if err != nil {
		fmt.Println("初始化国家代码映射表失败:", err)
		panic(err)
	}
}

func readCountryCSV() error {
	file, err := os.Open("resources/ISO-3166.csv")
	if err != nil {
		return fmt.Errorf("打开 ISO-3166.csv 文件失败: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	//设置允许字段中带引号
	reader.LazyQuotes = true
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("读取 CSV 记录失败: %w", err)
		}

		if len(record) < 2 { // 检查字段数量，至少要有两个字段
			fmt.Printf("CSV 记录格式错误，字段数量不足: %v\n", record)
			continue // 跳过此条记录，继续读取下一条
		}
		country := record[0]
		code := record[1]

		countryMap[code] = country
	}

	return nil
}

func GetCountry(code string) string {
	if countryMap == nil {
		return "国家代码映射表未初始化"
	}
	value, ok := countryMap[code]
	if ok {
		return value
	}
	return "未找到该国家代码"
}
