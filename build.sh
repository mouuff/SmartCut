# go build -ldflags -H=windowsgui cmd/smartcuts/main.go

binary_name=smartcuts_windows_amd64.exe
zip_name=binaries_windows.zip

fyne package -app-id com.mouuff.smartcuts -os windows -icon images/icon.png -release

mv SmartCuts.exe $binary_name
rm $zip_name
zip -r $zip_name $binary_name