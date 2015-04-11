// Package dmidecode extracts useful info from that command on linux
//
// Csv output is particularly supported, so that a csvfile-based enterprise's ETL tools can also monitor its servers and desktops
package dmidecode

import (
	"fmt"
	"github.com/LDCS/genutil"
	"sort"
	"strings"
)

// Dmidecodedata hold dmidecode data
type Dmidecodedata struct {
	// "path":"size":"transport-serialnumber":"logical-sector-size":"physical-sector-size":"partition-table-serialnumber":"model-name";
	Manufacturer_ string
	Productname_  string
	Serialnumber_ string
	Uuid_         string
}

// SortedKeys_String2PtrDmidecodedata is generic
func SortedKeys_String2PtrDmidecodedata(_mp *map[string]*Dmidecodedata) []string {
	keys := make([]string, len(*_mp))
	ii := 0
	for kk := range *_mp {
		keys[ii] = kk
		ii++
	}
	sort.Strings(keys)
	return keys
}

// Keys_String2PtrDmidecodedata is generic
func Keys_String2PtrDmidecodedata(_mp *map[string]*Dmidecodedata) []string {
	keys := make([]string, len(*_mp))
	ii := 0
	for kk := range *_mp {
		keys[ii] = kk
		ii++
	}
	return keys
}

// Header is generic
func Header() string { return fmt.Sprintf("dc.Manufacturer,dc.Productname,dc.Serialnumber,dc.Uuid") }

// Csv is generic
func (self *Dmidecodedata) Csv() string {
	if self == nil {
		return ",,,"
	}
	return fmt.Sprintf("%s,%s,%s,%s",
		self.Manufacturer_, self.Productname_, self.Serialnumber_, self.Uuid_)
}

// Print is generic
func (self *Dmidecodedata) Print() {
	if self == nil {
		return
	}
	fmt.Printf("Manufacturer=%s Productname=%s Serialnumber=%s Uuid=%s\n",
		self.Manufacturer_, self.Productname_, self.Serialnumber_, self.Uuid_)
}

// Sprint is generic
func (self *Dmidecodedata) Sprint() string {
	if self == nil {
		return ""
	}
	return fmt.Sprintf("Manufacturer=%s Productname=%s Serialnumber=%s Uuid=%s\n",
		self.Manufacturer_, self.Productname_, self.Serialnumber_, self.Uuid_)
}

// cleanItem cleans an item
func cleanItem(_str string) string {
	_str = strings.TrimSpace(_str)
	_str = strings.Replace(_str, ",", ";", -1)
	_str = strings.Replace(_str, "To be filled by ", "", -1)
	_str = strings.Replace(_str, "O.E.M.", "OEM", -1)
	if strings.HasPrefix(_str, "Gigabyte") {
		_str = "Gigabyte"
	}
	return _str
}

// DoListDmidecodedata extracts dmidecode data
func DoListDmidecodedata(_verbose bool) *Dmidecodedata {
	out := genutil.BashExecOrDie(_verbose, "/usr/bin/timeout 10 /usr/sbin/dmidecode -t System", ".")
	if _verbose {
		fmt.Println(out)
	}
	lines := genutil.CleanAndSplitOnSpaces(out, " ")
	var lastdc *Dmidecodedata
	seenSystemInformation := false
	lastdc = new(Dmidecodedata)
	for ii, lineraw := range lines {
		line := strings.Replace(strings.TrimSpace(lineraw), ",", ";", -1)
		items := strings.Split(line, ":")
		nitems := len(items)
		if nitems < 1 {
			continue
		}
		if _verbose {
			fmt.Printf("line%d=%s\n", ii, line)
		}
		if !seenSystemInformation {
			switch {
			case (nitems > 0) && (items[0] == "System Information"):
				seenSystemInformation = true
				continue
			default:
				continue
			}
		} else {
			if strings.HasPrefix(items[0], "Handle") {
				break
			}
		}
		if _verbose {
			fmt.Printf("LINE%d=%s\n", ii, line)
		}
		switch items[0] {
		case "Manufacturer":
			lastdc.Manufacturer_ = cleanItem(items[1])
		case "Product Name":
			lastdc.Productname_ = cleanItem(items[1])
		case "Serial Number":
			lastdc.Serialnumber_ = cleanItem(items[1])
		case "UUID":
			lastdc.Uuid_ = cleanItem(items[1])
		default:
			if _verbose {
				fmt.Printf("line%d: item0(%s) %s\n", ii, items[0], strings.Join(items, "~"))
			}
		}
	}
	return lastdc
}
