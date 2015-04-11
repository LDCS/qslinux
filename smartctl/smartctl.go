// Package parted extracts smartctl information on linux (hp and supermicro are handled better)
//
// Csv output is particularly supported, so that a csvfile-based enterprise's ETL tools can also monitor its servers and desktops
package smartctl

import (
	"fmt"
	"github.com/LDCS/genutil"
	"github.com/LDCS/qslinux/df"
	"github.com/LDCS/qslinux/dmidecode"
	"github.com/LDCS/qslinux/hp"
	"github.com/LDCS/qslinux/parted"
	"github.com/LDCS/qslinux/scsi"
	"sort"
	"strings"
)

// Smartctldata holds smartctl data
type Smartctldata struct {
	Vendor_                                string
	Product_                               string
	Revision_                              string
	Usercapacity_                          string
	Logicalblocksize_                      string
	Logicalunitid_                         string
	Serialnumber_                          string
	Devicetype_                            string
	Transportprotocol_                     string
	Localtimeis_                           string
	Smarthealthstatus_                     string
	Currentdrivetemperature_               string
	Drivetriptemperature_                  string
	Specifiedcyclecountoverdevicelifetime_ string
	Accumulatedstartstopcycles_            string
	Elementsingrowndefectlist_             string
	Nonmediumerrorcount_                   string
	Errorsread_                            string
	Errorswrite_                           string
}

const (
	names     = "Vendor,Product,Revision,Usercapacity,Logicalblocksize,Logicalunitid,Serialnumber,Devicetype,Transportprotocol,Localtimeis,Smarthealthstatus,Currentdrivetemperature,Drivetriptemperature,Specifiedcyclecountoverdevicelifetime,Accumulatedstartstopcycles,Elementsingrowndefectlist,Nonmediumerrorcount,Errorsread,Errorswrite"
	hdrprefix = ",sc."
)

var (
	headerString  string
	commaString   string
	pctString     string
	namePctString string
)

// init  is generic
func init() {
	headerString = (hdrprefix + strings.Join(strings.Split(names, ","), hdrprefix))[1:] // fmt.Sprintf("sc.Vendor,sc.Product,sc.Revision,sc.Usercapacity,sc.Logicalblocksize,sc.Logicalunitid,sc.Serialnumber,sc.Devicetype,sc.Transportprotocol,sc.Localtimeis,sc.Smarthealthstatus,sc.Currentdrivetemperature,sc.Drivetriptemperature,sc.Specifiedcyclecountoverdevicelifetime,sc.Accumulatedstartstopcycles,sc.Elementsingrowndefectlist,sc.Nonmediumerrorcount,sc.Errorsread,sc.Errorswrite") }
	commaString = strings.Repeat(",", strings.Count(headerString, ","))
	pctString = strings.Repeat(",%s", 1+strings.Count(headerString, ","))[1:]
	namePctString = strings.Replace(names, ",", "=%s ", -1) + "=%s\n"
}

// SortedKeys_String2PtrSmartctldata is generic
func SortedKeys_String2PtrSmartctldata(_mp *map[string]*Smartctldata) []string {
	keys := make([]string, len(*_mp))
	ii := 0
	for kk := range *_mp {
		keys[ii] = kk
		ii++
	}
	sort.Strings(keys)
	return keys
}

// Keys_String2PtrSmartctldata is generic
func Keys_String2PtrSmartctldata(_mp *map[string]*Smartctldata) []string {
	keys := make([]string, len(*_mp))
	ii := 0
	for kk := range *_mp {
		keys[ii] = kk
		ii++
	}
	return keys
}

// Header is generic
func Header() string { return headerString }

