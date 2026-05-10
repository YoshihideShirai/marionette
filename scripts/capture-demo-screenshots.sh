#!/usr/bin/env bash
set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
lib_root="${MARIONETTE_BROWSER_LIB_ROOT:-/tmp/marionette-browser-libs/root}"

ensure_browser_libs() {
  mkdir -p "$(dirname "$lib_root")/debs" "$lib_root"
  if [[ -f "$lib_root/usr/lib/x86_64-linux-gnu/libnss3.so" && -f "$lib_root/usr/lib/x86_64-linux-gnu/libasound.so.2" ]]; then
    return
  fi

  (
    cd "$(dirname "$lib_root")/debs"
    apt-get download libnss3 libnspr4 libasound2t64 >/dev/null
    for deb in ./*.deb; do
      dpkg-deb -x "$deb" "$lib_root"
    done
  )
}

ensure_browser_libs

tmpdir="$(mktemp -d)"
trap 'rm -rf "$tmpdir"' EXIT

cat >"$tmpdir/go.mod" <<'EOF'
module marionette-demo-screenshots

go 1.25

require github.com/go-rod/rod v0.116.2
EOF

cat >"$tmpdir/main.go" <<'EOF'
package main

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

type demo struct {
	name    string
	command []string
	url     string
	output  string
	height  int
	prepare func(*rod.Page)
}

func main() {
	repo := os.Getenv("MARIONETTE_REPO_ROOT")
	if repo == "" {
		panic("MARIONETTE_REPO_ROOT is required")
	}

	browserPath := launcher.NewBrowser().MustGet()
	browserURL := launcher.New().
		Bin(browserPath).
		Headless(true).
		Set("no-sandbox").
		Set("disable-dev-shm-usage").
		MustLaunch()
	browser := rod.New().ControlURL(browserURL).MustConnect()
	defer browser.MustClose()

	demos := []demo{
		{
			name:    "full Marionette demo",
			command: []string{"go", "run", "./cmd/marionette"},
			url:     "http://127.0.0.1:8080",
			output:  "docs/assets/marionette-demo.png",
			height:  1000,
		},
		{
			name:    "minimal sample",
			command: []string{"go", "run", "./cmd/simple-sample"},
			url:     "http://127.0.0.1:8081",
			output:  "docs/assets/simple-sample.png",
			height:  920,
			prepare: func(page *rod.Page) {
				page.MustElement("#task-name").MustInput("Prepare demo gallery")
				page.MustElementR("button", "Add Task").MustClick()
				time.Sleep(250 * time.Millisecond)
			},
		},
		{
			name:    "admin sample",
			command: []string{"go", "run", "./cmd/admin-sample"},
			url:     "http://127.0.0.1:8082",
			output:  "docs/assets/admin-sample.png",
			height:  920,
			prepare: func(page *rod.Page) {
				page.MustElementR("button", "Continue with Demo SSO").MustClick()
				time.Sleep(1200 * time.Millisecond)
			},
		},
		{
			name:    "DashWind demo",
			command: []string{"go", "run", "./cmd/dashwind-demo"},
			url:     "http://127.0.0.1:8083",
			output:  "docs/assets/dashwind-demo.png",
			height:  1000,
		},
		{
			name:    "AI chat sample",
			command: []string{"go", "run", "./cmd/ai-chat-sample"},
			url:     "http://127.0.0.1:8084",
			output:  "docs/assets/ai-chat-sample.png",
			height:  1000,
			prepare: func(page *rod.Page) {
				page.MustElement("#chat-prompt").MustInput("Show me the Marionette demo flow.")
				page.MustElementR("button", "Send").MustClick()
				time.Sleep(250 * time.Millisecond)
			},
		},
		{
			name:    "custom JavaScript sample",
			command: []string{"go", "run", "./cmd/custom-javascript-sample"},
			url:     "http://127.0.0.1:8082",
			output:  "docs/assets/custom-javascript-sample.png",
			height:  920,
			prepare: func(page *rod.Page) {
				time.Sleep(2 * time.Second)
			},
		},
	}

	for _, d := range demos {
		fmt.Printf("capturing %s\n", d.name)
		capture(repo, browser, d)
	}
}

func capture(repo string, browser *rod.Browser, d demo) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var stderr bytes.Buffer
	cmd := exec.CommandContext(ctx, d.command[0], d.command[1:]...)
	cmd.Dir = repo
	cmd.Stderr = &stderr
	cmd.Stdout = &stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	if err := cmd.Start(); err != nil {
		panic(fmt.Sprintf("%s: start failed: %v", d.name, err))
	}
	defer func() {
		cancel()
		if cmd.Process != nil {
			_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGTERM)
		}
		_ = cmd.Wait()
	}()

	if err := waitForHTTP(d.url, 30*time.Second); err != nil {
		panic(fmt.Sprintf("%s: server did not become ready: %v\n%s", d.name, err, stderr.String()))
	}

	page := browser.MustPage()
	defer page.MustClose()
	page.MustSetViewport(1440, d.height, 1, false)
	page.MustNavigate(d.url)
	page.MustWaitLoad()
	time.Sleep(600 * time.Millisecond)
	if d.prepare != nil {
		d.prepare(page)
	}

	output := filepath.Join(repo, d.output)
	if err := os.MkdirAll(filepath.Dir(output), 0o755); err != nil {
		panic(err)
	}
	page.MustScreenshot(output)
}

func waitForHTTP(url string, timeout time.Duration) error {
	client := &http.Client{Timeout: 700 * time.Millisecond}
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		resp, err := client.Get(url)
		if err == nil {
			_ = resp.Body.Close()
			if resp.StatusCode >= 200 && resp.StatusCode < 500 {
				return nil
			}
		}
		time.Sleep(250 * time.Millisecond)
	}
	return fmt.Errorf("timeout waiting for %s", url)
}
EOF

(
  cd "$tmpdir"
  go mod tidy
  MARIONETTE_REPO_ROOT="$repo_root" \
  LD_LIBRARY_PATH="$lib_root/usr/lib/x86_64-linux-gnu${LD_LIBRARY_PATH:+:$LD_LIBRARY_PATH}" \
  go run .
)
