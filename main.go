package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/cloudflare/cloudflare-go"
)

var (
	zoneName  = os.Getenv("ZONE_NAME")
	subDomain = os.Getenv("SUB_DOMAIN")
	apiKey    = os.Getenv("CF_API_KEY")
	apiEmail  = os.Getenv("CF_API_EMAIL")
)

// GetPublicIP
// https://ip.seeip.org
func GetPublicIP() (string, error) {
	resp, err := http.Get("https://ip.seeip.org/json")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var body map[string]string

	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		return "", err
	}
	ip, ok := body["ip"]
	if ok {
		return ip, nil
	}
	return "", errors.New("not found")
}

func job(api *cloudflare.API, id string, record cloudflare.DNSRecord) {
	// 获取外网IP地址
	ip, err := GetPublicIP()
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("ip %s \n", ip)

	// 查询DNS记录是否存在
	records, err := api.DNSRecords(id, record)
	if err != nil {
		log.Println(err)
		return
	}

	rr := records[0]

	log.Printf("Name %s ,Content %s\n", rr.Name, rr.Content)

	// 不存在插入
	if len(records) == 0 {
		res, err := api.CreateDNSRecord(id, record)
		if err != nil {
			log.Println(err)
		}
		fmt.Println(res)
		return
	}

	if rr.Content != ip {
		// 加一个IP是否更改的判断
		record.Content = ip
		err := api.UpdateDNSRecord(id, rr.ID, record)
		if err != nil {
			log.Println(err)
		}
		log.Printf("UpdateDNSRecord success \n")
	}
}

func main() {

	// https://api.cloudflare.com/
	// https://pkg.go.dev/github.com/cloudflare/cloudflare-go
	// https://github.com/cloudflare/cloudflare-go
	// 获取cloudflare存储结果
	api, err := cloudflare.New(apiKey, apiEmail)
	if err != nil {
		log.Fatal(err)
	}

	// Fetch the zone ID
	id, err := api.ZoneIDByName(zoneName) // Assuming example.com exists in your Cloudflare account already
	if err != nil {
		log.Fatal(err)
	}

	record := cloudflare.DNSRecord{
		Type: "A",
		Name: subDomain,
		TTL:  120,
	}

	ticker := time.NewTicker(15 * time.Second)

	for {
		select {
		case <-ticker.C:
			log.Println("run job")
			go job(api, id, record)
		}
	}
}