// Csv is generic
func (self *Smartctldata) Csv() string {
	if self == nil {
		return commaString
	}
	return fmt.Sprintf(pctString, self.Vendor_, self.Product_, self.Revision_, self.Usercapacity_, self.Logicalblocksize_, self.Logicalunitid_, self.Serialnumber_, self.Devicetype_, self.Transportprotocol_, self.Localtimeis_, self.Smarthealthstatus_, self.Currentdrivetemperature_, self.Drivetriptemperature_, self.Specifiedcyclecountoverdevicelifetime_, self.Accumulatedstartstopcycles_, self.Elementsingrowndefectlist_, self.Nonmediumerrorcount_, self.Errorsread_, self.Errorswrite_)
}

// Sprint is generic
func (self *Smartctldata) Sprint() string {
	if self == nil {
		return ""
	}
	return fmt.Sprintf(namePctString, self.Vendor_, self.Product_, self.Revision_, self.Usercapacity_, self.Logicalblocksize_, self.Logicalunitid_, self.Serialnumber_, self.Devicetype_, self.Transportprotocol_, self.Localtimeis_, self.Smarthealthstatus_, self.Currentdrivetemperature_, self.Drivetriptemperature_, self.Specifiedcyclecountoverdevicelifetime_, self.Accumulatedstartstopcycles_, self.Elementsingrowndefectlist_, self.Nonmediumerrorcount_, self.Errorsread_, self.Errorswrite_)
}

// Print is generic
func (self *Smartctldata) Print() {
	if self == nil {
		return
	}
	fmt.Printf(fmt.Sprint())
}

// mergeSmartctl is generic
func mergeSmartctl(_jj int, _sum, _add *Smartctldata) {
	_sum.Vendor_ += ";" + _add.Vendor_
	_sum.Product_ += ";" + _add.Product_
	_sum.Revision_ += ";" + _add.Revision_
	_sum.Usercapacity_ += ";" + _add.Usercapacity_
	_sum.Logicalblocksize_ += ";" + _add.Logicalblocksize_
	_sum.Logicalunitid_ += ";" + _add.Logicalunitid_
	_sum.Serialnumber_ += ";" + _add.Serialnumber_
	_sum.Devicetype_ += ";" + _add.Devicetype_
	_sum.Transportprotocol_ += ";" + _add.Transportprotocol_
	_sum.Localtimeis_ += ";" + _add.Localtimeis_
	_sum.Smarthealthstatus_ += ";" + _add.Smarthealthstatus_
	_sum.Currentdrivetemperature_ += ";" + _add.Currentdrivetemperature_
	_sum.Drivetriptemperature_ += ";" + _add.Drivetriptemperature_
	_sum.Specifiedcyclecountoverdevicelifetime_ += ";" + _add.Specifiedcyclecountoverdevicelifetime_
	_sum.Accumulatedstartstopcycles_ += ";" + _add.Accumulatedstartstopcycles_
	_sum.Elementsingrowndefectlist_ += ";" + _add.Elementsingrowndefectlist_
	_sum.Nonmediumerrorcount_ += ";" + _add.Nonmediumerrorcount_
	_sum.Errorsread_ += ";" + _add.Errorsread_
	_sum.Errorswrite_ += ";" + _add.Errorswrite_
}

