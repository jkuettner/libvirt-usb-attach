package main

import (
	"bufio"
	"fmt"
	"github.com/jkuettner/libvirt-usb-attach/pkg/vm"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strings"
)

var (
	libvirtdAddress string

	vendor    string
	productId string
	bus       int
	device    int

	parseLsusbLine bool
	verbose        bool
)

func init() {
	rootCmd.Flags().StringVarP(&libvirtdAddress, "address", "a", "qemu:///system", "libvirtd address")
	rootCmd.Flags().StringVarP(&vendor, "vendor-id", "V", "", "Vendor ID")
	rootCmd.Flags().StringVarP(&productId, "product-id", "P", "", "Product ID")

	rootCmd.Flags().IntVarP(&bus, "bus", "b", 0, "Bus")
	rootCmd.Flags().IntVarP(&device, "device", "d", 0, "Device")

	rootCmd.Flags().BoolVarP(&parseLsusbLine, "parse-lsusb-line", "p", false, `Parse vendor-id, product-id, bus and device from "lsusb" output`)
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, `enable verbose output`)
}

var rootCmd = &cobra.Command{
	Use:  "libvirtusbattach [VM/DOMAIN NAME]",
	Long: "attach an usb device by to the virtual machine/domain",
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		vmName := args[0]
		attacher, err := vm.NewUSBAttacher(libvirtdAddress)
		if err != nil {
			log.Fatal(err)
		}
		usbDev := &vm.USBDevice{
			Vendor:    vendor,
			ProductId: productId,
			Bus:       bus,
			Device:    device,
		}
		if parseLsusbLine {
			reader := bufio.NewReader(os.Stdin)
			lsusbLine, err := reader.ReadString('\n')
			if err != nil {
				log.Fatal(fmt.Errorf("error reading lsusb-line: %v", err))
			}
			if strings.HasSuffix(lsusbLine, "\n") {
				lsusbLine = lsusbLine[:len(lsusbLine)-1]
			}
			if err := usbDev.ParseFromLsusbLine(lsusbLine); err != nil {
				log.Fatal(fmt.Errorf("error parsing lsusb line: %s", err))
			}
		}
		if err := usbDev.Validate(); err != nil {
			log.Fatal(err)
		}

		log.Printf("attaching %q (%d.%d %s:%s) to vm %s", usbDev.Description, usbDev.Bus, usbDev.Device, usbDev.Vendor, usbDev.ProductId, vmName)
		if err := attacher.AttachDevice(vmName, usbDev); err != nil {
			log.Fatal(err)
		}
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
