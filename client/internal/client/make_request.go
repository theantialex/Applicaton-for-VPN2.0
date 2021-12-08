package client

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"math/rand"
	"net"
	"sync"

	"vpn2.0/app/client/internal/config"
	"vpn2.0/app/client/internal/tun"
	commands "vpn2.0/app/lib/cmd"
	"vpn2.0/app/lib/ctxmeta"
	"vpn2.0/app/lib/localnet"
)

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

func (c *Manager) MakeConnectRequest(ctx context.Context, name string, pass string, errCh chan error) (string, error) {
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


	go c.processConnectResponse(ctx, conn, resp, errCh)

	return resp, nil
}

func (c *Manager) processConnectResponse(ctx context.Context, conn net.Conn, respStr string, errCh chan error) {
	logger := ctxmeta.GetLogger(ctx)

	resp := commands.GetWords(respStr)
	if resp[0] != commands.SuccessResponse {
		logger.Error("got error in server resp")
		errCh <- errors.New(resp[0])
		return
	}

	id, err := localnet.GetIDFromIp(ctx, resp[1])
	if err != nil {
		errCh <- err
		return
	}

	c.SetClientID(id)

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

	var wg sync.WaitGroup
	wg.Add(1)
	go c.HandleTunEvent(ctx, tunIf, &wg, conn, errCh)
	wg.Add(1)
	go c.HandleConnTunEvent(ctx, tunIf, &wg, conn, errCh)
	wg.Wait()
}

func (c *Manager) MakeLeaveRequest(ctx context.Context, name string, pass string) (string, error) {
	logger := ctxmeta.GetLogger(ctx)

	conn, err := net.Dial("tcp", config.ADDR+":"+c.Config.ServerPort)
	if err != nil {
		logger.Error("failed to connect to server", zap.Error(err))
		return "", err
	}

	msg := fmt.Sprintf("%s %s %s %d", commands.LeaveCmd, name, pass, c.ID)

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

	respStrings := commands.GetWords(resp)
	if respStrings[0] != commands.SuccessResponse {
		logger.Error("got error in server resp")
		return "", nil
	}

	tunName := tun.GetTunName("client", 1, c.ID)
	c.SetClientID(IDUndefined)

	err = tun.SetTunDown(ctx, tunName)
	if err != nil {
		return "", err
	}

	return resp, nil
}

func (c *Manager) MakeDeleteRequest(ctx context.Context, name string, pass string) (string, error) {
	logger := ctxmeta.GetLogger(ctx)

	conn, err := net.Dial("tcp", config.ADDR+":"+c.Config.ServerPort)
	if err != nil {
		logger.Error("failed to connect to server", zap.Error(err))
		return "", err
	}

	msg := fmt.Sprintf("%s %s %s", commands.LeaveCmd, name, pass)

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

	respStrings := commands.GetWords(resp)
	if respStrings[0] != commands.SuccessResponse {
		logger.Error("got error in server resp")
		return "", nil
	}

	tunName := tun.GetTunName("client", 1, c.ID)
	c.SetClientID(IDUndefined)

	err = tun.SetTunDown(ctx, tunName)
	if err != nil {
		return "", err
	}

	return resp, nil
}
