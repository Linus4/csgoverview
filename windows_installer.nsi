;General

!define APP_NAME "csgoverview"
Name "${APP_NAME}"
OutFile "${APP_NAME}_windows_v0.7.1.exe"
LicenseData "LICENSE" ;FIXME correct?
RequestExecutionLevel admin
Unicode True
InstallDir "$PROGRAMFILES\${APP_NAME}"
InstallDirRegKey HKCU "Software\${APP_NAME}" ""

;Pages

Page license
Page components
Page directory
Page instfiles
UninstPage uninstConfirm
UninstPage instfiles

Section "csgoverview" SecCSGOverview

    SetOutPath "$INSTDIR"

    FILE csgoverview.exe
    FILE DejaVuSans.ttf
    FILE LICENSE ;FIXME?

    CreateDirectory $INSTDIR\assets\maps

    ;Store installation folder
    WriteRegStr HKCU "Software\${APP_NAME}" "" $INSTDIR

    ;Create uninstaller
    WriteUninstaller "$INSTDIR\Uninstall.exe"

SectionEnd

Section "maps" SecMaps

    SetOutPath "$INSTDIR\assets\maps"

    ; DOWNLOADS
    ;inetc::get /POPUP "" /CAPTION "de_overpass.jpg" "https://raw.githubusercontent.com/zoidbergwill/csgo-overviews/master/overviews/de_overpass.jpg"
    ;Pop $0
    ;MessageBox MB_OK "Download Status: $0"

SectionEND

;Descriptions

LangString DESC_SecCSGOverview ${LANG_ENGLISH} "Install csgoverview program."
LangString DESC_SecMaps ${LANG_ENGLISH} "Download overview maps (requires internet connection)."
;FIXME need to assign descriptions to sections

;Uninstaller Section

Section "un.Uninstall"

    ;add files

    RMDir "$INSTDIR"

    DeleteRegKey /ifempty HKCU "Software\${APP_NAME}"

SectionEnd
