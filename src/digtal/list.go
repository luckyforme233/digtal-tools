package digtal

import (
	"context"
	"github.com/digitalocean/godo"
	"time"
)

type DropLetItem struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Memory int    `json:"memory"`
	Vcpus  int    `json:"vcpus"`
	Disk   int    `json:"disk"`
	Region struct {
		Slug      string   `json:"slug"`
		Name      string   `json:"name"`
		Sizes     []string `json:"sizes"`
		Available bool     `json:"available"`
		Features  []string `json:"features"`
	} `json:"region"`
	Image struct {
		ID            int       `json:"id"`
		Name          string    `json:"name"`
		Type          string    `json:"type"`
		Distribution  string    `json:"distribution"`
		Slug          string    `json:"slug"`
		Public        bool      `json:"public"`
		Regions       []string  `json:"regions"`
		MinDiskSize   int       `json:"min_disk_size"`
		SizeGigabytes float64   `json:"size_gigabytes"`
		CreatedAt     time.Time `json:"created_at"`
		Description   string    `json:"description"`
		Status        string    `json:"status"`
	} `json:"image"`
	Size struct {
		Slug         string   `json:"slug"`
		Memory       int      `json:"memory"`
		Vcpus        int      `json:"vcpus"`
		Disk         int      `json:"disk"`
		PriceMonthly int      `json:"price_monthly"`
		PriceHourly  float64  `json:"price_hourly"`
		Regions      []string `json:"regions"`
		Available    bool     `json:"available"`
		Transfer     int      `json:"transfer"`
		Description  string   `json:"description"`
	} `json:"size"`
	SizeSlug string   `json:"size_slug"`
	Features []string `json:"features"`
	Status   string   `json:"status"`
	Networks struct {
		V4 []struct {
			IPAddress string `json:"ip_address"`
			Netmask   string `json:"netmask"`
			Gateway   string `json:"gateway"`
			Type      string `json:"type"`
		} `json:"v4"`
	} `json:"networks"`
	CreatedAt time.Time     `json:"created_at"`
	VolumeIds []interface{} `json:"volume_ids"`
	VpcUUID   string        `json:"vpc_uuid"`
}

func DropletList() ([]godo.Droplet, error) {
	todo := context.TODO()
	// create a list to hold our droplets
	list := make([]godo.Droplet, 0)

	// create options. initially, these will be blank
	opt := &godo.ListOptions{}
	for {
		droplets, resp, err := client.Droplets.List(todo, opt)
		if err != nil {
			return nil, err
		}
		// append the current page's droplets to our list
		list = append(list, droplets...)
		// if we are at the last page, break out the for loop
		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}
		page, err := resp.Links.CurrentPage()
		if err != nil {
			return nil, err
		}
		// set the page we want for the next request
		opt.Page = page + 1
	}
	return list, nil
}
