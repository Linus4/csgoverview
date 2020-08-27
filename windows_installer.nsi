;General

Unicode True
!define APP_NAME "csgoverview"
Name "${APP_NAME}"
Icon "csgoverview.ico"
OutFile "${APP_NAME}_windows_v1.0.0_install.exe"
LicenseData "LICENSE"
RequestExecutionLevel admin
; Set default InstallDir
InstallDir "$PROGRAMFILES64\${APP_NAME}"
; check string in registry and use it as the install dir if that string is valid
InstallDirRegKey HKLM "Software\${APP_NAME}" "InstallLocation"

;Pages

Page license
Page components
Page directory
Page instfiles
UninstPage uninstConfirm
UninstPage instfiles

Section "Install CSGOverview" SecCSGOverview

	SetRegView 64

    SetOutPath "$INSTDIR"

    FILE csgoverview.exe
    FILE DejaVuSans.ttf
    FILE LICENSE

    CreateDirectory $INSTDIR\assets\maps

    ;Store installation folder
    WriteRegStr HKLM "Software\${APP_NAME}" "InstallLocation" $INSTDIR

    ;Create uninstaller
    WriteUninstaller "$INSTDIR\Uninstall.exe"
	WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APP_NAME}" "DisplayName" "CSGOverview - CS:GO Demo Viewer"
	WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APP_NAME}" "UninstallString" "$\"$INSTDIR\Uninstall.exe$\""
	WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APP_NAME}" "DisplayIcon" "$\"$INSTDIR\csgoverview.exe$\""

SectionEnd

Section "Register .dem Files" SecFile

	SetRegView 64

    ;register file association
	WriteRegStr HKLM "Software\${APP_NAME}" "RegisterDemFiles" "Yes"
	;just overwrite lul
	WriteRegStr HKCR ".dem" "" "CSGOverview.demo"
	WriteRegStr HKCR "CSGOverview.demo" "" "CS:GO Demo File"
	WriteRegStr HKCR "CSGOverview.demo\DefaultIcon" "" "$INSTDIR\csgoverview.exe,0"
	WriteRegStr HKCR "CSGOverview.demo\shell" "" "open"
	WriteRegStr HKCR "CSGOverview.demo\shell\open\command" "" '"$INSTDIR\csgoverview.exe" "%1"'

SectionEnd

Section "Download Maps" SecMaps

    SetOutPath "$INSTDIR\assets\maps"

    ; DOWNLOADS
    inetc::get "https://raw.githubusercontent.com/zoidbergwill/csgo-overviews/master/overviews/de_overpass.jpg" de_overpass.jpg
    inetc::get "https://raw.githubusercontent.com/zoidbergwill/csgo-overviews/master/overviews/de_mirage.jpg" de_mirage.jpg
    inetc::get "https://raw.githubusercontent.com/zoidbergwill/csgo-overviews/master/overviews/de_vertigo.jpg" de_vertigo.jpg
    inetc::get "https://raw.githubusercontent.com/zoidbergwill/csgo-overviews/master/overviews/de_vertigo_lower.jpg" de_vertigo_lower.jpg
    inetc::get "https://raw.githubusercontent.com/zoidbergwill/csgo-overviews/master/overviews/de_nuke.jpg" de_nuke.jpg
    inetc::get "https://raw.githubusercontent.com/zoidbergwill/csgo-overviews/master/overviews/de_nuke_lower.jpg" de_nuke_lower.jpg
    inetc::get "https://raw.githubusercontent.com/zoidbergwill/csgo-overviews/master/overviews/de_cache.jpg" de_cache.jpg
    inetc::get "https://raw.githubusercontent.com/zoidbergwill/csgo-overviews/master/overviews/de_inferno.jpg" de_inferno.jpg
    inetc::get "https://raw.githubusercontent.com/zoidbergwill/csgo-overviews/master/overviews/de_train.jpg" de_train.jpg
    inetc::get "https://raw.githubusercontent.com/zoidbergwill/csgo-overviews/master/overviews/de_dust2.jpg" de_dust2.jpg

SectionEnd

;Uninstaller Section

Section "un.Uninstall"

	SetRegView 64

    Delete "$INSTDIR\assets\maps\de_overpass.jpg"
    Delete "$INSTDIR\assets\maps\de_mirage.jpg"
    Delete "$INSTDIR\assets\maps\de_vertigo.jpg"
    Delete "$INSTDIR\assets\maps\de_vertigo_lower.jpg"
    Delete "$INSTDIR\assets\maps\de_nuke.jpg"
    Delete "$INSTDIR\assets\maps\de_nuke_lower.jpg"
    Delete "$INSTDIR\assets\maps\de_cache.jpg"
    Delete "$INSTDIR\assets\maps\de_inferno.jpg"
    Delete "$INSTDIR\assets\maps\de_train.jpg"
    Delete "$INSTDIR\assets\maps\de_dust2.jpg"

    RMDIR "$INSTDIR\assets\maps"
    RMDIR "$INSTDIR\assets"

    Delete "$INSTDIR\csgoverview.exe"
    Delete "$INSTDIR\DejaVuSans.ttf"
    Delete "$INSTDIR\LICENSE"
    Delete "$INSTDIR\Uninstall.exe"

    RMDir "$INSTDIR"

    ;unregister file association
	ReadRegStr $1 HKLM "Software\${APP_NAME}" "RegisterDemFiles"
	StrCmp $1 "Yes" UnregisterFile NoUnregisterFile
	UnregisterFile:
	DeleteRegKey HKCR ".dem"
	DeleteRegKey HKCR "CSGOverview.demo"
	NoUnregisterFile:

    DeleteRegKey HKLM "Software\${APP_NAME}"
	DeleteRegKey HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APP_NAME}"

SectionEnd