// DoListSmartctldataOne is a loop over smartctl controllers
func DoListSmartctldataOne(_df *df.Dfdata, _scsi *scsi.Scsidata, _parted *parted.Parteddata, _dmidecode *dmidecode.Dmidecodedata, _hp *hp.Hpdata, _verbose bool) *Smartctldata {
	switch {
	case _dmidecode == nil:
		if _verbose {
			fmt.Println("DoListSmartctldataOne: Skipping: dmidecode is nil\n")
		}
	case _scsi == nil:
		if _verbose {
			fmt.Println("DoListSmartctldataOne: Skipping: scsi is nil\n")
		}
	case (_df != nil) && (_df.Type_ == "network"):
		if _verbose {
			fmt.Println("DoListSmartctldataOne: Skipping: df.Type is network\n")
		}
	case (_df != nil) && (_df.Type_ == "tmpfs"):
		if _verbose {
			fmt.Println("DoListSmartctldataOne: Skipping: df.Type is tmpfs\n")
		}
	case (_df != nil) && (_df.Type_ == "none"):
		if _verbose {
			fmt.Println("DoListSmartctldataOne: Skipping: df.Type is none\n")
		}
	default:
		if _verbose {
			fmt.Println("DoListSmartctldataOne: found default\n")
		}
		if (_dmidecode.Manufacturer_ == "HP") && (_scsi.Devicetype_ == "storage") && (_scsi.Vendor_ == "HP") { // Storage controller such as P410
			if _verbose {
				fmt.Println("DoListSmartctldataOne: found HP controller \n")
			}
			return DoListSmartctldataOneHP(_df, _scsi, _parted, _dmidecode, _verbose)
		}
		if (_dmidecode.Manufacturer_ == "HP") && (_scsi.Devicetype_ == "disk") && (_scsi.Vendor_ == "HP") && (len(_scsi.Generic_) > 0) && (_parted != nil) /* && (_parted.Type_ == "rawdevice")*/ { // Logical device on a storage controller such as P410
			if _verbose {
				fmt.Println("DoListSmartctldataOne: found HP logical disk\n")
			}
			return DoListSmartctldataOneHPDisk(_df, _scsi, _parted, _dmidecode, _verbose)
		}
		if (_dmidecode.Manufacturer_ == "Supermicro") && (_scsi.Devicetype_ == "disk") && (_parted != nil) /* && (!genutil.StrSin(_parted.Type_, "partition|softraid"))*/ {
			if _verbose {
				fmt.Println("DoListSmartctldataOne: found Supermicro \n")
			}
			return DoListSmartctldataOneSupermicro(_df, _scsi, _parted, _dmidecode, _verbose)
		}
	}
	return new(Smartctldata)
}

// DoListSmartctldataOneHP does one HP controller
func DoListSmartctldataOneHP(_df *df.Dfdata, _scsi *scsi.Scsidata, _parted *parted.Parteddata, _dmidecode *dmidecode.Dmidecodedata, _verbose bool) *Smartctldata {
	sc := new(Smartctldata)

	doStop := false
	for jj := 0; jj < 100; jj++ {
		cmdstr := fmt.Sprintf("/usr/bin/timeout 10 /usr/sbin/smartctl -a -d cciss,%d %s", jj, _scsi.Generic_)
		if _verbose {
			fmt.Println("cmdstr is ", cmdstr)
		}
		out := genutil.BashExecOrDie(_verbose, cmdstr, ".")
		if _verbose {
			fmt.Println(out)
		}
		lines := genutil.CleanAndSplitOnSpaces(out, " ")
		lastsc := new(Smartctldata)
		for ii, lineraw := range lines {
			line := strings.Replace(strings.TrimSpace(lineraw), ",", ";", -1)
			if _verbose {
				fmt.Printf("jj=%d LINE%d=%s\n", jj, ii, line)
			}
			if strings.Contains(line, "mandatory SMART command failed:") {
				doStop = true
				break
			}
			items := strings.Split(line, ":")
			if len(items) < 1 {
				continue
			}
			switch {
			case items[0] == "Vendor":
				lastsc.Vendor_ = items[1]
			case items[0] == "Product":
				lastsc.Product_ = items[1]
			case items[0] == "Revision":
				lastsc.Revision_ = items[1]
			case items[0] == "User Capacity":
				lastsc.Usercapacity_ = items[1]
			case items[0] == "Logical block size":
				lastsc.Logicalblocksize_ = items[1]
			case items[0] == "Logical Unit id":
				lastsc.Logicalunitid_ = items[1]
			case items[0] == "Serial number":
				lastsc.Serialnumber_ = items[1]
			case items[0] == "Device type":
				lastsc.Devicetype_ = items[1]
			case items[0] == "Transport protocol":
				lastsc.Transportprotocol_ = items[1]
			case items[0] == "Local Time is": // lastsc.Localtimeis_	= items[1]
			case items[0] == "SMART Health Status":
				lastsc.Smarthealthstatus_ = items[1]
			case items[0] == "Current Drive Temperature":
				lastsc.Currentdrivetemperature_ = items[1]
			case items[0] == "Drive Trip Temperature":
				lastsc.Drivetriptemperature_ = items[1]
			case items[0] == "Specified cycle count over device lifetime":
				lastsc.Specifiedcyclecountoverdevicelifetime_ = items[1]
			case items[0] == "Accumulated start-stop cycles":
				lastsc.Accumulatedstartstopcycles_ = items[1]
			case items[0] == "Elements in grown defect list":
				lastsc.Elementsingrowndefectlist_ = items[1]
			case items[0] == "Non-medium error count":
				lastsc.Nonmediumerrorcount_ = items[1]
			case items[0] == "read":
				lastsc.Errorsread_ = items[1]
			case items[0] == "write":
				lastsc.Errorswrite_ = items[1]
			default:
				if _verbose {
					fmt.Printf("line%d: %s\n", ii, strings.Join(items, "!"))
				}
			}
		}
		if doStop {
			if (jj > 0) && (_df != nil) && (len(_df.Type_) == 0) {
				_df.Type_ = "controller"
			}
			break
		}
		mergeSmartctl(jj, sc, lastsc)
	}
	return sc
}

