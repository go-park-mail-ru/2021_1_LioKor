package repository

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"io"
	"liokor_mail/internal/pkg/common"
	"os"
	"testing"
	"time"
)

var (
	repo common.PostgresDataBase
)

var containerConfig = common.Config{
	DBHost: "127.0.0.1",
	DBPort: 15000,
	DBUser: "liokor",
	DBPassword: "Qwerty123",
	DBDatabase: "liokor_mail",
	DBConnectTimeout: 10,
}

func TestMain(m *testing.M) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	reader, err := cli.ImagePull(ctx, "docker.io/library/postgres", types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, reader)

	//bind := "C:/Users/Altana/GolandProjects/2021_1_LioKor/internal/pkg/mail/repository/migrations/:/docker-entrypoint-initdb.d/"
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: "postgres",
		Env: []string{
			"POSTGRES_DB=liokor_mail",
			"POSTGRES_USER=liokor",
			"POSTGRES_PASSWORD=Qwerty123",
		},
	}, &container.HostConfig{
		//Binds: []string {bind},
		PortBindings: nat.PortMap{"5432/tcp": []nat.PortBinding{{"127.0.0.1", "15000"}}},
	}, nil, nil, "")
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	defer func() {
		if err = cli.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions{Force: true}); err != nil {
			panic(err)
		}
	}()

	time.Sleep(time.Second * 10)

	repo, err = common.NewPostgresDataBase(containerConfig)
	if err != nil {
		panic(err)
	}

	_, err = repo.DBConn.Exec(
		ctx,
		"CREATE EXTENSION IF NOT EXISTS CITEXT;" +
			"CREATE TABLE IF NOT EXISTS users (" +
			"id SERIAL PRIMARY KEY," +
			"username CITEXT UNIQUE NOT NULL," +
			"password_hash CITEXT NOT NULL," +
			"avatar_url VARCHAR(128)," +
			"fullname VARCHAR(128)," +
			"reserve_email CITEXT" +
			");",
		)
	if err != nil {
		panic(err)
	}
}