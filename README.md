# libvirt-usb-attach

`libvirt-usb-attach` is a command-line tool to attach a USB device to a specified libvirt virtual machine (VM) or domain. It simplifies the process of connecting USB devices to VMs managed by libvirt.

## Usage

```bash
attach an usb device by to the virtual machine/domain

Usage:
  libvirtusbattach [VM/DOMAIN NAME] [flags]

Flags:
  -a, --address string      libvirtd address (default "qemu:///system")
  -b, --bus int             Bus
  -d, --device int          Device
  -h, --help                help for libvirtusbattach
  -p, --parse-lsusb-line    Parse vendor-id, product-id, bus and device from "lsusb" output
  -P, --product-id string   Product ID
  -V, --vendor-id string    Vendor ID
  -v, --verbose             enable verbose output
```

## Examples

In this example we want to attach the following usb device (output from `lsusb`):
```bash
Bus 003 Device 002: ID 2357:0604 TP-Link TP-Link UB500 Adapter
```
### Attach device by bus and device number, vendor- and product-id
```bash
libvirtusbattach <my-vm-name> --bus 3 --device 2 --vendor-id 2357 --product-id 0604
```
or
```bash
libvirtusbattach <my-vm-name> -b 3 -d 2 -V 2357 -P 0604
```

### Attach a USB Device Using `lsusb` output line

```bash
echo "Bus 003 Device 002: ID 2357:0604 TP-Link TP-Link UB500 Adapter" | libvirtusbattach <my-vm-name> -p
```

### Parse USB Device Information from `lsusb` Output

### Specify a Custom Libvirt Address

```bash
libvirtusbattach my-vm-name \
  --address qemu+ssh://user@remote-system/system \
  --vendor-id 1234 \
  --product-id 5678
```

### Use in a Waybar Script

The tool can be used to create a script for Waybar. For example:

```json
"custom/myvm": {
      "on-click": "lsusb | wofi -i --dmenu | libvirtusbattach <my-vm> -p",
},
```
