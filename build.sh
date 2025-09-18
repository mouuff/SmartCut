# go build -ldflags -H=windowsgui cmd/smartcut/main.go

binary_name=smartcut_windows_amd64.exe
zip_name=binaries_windows.zip

fyne package -app-id com.mouuff.smartcut -os windows -icon images/icon.png -release

mv SmartCut.exe $binary_name
rm $zip_name
zip -r $zip_name $binary_name