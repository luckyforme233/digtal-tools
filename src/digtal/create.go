package digtal

import (
	"context"
	"digtal/pkg/digtal-domain/config"
	"fmt"
	"github.com/digitalocean/godo"
	"io/ioutil"
	"log"
)

func CreateDroplet() (*godo.Droplet, *godo.Response, error) {
	dropletName := "super-cool-droplet"
	keys := make([]godo.DropletCreateSSHKey, 0)
	key, err := CreateKeys()
	if err != nil {
		return nil, nil, err
	}
	keys = append(keys, *key)
	createRequest := &godo.DropletCreateRequest{
		Name:   dropletName,
		Region: "SFO3",
		Size:   "s-1vcpu-1gb",
		Image: godo.DropletCreateImage{
			Slug: "ubuntu-20-04-x64",
		},
		SSHKeys: keys,
	}
	ctx := context.TODO()
	return client.Droplets.Create(ctx, createRequest)
}

func DeleteDroplet(id int) {
	ctx := context.TODO()
	response, err := client.Droplets.Delete(ctx, id)
	if err != nil {
		log.Println("删除失败", err)
	}
	fmt.Println(response.String())
}

func GetDropLetInfo(id int) (*godo.Droplet, *godo.Response, error) {
	ctx := context.TODO()
	return client.Droplets.Get(ctx, id)
}

func CreateKeys() (*godo.DropletCreateSSHKey, error) {

	ctx := context.TODO()

	list, _, err := client.Keys.List(ctx, &godo.ListOptions{
		Page:    1,
		PerPage: 20,
	})
	if len(list) > 0 && err == nil {
		return &godo.DropletCreateSSHKey{
			Fingerprint: list[0].Fingerprint,
			ID:          list[0].ID,
		}, nil
	}

	key, err := ioutil.ReadFile(config.C.PubKeyPath)
	if err != nil {
		log.Fatalf("unable to read private key: %v", err)
		return nil, err
	}
	Key, _, err := client.Keys.Create(ctx, &godo.KeyCreateRequest{
		Name:      "localhost",
		PublicKey: string(key),
	})
	if err != nil {
		return nil, err
	}
	return &godo.DropletCreateSSHKey{
		ID:          Key.ID,
		Fingerprint: Key.Fingerprint,
	}, nil
}
