csgoverview installation instructions
=====================================

1. Download latest installer from the [releases page](https://github.com/Linus4/csgoverview/releases).
2. Important: To verify integrity of the installer and the program:
    - Run `Command Prompt` from Windows Start Menu.
    - Navigate to the installer you just downloaded using the following commands:
        + `c:`: switches to C partition (for example)
        + `cd <dir>`: changes the directory (e.g. `cd C:\Users\Linus\Downloads`)
    - Run `CertUtil -hashfile csgoverview_windows_v1.0.0_install.exe SHA256 | findstr -v "hash"` on the command
  line and make sure the checksum matches the one provided [here](https://github.com/Linus4/csgoverview/blob/7747114bf2419d19aae4f4435428c07049bd8412/README_WINDOWS.txt#L20).
    - If the numbers don't match, try downloading the installer again and if they still don't match please let me know in the [chat](https://gitter.im/csgoverview/community?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge) and do NOT install csgoverview.
3. Run the installer.
4. Right click a demo and select 'Open with' to open it with csgoverview.

Checksums
=========

v1.0.0: 87f439d4e1097d534e799d576213b653a744e77b97b62cf3225e0dde614e1cfa


Updates
=======

Download the latest release from https://github.com/Linus4/csgoverview/releases.
Uninstall the old version. Dont' forget to verify the checksum (!) and then run 
the installer you just downloaded.

You can watch the releases of this project on the github page to be notified when
a new version is released.

When there is an update to a map that changes the layout, you need to download
the overview image from https://github.com/zoidbergwill/csgo-overviews again or
run the installer again (it might take a while before the new overview images are
available online.
