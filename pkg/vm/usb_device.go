package vm

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

var (
	defaultLsusbPattern = regexp.MustCompile(`^Bus (?P<bus>\d+) Device (?P<device>\d+): ID (?P<vendor>[0-9A-Fa-f]+):(?P<productId>[0-9A-Fa-f]+) (?P<description>.+)$`)
)

type USBDevice struct {
	Description string
	Vendor      string
	ProductId   string
	Bus         int
	Device      int
}

func (dev *USBDevice) Validate() error {
	if len(dev.Vendor) == 0 {
		return errors.New("vendor is empty")
	}
	if len(dev.ProductId) == 0 {
		return errors.New("productId is empty")
	}
	if dev.Bus == 0 {
		return errors.New("bus is missing")
	}
	if dev.Device == 0 {
		return errors.New("device is missing")
	}

	return nil
}

func (dev *USBDevice) ParseFromLsusbLine(in string) error {
	match := defaultLsusbPattern.FindStringSubmatch(in)

	var err error
	for i, name := range defaultLsusbPattern.SubexpNames() {
		value := match[i]
		switch name {
		case "bus":
			dev.Bus, err = strconv.Atoi(value)
		case "device":
			dev.Device, err = strconv.Atoi(value)
		case "vendor":
			dev.Vendor = value
		case "productId":
			dev.ProductId = value
		case "description":
			dev.Description = value
		}
		if err != nil {
			return fmt.Errorf("error converting bus %q to int: %s", value, err)
		}
	}
	if err := dev.Validate(); err != nil {
		return fmt.Errorf("error validating lsusb line %q: %s", in, err)
	}

	return nil
}
