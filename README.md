# noolite
Package noolite provide class for control Noolite Adapters PC11xx.

Protocol described on url: http://www.noo.com.by/assets/files/software/PC11xx_HID_API.pdf

To have access on device from common user add the next rule to udev. For example to /etc/udev/rules.d/50-noolite.rules next line:
```
ATTRS{idVendor}=="16c0", ATTRS{idProduct}=="05df", SUBSYSTEMS=="usb", ACTION=="add", MODE="0666", GROUP="noolite"
```
Then add your user to `noolite` group:
```
sudo usermod <user> -aG noolite
```
