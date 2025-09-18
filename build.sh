
binary_name=smartcut_windows_amd64.exe
go build -ldflags -H=windowsgui cmd/smartcut/main.go
mv main.exe smartcut_windows_amd64.exe
rm binaries.zip
zip -r binaries_windows.zip smartcut_windows_amd64.exe