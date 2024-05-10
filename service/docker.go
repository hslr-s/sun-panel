package main

import (
	"context"
	"fmt"
	"sun-panel/lib/docker"

	"github.com/docker/docker/api/types"
)

func main() {
	if client, err := docker.NewDockerClient(); err != nil {
		fmt.Println("docker连接错误", err)
	} else {
		options := types.ImageListOptions{}
		if list, err := client.ImageList(context.Background(), options); err != nil {
			fmt.Println("docker查询错误", err)
		} else {
			for _, v := range list {
				fmt.Println(v.ID, v.Labels, v.Containers, v.Size)
			}
		}
	}

}
