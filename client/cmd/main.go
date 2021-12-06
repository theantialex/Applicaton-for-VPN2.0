package client

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"vpn2.0/app/client/config"
	"vpn2.0/app/client/logs"
	"vpn2.0/app/lib/localnet"

	"go.uber.org/zap"

	commands "vpn2.0/app/lib/cmd"
	"vpn2.0/app/lib/ctxmeta"
	"vpn2.0/app/lib/tun"
)

func SetUpClient() (*Manager, context.Context) {
	conf, err := config.New()
	if err != nil {
		panic(err)
	}

	logger := logs.BuildLogger(conf)
	ctx := ctxzap.ToContext(context.Background(), logger)
	logger.Info("client starting...")

	c := Manager{Config: conf}
	return &c, ctx
}

func (c *Manager) MakeCreateRequest(ctx context.Context, name string, pass string) (string, error) {
	logger := ctxmeta.GetLogger(ctx)

	conn, err := net.Dial("tcp", config.ADDR+":"+c.Config.ServerPort)
	if err != nil {
		logger.Error("failed to connect to server", zap.Error(err))
		return "", err
	}

	msg := fmt.Sprintf("%s %s %s", commands.CreateCmd, name, pass)

	_, err = conn.Write([]byte(msg + "\n"))
	if err != nil {
		logger.Error("failed to write to conn", zap.Error(err))
		return "", err
	}

	clientReader := bufio.NewReader(conn)
	resp, err := clientReader.ReadString('\n')
	if err != nil {
		logger.Error("failed to read from conn", zap.Error(err))
		return "", err
	}
	return resp, nil
}

func (c *Manager) MakeConnectRequest(ctx context.Context, name string, pass string) (string, error) {
	logger := ctxmeta.GetLogger(ctx)

	conn, err := net.Dial("tcp", config.ADDR+":"+c.Config.ServerPort)
	if err != nil {
		logger.Error("failed to connect to server", zap.Error(err))
		return "", err
	}

	msg := fmt.Sprintf("%s %s %s", commands.ConnectCmd, name, pass)

	_, err = conn.Write([]byte(msg + "\n"))
	if err != nil {
		logger.Error("failed to write to conn", zap.Error(err))
		return "", err
	}

	clientReader := bufio.NewReader(conn)
	resp, err := clientReader.ReadString('\n')
	if err != nil {
		logger.Error("failed to read from conn", zap.Error(err))
		return "", err
	}

	go KeepConnection(ctx, conn, resp)
	return resp, nil
}

func KeepConnection(ctx context.Context, conn net.Conn, resp string) error {
	errCh := make(chan error, 1)
	go ProcessResp(ctx, conn, resp, errCh)
	close(errCh)
	if err := <-errCh; err != nil {
		return err
	}
	return nil
}

func ProcessResp(ctx context.Context, conn net.Conn, respStr string, errCh chan error) {
	logger := ctxmeta.GetLogger(ctx)

	resp := commands.GetWords(respStr)

	rand.Seed(time.Now().UnixNano())
	tunName := tun.GetTunName("client", 1, 10+rand.Intn(191))
	tunIf, err := tun.ConnectToTun(ctx, tunName)
	if err != nil {
		errCh <- err
	}
	logger.Debug("connected to tun", zap.String("tun_name", tunName))

	brd := localnet.GetBrdFromIp(ctx, resp[1])
	if brd == "" {
		errCh <- errors.New("failed to get brd")
	}

	err = tun.SetTunUp(ctx, resp[1], brd, tunName)
	if err != nil {
		errCh <- err
	}
	logger.Debug("set tun up", zap.String("tun_name", tunName))

	go tun.HandleTunEvent(ctx, tunIf, conn, errCh)
	go tun.HandleConnTunEvent(ctx, tunIf, conn, errCh)
}
