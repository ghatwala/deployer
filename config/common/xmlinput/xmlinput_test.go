package xmlinput

import (
	"fmt"
	"testing"
)

var xmldata = []byte(`<?xml version="1.0" encoding="UTF-8"?>
<input_data>
  <cpu>
  <configure>true</configure>
 	<min>1</min>
	<max>16</max>
	<default_value>1</default_value>
  </cpu>
  <ram>
  	<configure>true</configure>
  	<min>2500</min>
    <max>0</max>
	<default_value>2500</default_value>
  </ram>
  <networks>
    <configure>true</configure>
	<network name="Management" max_ifaces="1">
	    <mode type="bridged" vnic_driver="e1000"/>
		<mode type="direct" vnic_driver="e1000"/>
	</network>
	<network name="Traffic" max_ifaces="9"> 
		<mode type="bridged" vnic_driver="virtio"/>
		<mode type="passthrough"/>
		<ui_mode_selection>
			<appearance mode_type="bridged" appear="virtio"/>
			<appearance mode_type="passthrough" appear="pass-through"/>
		</ui_mode_selection>
	</network>
	<network name="Bkp" max_ifaces="2"> 
		<mode type="bridged" vnic_driver="virtio"/>
		<mode type="passthrough"/>
	</network>
  </networks>
  <nics>
	<!-- Allowed vendors and models -->
	<allow vendor="Intel" model=""/>
	<allow vendor="Broadcom" model=""/>
	<!-- Denied vendors and models -->
	<deny vendor="Broadcom" model=""/>
  </nics>
</input_data>`)

func TestParseXMLInput(t *testing.T) {
	d, err := ParseXMLInputMock(xmldata)
	if err != nil {
		t.Fatal(err)
	}
	for _, nic := range d.Allowed {
		fmt.Printf("\nAllowed : NIC Vendor =>%s|NIC Model => %s\n",
			nic.Vendor, nic.Model)
	}
	for _, nic := range d.Denied {
		fmt.Printf("\nDenied : NIC Vendor =>%s|NIC Model => %s\n",
			nic.Vendor, nic.Model)
	}
}

var bad_xmldata = []byte(`<?xml version="1.0" encoding="UTF-8"?>
<input_data>
  <cpu>
  <configure>true</configure>
 	<min>1</min>
	<max>16</max>
	<default_value>1</default_value>
  </cpu>
  <ram>
  	<configure>true</configure>
  	<min>2500</min>
    <max>0</max>
	<default_value>2500</default_value>
  </ram>
  <networks>
    <configure>true</configure>
	<network name="Management" max_ifaces="1">
	    <mode type="bridged" vnic_driver="e1000"/>
		<mode type="direct" vnic_driver="e1000"/>
		<ui_mode_selection enable="false"/>
	</network>
	<network name="Traffic" max_ifaces="9"> 
		<mode type="bridged" vnic_driver="virtio"/>
		<mode type="passthrough"/>
		<mode type="direct"/>
		<ui_mode_selection enable="true">
			<appearance mode_type="bridged" appear="virtio"/>
			<appearance mode_type="passthrough" appear="passthrough"/>
		</ui_mode_selection>
	</network>
	<network name="Bkp" max_ifaces="2"> 
		<mode type="bridged" vnic_driver="virtio"/>
		<mode type="passthrough"/>
	</network>
  </networks>
  <nics>
	<!-- Allowed vendors and models -->
	<allow vendor="Intel" model=""/>
	<allow vendor="Broadcom" model=""/>
	<!-- Denied vendors and models -->
	<deny vendor="Broadcom" model=""/>
  </nics>
</input_data>`)

func TestParseXMLInputBad(t *testing.T) {
	_, err := ParseXMLInputMock(bad_xmldata)
	if err == nil {
		t.Fatalf("supposed to produce an error")
	}
}