package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"sync/atomic"
	"syscall"
	"time"
)

var (
	p = atomic.Value{}
)

func getEnv(name, defaultValue string) string {
	if v := os.Getenv(name); v != "" {
		return v
	}
	return defaultValue
}

func Password() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if !strings.EqualFold(req.Method, http.MethodPost) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer req.Body.Close()
		var pass = struct {
			P string `json:"p"`
		}{}
		if err := json.NewDecoder(req.Body).Decode(&pass); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if pass.P == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		p.Store(pass.P)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}
}

func runHelixRelayer(ctx context.Context) *exec.Cmd {
	subProcess := exec.CommandContext(ctx, "node", "./dist/src/main")
	subProcess.Dir = getEnv("HELIX_RELAYER_DIR", "./relayer")
	_ = os.Setenv("HELIX_RELAYER_PASSWORD", p.Load().(string))
	if os.Getenv("LP_BRIDGE_PATH") == "" {
		_ = os.Setenv("LP_BRIDGE_PATH", "./.maintain/configure.json")
	}

	if os.Getenv("LP_BRIDGE_STORE_PATH") == "" {
		_ = os.Setenv("LP_BRIDGE_STORE_PATH", "./.maintain/db")
	}

	subProcess.Env = os.Environ()

	subProcess.Stdout = os.Stdout
	subProcess.Stderr = os.Stderr
	subProcess.Stdin = os.Stdin
	if err := subProcess.Start(); err != nil {
		Panicf("cmd.Run() failed with %s\n", err)
	}
	go func() {
		if err := subProcess.Wait(); err != nil {
			Panicf("cmd.Run() failed with %s\n", err)
		}
	}()
	return subProcess
}

func healthCheck(ctx context.Context, cmd *exec.Cmd) {
	var (
		t = time.NewTicker(time.Minute)
	)
	defer t.Stop()
	defer func() {
		if err := recover(); err != nil {
			Warnf("recovered from panic: %s\n", err)
		}
		go healthCheck(ctx, runHelixRelayer(ctx))
	}()
	serveHealthCheck := func() error {
		client := &http.Client{
			Timeout: time.Second * 5,
		}
		resp, err := client.Get("http://127.0.0.1:3000/")
		if err != nil {
			return fmt.Errorf("http.Get() failed with %s\n", err)
		}
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("http.Get() failed with %s\n", resp.Status)
		}
		defer resp.Body.Close()
		return nil
	}
	var serveHealthCheckErrorCount = 0
	for {
		select {
		case <-t.C:
			if err := serveHealthCheck(); err != nil {
				serveHealthCheckErrorCount++
				if serveHealthCheckErrorCount > 3 {
					return
				}
			}
			if cmd.ProcessState != nil && cmd.ProcessState.Exited() {
				Warnf("helix relayer exited with %s, restarting...\n", cmd.ProcessState.String())
				cmd = runHelixRelayer(ctx)
				continue
			}
			serveHealthCheckErrorCount = 0
			Debugf("helix relayer is running\n")
		case <-ctx.Done():
			return
		}
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.TODO())

	server := &http.Server{Addr: getEnv("SERVER_ADDR", ":8080"), Handler: Password()}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			Panicf("http server error: %s\n", err)
		}
	}()
	for {
		if p.Load() == nil || p.Load().(string) == "" {
			Debugf("waiting for password...\n")
			time.Sleep(time.Second)
			continue
		}
		break
	}
	go healthCheck(ctx, runHelixRelayer(ctx))

	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	_ = server.Shutdown(context.TODO())
	cancel()
	time.Sleep(time.Second * 3)
}
