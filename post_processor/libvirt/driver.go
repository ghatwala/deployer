package libvirt

import (
	"fmt"
	"strings"
	"sync"

	"github.com/clbanning/mxj"
	"github.com/dorzheh/deployer/utils"
	ssh "github.com/dorzheh/infra/comm/common"
)

type Driver struct {
	sync.Mutex
	run func(string) (string, error)
}

func NewDriver(config *ssh.Config) *Driver {
	d := new(Driver)
	d.run = utils.RunFunc(config)
	return d
}

func (d *Driver) DefineDomain(domainConfig string) error {
	d.Lock()
	defer d.Unlock()

	if _, err := d.run("virsh define " + domainConfig); err != nil {
		return err
	}
	return nil
}

func (d *Driver) StartDomain(name string) error {
	d.Lock()
	defer d.Unlock()

	if _, err := d.run("virsh start " + name); err != nil {
		return err
	}
	return nil
}

func (d *Driver) DestroyDomain(name string) error {
	d.Lock()
	defer d.Unlock()

	if _, err := d.run("virsh destroy " + name); err != nil {
		return err
	}
	return nil
}

func (d *Driver) UndefineDomain(name string) error {
	d.Lock()
	defer d.Unlock()

	if _, err := d.run("virsh undefine " + name); err != nil {
		return err
	}
	return nil
}

func (d *Driver) SetAutostart(name string) error {
	d.Lock()
	defer d.Unlock()

	if _, err := d.run("virsh autostart " + name); err != nil {
		return err
	}
	return nil
}

func (d *Driver) DomainExists(name string) bool {
	d.Lock()
	defer d.Unlock()

	if _, err := d.run("virsh dominfo " + name); err != nil {
		return false
	}
	return true
}

// Emulator returns appropriate path to QEMU emulator for a given architecture
func (d *Driver) Emulator(arch string) (string, error) {
	switch arch {
	case "x86_64":
	case "i686":
	default:
		return "", fmt.Errorf("Unsupported architecture(%s).Supported i686 and x86_64 only", arch)
	}

	out, err := d.run("virsh capabilities")
	if err != nil {
		return "", err
	}

	m, err := mxj.NewMapXml([]byte(out))
	if err != nil {
		return "", err
	}

	v, _ := m.ValuesForPath("capabilities.guest.arch", "-name:"+arch)
	return v[0].(map[string]interface{})["emulator"].(string), nil
}

// Version returns libvirt API version
func (d *Driver) Version() (string, error) {
	out, err := d.run("libvirtd --version")
	if err != nil {
		return "", err
	}
	return strings.Split(out, " ")[2], nil

}
