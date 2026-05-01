package main

import (
	"fmt"
	"log"
	"net/http"

	"blog-mcp/internal/config"
	blogmcp "blog-mcp/internal/mcp"

	gomcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	server := gomcp.NewServer(&gomcp.Implementation{
		Name:    config.MustGetServerName(),
		Version: config.MustGetServerVersion(),
	}, nil)
	runner := blogmcp.NewRunner(config.MustGetRepoRoot())

	blogmcp.RegisterTools(server, runner)

	listenAddr := fmt.Sprintf(":%s", config.MustGetPort())

	if err := runServer(server, listenAddr); err != nil {
		log.Fatalf("サーバーの起動に失敗しました: %v", err)
	}
}

func runServer(server *gomcp.Server, listenAddr string) error {
	mux := http.NewServeMux()
	mux.Handle("/", gomcp.NewStreamableHTTPHandler(func(*http.Request) *gomcp.Server {
		return server
	}, nil))

	httpServer := &http.Server{
		Addr:    listenAddr,
		Handler: mux,
	}

	log.Printf("MCP サーバーを http://%s で待ち受けます", listenAddr)
	return httpServer.ListenAndServe()
}
