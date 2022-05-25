package cloudflare

import (
	"context"
	"digtal/pkg/digtal-domain/config"
	"fmt"
	"github.com/cloudflare/cloudflare-go"
	"log"
	"math/rand"
	"os"
	"time"
	"unsafe"
)

func ShowAllDns(domain string) error {
	api, err := cloudflare.New(config.C.CLApiKey, config.C.CLEmail)
	if err != nil {
		log.Fatal(err)
	}
	zoneID, err := api.ZoneIDByName(domain)
	if err != nil {
		fmt.Println(err)
		return err
	}

	records, err := api.DNSRecords(context.Background(), zoneID, cloudflare.DNSRecord{})
	if err != nil {
		return err
	}
	for _, record := range records {
		api.DeleteDNSRecord(context.Background(), zoneID, record.ID)
	}
	return nil
}

func CreateDns(domain string, ip string) (string, error) {

	api, err := cloudflare.New(config.C.CLApiKey, config.C.CLEmail)
	if err != nil {
		log.Fatal(err)
	}
	zoneID, err := api.ZoneIDByName(domain)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	isProxied := false
	recordName := randStr(6)
	record := cloudflare.DNSRecord{
		Name:    recordName,
		Type:    "A",
		Content: ip,
		TTL:     60,
		Proxied: &isProxied,
	}

	_, err = api.CreateDNSRecord(context.Background(), zoneID, record)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating DNS record: ", err)
		return "", err
	}

	return recordName + "." + domain, nil
}

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var src = rand.NewSource(time.Now().UnixNano())

const (
	// 6 bits to represent a letter index
	letterIdBits = 6
	// All 1-bits as many as letterIdBits
	letterIdMask = 1<<letterIdBits - 1
	letterIdMax  = 63 / letterIdBits
)

func randStr(n int) string {
	b := make([]byte, n)
	// A rand.Int63() generates 63 random bits, enough for letterIdMax letters!
	for i, cache, remain := n-1, src.Int63(), letterIdMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdMax
		}
		if idx := int(cache & letterIdMask); idx < len(letters) {
			b[i] = letters[idx]
			i--
		}
		cache >>= letterIdBits
		remain--
	}
	return *(*string)(unsafe.Pointer(&b))
}