// DoListSmartctldataOneHPDisk extracts smartctl for HP
func DoListSmartctldataOneHPDisk(_df *df.Dfdata, _scsi *scsi.Scsidata, _parted *parted.Parteddata, _dmidecode *dmidecode.Dmidecodedata, _verbose bool) *Smartctldata {
	sc := new(Smartctldata)

	if true {
		cmdstr := fmt.Sprintf("/usr/bin/timeout 10 /usr/sbin/smartctl -iH %s", _scsi.Generic_)
		if _verbose {
			fmt.Println("cmdstr is ", cmdstr)
		}
		out := genutil.BashExecOrDie(_verbose, cmdstr, ".")
		if _verbose {
			fmt.Println(out)
		}
		lines := genutil.CleanAndSplitOnSpaces(out, " ")
		for ii, lineraw := range lines {
			line := strings.Replace(strings.TrimSpace(lineraw), ",", ";", -1)
			if _verbose {
				fmt.Printf("LINE%d=%s\n", ii, line)
			}
			if strings.Contains(line, "device is NOT READY") {
				sc.Smarthealthstatus_ = "FAILED"
				break
			}
			items := strings.Split(line, ":")
			if len(items) < 2 {
				continue
			}
			item1 := strings.TrimSpace(items[1])
			switch {
			case items[0] == "Vendor":
				sc.Vendor_ = item1
			case items[0] == "Product":
				sc.Product_ = item1
			case items[0] == "Revision":
				sc.Revision_ = item1
			case items[0] == "User Capacity":
				sc.Usercapacity_ = item1
			case items[0] == "Logical block size":
				sc.Logicalblocksize_ = item1
			case items[0] == "Logical Unit id":
				sc.Logicalunitid_ = item1
			case items[0] == "Serial number":
				sc.Serialnumber_ = item1
			case items[0] == "Device type":
				sc.Devicetype_ = item1
			case items[0] == "Transport protocol":
				sc.Transportprotocol_ = item1
			case items[0] == "Local Time is": // sc.Localtimeis_	= item1
			case items[0] == "SMART Health Status":
				sc.Smarthealthstatus_ = item1
			case items[0] == "Current Drive Temperature":
				sc.Currentdrivetemperature_ = item1
			case items[0] == "Drive Trip Temperature":
				sc.Drivetriptemperature_ = item1
			case items[0] == "Specified cycle count over device lifetime":
				sc.Specifiedcyclecountoverdevicelifetime_ = item1
			case items[0] == "Accumulated start-stop cycles":
				sc.Accumulatedstartstopcycles_ = item1
			case items[0] == "Elements in grown defect list":
				sc.Elementsingrowndefectlist_ = item1
			case items[0] == "Non-medium error count":
				sc.Nonmediumerrorcount_ = item1
			case items[0] == "read":
				sc.Errorsread_ = item1
			case items[0] == "write":
				sc.Errorswrite_ = item1
			default:
				if _verbose {
					fmt.Printf("line%d: %s\n", ii, strings.Join(items, "!"))
				}
			}
		}
	}

	if sc.Smarthealthstatus_ == "FAILED" {
		// Probably should identify the failed physical drives, and blink their LED
	}
	return sc
}

