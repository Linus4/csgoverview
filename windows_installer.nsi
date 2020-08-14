;General

Unicode True
!include "FileAssociation.nsh"
!define APP_NAME "csgoverview"
Name "${APP_NAME}"
OutFile "${APP_NAME}_windows_v1.0.0_install.exe"
LicenseData "LICENSE"
RequestExecutionLevel admin
;set default InstallDir
InstallDir "$PROGRAMFILES\${APP_NAME}"
; check string in registry and use it as the install dir if that string is valid
InstallDirRegKey HKCU "Software\${APP_NAME}" "InstallLocation"

;Pages

Page license
Page components
Page directory
Page instfiles
UninstPage uninstConfirm
UninstPage instfiles

Section "Install csgoverview" SecCSGOverview

    SetOutPath "$INSTDIR"

    FILE csgoverview.exe
    FILE DejaVuSans.ttf
    FILE LICENSE

    CreateDirectory $INSTDIR\assets\maps

    ;Store installation folder
    WriteRegStr HKCU "Software\${APP_NAME}" "InstallLocation" $INSTDIR

    ;register file association
    ${registerExtension} "$INSTDIR\${APP_NAME}.exe" ".dem" "DEM_File"

    ;Create uninstaller
    WriteUninstaller "$INSTDIR\Uninstall.exe"

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

    DeleteRegKey /ifempty HKCU "Software\${APP_NAME}"

    ;unregister file association
    ${unregisterExtension} ".dem" "DEM_File"

SectionEnd
