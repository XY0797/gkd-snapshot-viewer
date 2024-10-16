package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
)

type Config struct {
	Port string
	Root string
}

func readConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("无法打开配置文件: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	config := &Config{}

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, ";") || strings.HasPrefix(line, "#") || line == "" {
			continue // skip comments and empty lines
		}

		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			// section handling can be added here if needed
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) < 2 {
			continue // ignore lines without an "="
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "port":
			config.Port = value
		case "root":
			config.Root = value
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("无法读取配置文件: %w", err)
	}

	return config, nil
}

var indexHTMLPath string

func snapshotHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, indexHTMLPath)
}

func main() {
	config, err := readConfig("config.ini")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	indexHTMLPath = path.Join(config.Root, "index.html")

	http.Handle("/", http.FileServer(http.Dir(config.Root)))
	http.HandleFunc("/snapshot/", snapshotHandler) // 添加新的路由处理器

	log.Printf("启动服务器在端口 %s, 网站文件来源目录为 '%s'", config.Port, config.Root)
	log.Printf("访问 http://127.0.0.1:%s 浏览网页", config.Port)

	// 在这里尝试打开浏览器
	url := fmt.Sprintf("http://127.0.0.1:%s/", config.Port)
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}
	err = cmd.Start()
	if err != nil {
		log.Printf("未能打开浏览器: %v", err)
	}

	err = http.ListenAndServe(":"+config.Port, nil)
	if err != nil {
		log.Fatalf("启动服务器失败: %v", err)
	}
}
