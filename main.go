package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"slices"

	"github.com/spf13/pflag"
	"libvirt.org/go/libvirtxml"

	vmSchema "kubevirt.io/api/core/v1"
)

const (
	bootAnnotation = "os.vm.kubevirt.io/boot"
)

type BootConfig struct {
	BootDevices []DomainBootDevice `json:"boot"`
	BootMenu    *DomainBootMenu    `json:"bootmenu,omitempty"`
}

// https://libvirt.org/formatdomain.html#operating-system-booting
var validateDevs = []string{"fd", "hd", "cdrom", "network"}

type DomainBootDevice struct {
	Dev string `json:"dev"`
}

type DomainBootMenu struct {
	Enable  string `json:"enable"`
	Timeout string `json:"timeout"` // milliseconds
}

func onDefineDomain(vmiJSON, domainXML []byte) (string, error) {
	vmiSpec := vmSchema.VirtualMachineInstance{}
	if err := json.Unmarshal(vmiJSON, &vmiSpec); err != nil {
		return "", fmt.Errorf("Failed to unmarshal given VMI spec: %s %s", err, string(vmiJSON))
	}

	domainSpec := libvirtxml.Domain{}
	if err := xml.Unmarshal(domainXML, &domainSpec); err != nil {
		return "", fmt.Errorf("Failed to unmarshal given Domain spec: %s %s", err, string(domainXML))
	}

	annotations := vmiSpec.GetAnnotations()
	bootConfigAnnotation, found := annotations[bootAnnotation]
	if !found {
		return string(domainXML), nil
	}

	var bootConfig BootConfig
	if err := json.Unmarshal([]byte(bootConfigAnnotation), &bootConfig); err != nil {
		return "", fmt.Errorf("Failed to unmarshal given bootAnnotation value: %s %s", err, bootConfigAnnotation)
	}

	if domainSpec.OS == nil {
		domainSpec.OS = &libvirtxml.DomainOS{}
	}
	if len(bootConfig.BootDevices) > 0 {
		bootDevices := make([]libvirtxml.DomainBootDevice, 0)
		for _, dev := range bootConfig.BootDevices {
			if !slices.Contains(validateDevs, dev.Dev) {
				continue
			}
			bootDevices = append(bootDevices, libvirtxml.DomainBootDevice{
				Dev: dev.Dev,
			})
		}
		domainSpec.OS.BootDevices = bootDevices
	}

	if bootConfig.BootMenu != nil {
		domainSpec.OS.BootMenu = &libvirtxml.DomainBootMenu{
			Enable:  bootConfig.BootMenu.Enable,
			Timeout: bootConfig.BootMenu.Timeout,
		}
	}

	newDomainXML, err := xml.Marshal(domainSpec)
	if err != nil {
		return "", fmt.Errorf("Failed to marshal new Domain spec: %s %+v", err, domainSpec)
	}

	return string(newDomainXML), nil
}

func main() {
	var vmiJSON, domainXML string
	pflag.StringVar(&vmiJSON, "vmi", "", "VMI to change in JSON format")
	pflag.StringVar(&domainXML, "domain", "", "Domain spec in XML format")
	pflag.Parse()

	logger := log.New(os.Stderr, "boot", log.Ldate)
	if vmiJSON == "" || domainXML == "" {
		logger.Printf("Bad input vmi=%d, domain=%d", len(vmiJSON), len(domainXML))
		os.Exit(1)
	}

	domainXML, err := onDefineDomain([]byte(vmiJSON), []byte(domainXML))
	if err != nil {
		logger.Printf("onDefineDomain failed: %s", err)
		panic(err)
	}
	fmt.Println(domainXML)
}