// DoListSmartctldataOneSupermicro extracts smartctl for supermicro
func DoListSmartctldataOneSupermicro(_df *df.Dfdata, _scsi *scsi.Scsidata, _parted *parted.Parteddata, _dmidecode *dmidecode.Dmidecodedata, _verbose bool) *Smartctldata {
	sc := new(Smartctldata)

	if true {
		cmdstr := fmt.Sprintf("/usr/bin/timeout 10 /usr/sbin/smartctl -a %s", _scsi.Generic_)
		if _verbose {
			fmt.Println("cmdstr is ", cmdstr)
		}
		out := genutil.BashExecOrDie(_verbose, cmdstr, ".")
		if _verbose {
			fmt.Println(out)
		}
		lines := genutil.CleanAndSplitOnSpaces(out, " ")
		for ii, lineraw := range lines {
			line := strings.Replace(strings.TrimSpace(lineraw), ",", ";", -1)
			if _verbose {
				fmt.Printf("LINE%d=%s\n", ii, line)
			}
			if strings.Contains(line, "mandatory SMART command failed:") {
				break
			}
			items := strings.Split(line, ":")
			if len(items) < 2 {
				continue
			}
			item1 := strings.TrimSpace(items[1])
			switch {
			case items[0] == "Vendor":
				sc.Vendor_ = item1
			case items[0] == "Device Model":
				sc.Product_ = item1
			case items[0] == "Firmware Version":
				sc.Revision_ = item1
			case items[0] == "User Capacity":
				sc.Usercapacity_ = item1
			case items[0] == "Logical block size":
				sc.Logicalblocksize_ = item1
			case items[0] == "Logical Unit id":
				sc.Logicalunitid_ = item1
			case items[0] == "Serial Number":
				sc.Serialnumber_ = item1
			case items[0] == "ATA Standard is":
				sc.Devicetype_ = item1
			case items[0] == "Transport protocol":
				sc.Transportprotocol_ = item1
			case items[0] == "Local Time is": // sc.Localtimeis_	= item1
			case items[0] == "SMART overall-health self-assessment test result":
				sc.Smarthealthstatus_ = item1
			case items[0] == "Current Drive Temperature":
				sc.Currentdrivetemperature_ = item1
			case items[0] == "Drive Trip Temperature":
				sc.Drivetriptemperature_ = item1
			case items[0] == "Specified cycle count over device lifetime":
				sc.Specifiedcyclecountoverdevicelifetime_ = item1
			case items[0] == "Accumulated start-stop cycles":
				sc.Accumulatedstartstopcycles_ = item1
			case items[0] == "Elements in grown defect list":
				sc.Elementsingrowndefectlist_ = item1
			case items[0] == "Non-medium error count":
				sc.Nonmediumerrorcount_ = item1
			case items[0] == "read":
				sc.Errorsread_ = item1
			case items[0] == "write":
				sc.Errorswrite_ = item1
			default:
				if _verbose {
					fmt.Printf("line%d: %s\n", ii, strings.Join(items, "!"))
				}
			}
		}
	}
	return sc
}
