package vm

import (
	"fmt"
	"libvirt.org/go/libvirt"
)

type USBAttacher struct {
	libvirtConnect *libvirt.Connect
}

func NewUSBAttacher(address string) (*USBAttacher, error) {
	conn, err := libvirt.NewConnect(address)
	if err != nil {
		return nil, fmt.Errorf("error connecting to libvirt: %s", err)
	}

	return &USBAttacher{
		libvirtConnect: conn,
	}, nil
}

func (a *USBAttacher) AttachDevice(vmName string, device *USBDevice) error {
	vm, err := a.libvirtConnect.LookupDomainByName(vmName)
	if err != nil {
		return fmt.Errorf("error looking up domain %s: %s", vmName, err)
	}

	if err := vm.AttachDevice(a.buildAttachXML(device.Vendor, device.ProductId, device.Bus, device.Device)); err != nil {
		return fmt.Errorf("error attaching device %s:%s to vm/domain %s: %s", device.Vendor, device.ProductId, vmName, err)
	}
	return nil
}

func (a *USBAttacher) buildAttachXML(vendor, productId string, bus, device int) string {
	return fmt.Sprintf(`
<hostdev mode='subsystem' type='usb' managed='yes'>
      <source>
        <vendor id='0x%s'/>
        <product id='0x%s'/>
        <address bus='%d' device='%d'/>
      </source>
      <alias name='foobar'/>
</hostdev>
`, vendor, productId, bus, device)
}
