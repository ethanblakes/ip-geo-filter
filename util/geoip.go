package util

import (
	"fmt"
	"github.com/oschwald/geoip2-golang"
	"log"
	"net"
)

func GetISObyIP(ip string) string {
	db, err := geoip2.Open("resources/GeoLite2-City.mmdb")
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()
	// If you are using strings that may be invalid, check that ip is not nil
	validatedIP := net.ParseIP(ip)
	record, err := db.City(validatedIP)
	if err != nil {
		fmt.Println(err)
	}
	return record.Country.IsoCode
}
