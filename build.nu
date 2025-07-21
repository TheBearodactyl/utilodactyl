CGO_ENABLED=0 go build -ldflags="-s -w -extldflags '-static'" -trimpath
upx --best --lzma utilodactyl.exe
