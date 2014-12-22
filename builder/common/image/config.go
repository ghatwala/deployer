// Parses image configuration file (XML)

// Configuration example:
//
//<?xml version="1.0" encoding="UTF-8"?>
//<storage>
//  <config>
//	 <disk>
//	  	<size_gb>5</size_gb>
//    	<bootable>true</bootable>
//	 	 <fdisk_cmd>n\np\n1\n\n+3045M\nn\np\n2\n\n\nt\n2\n82\na\n1\nw\n</fdisk_cmd>
//   	 <description>Topology for release xxxx</description>
//  	 <partition>
//	 	    <sequence>1</sequence>
//	 	    <size_mb>3045</size_mb>
//   	    <label>SLASH</label>
//   	    <mount_point>/</mount_point>
//   	    <file_system>ext4</file_system>
//	 	    <file_system_args></file_system_args>
//	 	 </partition>
//	 	 <partition>
//	 	    <sequence>2</sequence>
//	 	    <size_mb>400</size_mb>
//   	    <label>SWAP</label>
//   	    <mount_point>SWAP</mount_point>
//   	    <file_system>swap</file_system>
//	 	    <file_system_args></file_system_args>
//	 	 </partition>
// 	 </disk>
// </config>
//</storage>`

package image

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
)

type StorageType string

const (
	StorageTypeRAW   StorageType = "raw"
	StorageTypeQCOW2 StorageType = "qcow2"
	StorageTypeVMDK  StorageType = "vmdk"
)

type ConfigIndex uint8

type Storage struct {
	Configs []*Config `xml:"config"`
}

type Config struct {
	Disks []*Disk `xml:"disk"`
}

type Disk struct {
	Path        string
	Type        StorageType  `xml:"storage_type"`
	SizeGb      int          `xml:"size_gb"`
	Bootable    bool         `xml:"bootable"`
	FdiskCmd    string       `xml:"fdisk_cmd"`
	Description string       `xml:"description"`
	Partitions  []*Partition `xml:"partition"`
}

type Partition struct {
	Sequence       int    `xml:"sequence"`
	BootFlag       bool   `xml: "boot_flag"`
	SizeMb         int    `xml:"size_mb"`
	Label          string `xml:"label"`
	MountPoint     string `xml:"mount_point"`
	FileSystem     string `xml:"file_system"`
	FileSystemArgs string `xml:"file_system_args"`
	Description    string `xml:"description"`
}

// ParseConfigFile is responsible for reading appropriate XML file
// and calling ParseConfig for further processing
func ParseConfigFile(xmlpath string) (*Storage, error) {
	fb, err := ioutil.ReadFile(xmlpath)
	if err != nil {
		return nil, err
	}
	return ParseConfig(fb)
}

// ParseConfig is responsible for processing XML content
func ParseConfig(fb []byte) (*Storage, error) {
	buf := bytes.NewBuffer(fb)
	p := new(Storage)
	decoded := xml.NewDecoder(buf)
	if err := decoded.Decode(p); err != nil {
		return nil, err
	}
	return p, nil
}

// DiskConfig returns disk configuration for appropriate index
func (s *Storage) IndexToConfig(index ConfigIndex) (*Config, error) {
	t := s.Configs[index]
	if t == nil {
		return nil, fmt.Errorf("no configuration found for index %d", index)
	}
	return t, nil
}
