// Package hp extracts useful info, using the hpaucli program, about all hp disk controllers on a linux system
//
// Csv output is particularly supported, so that a csvfile-based enterprise's ETL tools can also monitor its servers and desktops
package hp

import (
	"fmt"
	"github.com/LDCS/genutil"
	"sort"
	"strings"
)

// Hppd holds Hp physical disk data
type Hppd struct {
	Pdaddress_  string // physical disk address
	Pdstatus_   string
	Pdldnum_    string
	Pdsize_     string
	Pdspeed_    string
	Pdfirmware_ string
	Pdserial_   string
	Pdmodel_    string
	PdtempC_    string
	PdtempmaxC_ string
	Pdrate_     string
	Ld_         *Hpld
	Ctrl_       *Hpctrl
}

// Hpld holds Hp logical disk data
type Hpld struct {
	Pdas_     string // list of physical drives
	Ldnum_    string // logical drive number
	Ldstatus_ string
	Ldsize_   string
	Raid_     string
	Lddev_    string // Os device name
	Mountpts_ string
	Ctrl_     *Hpctrl
}

// Hpctrl holds hp controller data
type Hpctrl struct {
	Type_          string
	Slotnum_       string
	Ctrlserial_    string
	Ctrlstatus_    string
	Cachestatus_   string
	Batterystatus_ string
	Ldlist         []Hpld // list of all logical drives
	Pdlist_        []Hppd // list of all physical drives
}

// Hpdata holds all HPacucli data
type Hpdata struct {
	Ctrls_ []*Hpctrl
	Pdmap_ map[string]*Hppd
	Ldmap_ map[string]*Hpld
}

const (
	namesPd   = "Pdaddress,Pdstatus,Pdldnum,Pdsize,Pdspeed,Pdfirmware,Pdserial,Pdmodel,PdtempC,PdtempmaxC,Pdrate"
	namesLd   = "Pdas,Ldnum,Ldstatus,Ldsize,Raid,Lddev,Mountpts"
	namesCtrl = "Type,Slotnum,Ctrlserial,Ctrlstatus,Cachestatus,Batterystatus"
	hdrprefix = ",hp."
)

var (
	headerStringPd    string
	headerStringLd    string
	headerStringCtrl  string
	commaStringPd     string
	commaStringLd     string
	commaStringCtrl   string
	pctStringPd       string
	pctStringLd       string
	pctStringCtrl     string
	namePctStringPd   string
	namePctStringLd   string
	namePctStringCtrl string
)

// Nil is generic
func Nil() *Hpdata {
	return nil
}

// init  is generic
func init() {
	headerStringPd = (hdrprefix + strings.Join(strings.Split(namesPd, ","), hdrprefix))[1:]
	headerStringLd = (hdrprefix + strings.Join(strings.Split(namesLd, ","), hdrprefix))[1:]
	headerStringCtrl = (hdrprefix + strings.Join(strings.Split(namesCtrl, ","), hdrprefix))[1:]
	commaStringPd = strings.Repeat(",", strings.Count(headerStringPd, ","))
	commaStringLd = strings.Repeat(",", strings.Count(headerStringLd, ","))
	commaStringCtrl = strings.Repeat(",", strings.Count(headerStringCtrl, ","))
	pctStringPd = strings.Repeat(",%s", 1+strings.Count(headerStringPd, ","))[1:]
	pctStringLd = strings.Repeat(",%s", 1+strings.Count(headerStringLd, ","))[1:]
	pctStringCtrl = strings.Repeat(",%s", 1+strings.Count(headerStringCtrl, ","))[1:]
	namePctStringPd = strings.Replace(namesPd, ",", "=%s ", -1) + "=%s\n"
	namePctStringLd = strings.Replace(namesLd, ",", "=%s ", -1) + "=%s\n"
	namePctStringCtrl = strings.Replace(namesCtrl, ",", "=%s ", -1) + "=%s\n"
}

// SortedKeys_String2PtrHpdata is generic
func SortedKeys_String2PtrHpdata(_mp *map[string]*Hpdata) []string {
	keys := make([]string, len(*_mp))
	ii := 0
	for kk, _ := range *_mp {
		keys[ii] = kk
		ii += 1
	}
	sort.Strings(keys)
	return keys
}

// Keys_String2PtrHpdata is generic
func Keys_String2PtrHpdata(_mp *map[string]*Hpdata) []string {
	keys := make([]string, len(*_mp))
	ii := 0
	for kk, _ := range *_mp {
		keys[ii] = kk
		ii += 1
	}
	return keys
}

// Header is generic
func Header() string { return headerStringCtrl + "," + headerStringLd + "," + headerStringPd }

// Csv is generic
func (self *Hppd) Csv() string {
	if self == nil {
		return commaStringPd
	}
	return fmt.Sprintf(pctStringPd, self.Pdaddress_, self.Pdstatus_, self.Pdldnum_, self.Pdsize_, self.Pdspeed_, self.Pdfirmware_, self.Pdserial_, self.Pdmodel_, self.PdtempC_, self.PdtempmaxC_, self.Pdrate_)
}

// Csv is generic
func (self *Hpld) Csv() string {
	if self == nil {
		return commaStringLd
	}
	return fmt.Sprintf(pctStringLd, self.Pdas_, self.Ldnum_, self.Ldstatus_, self.Ldsize_, self.Raid_, self.Lddev_, self.Mountpts_)
}

// Csv is generic
func (self *Hpctrl) Csv() string {
	if self == nil {
		return commaStringCtrl
	}
	return fmt.Sprintf(pctStringCtrl, self.Type_, self.Slotnum_, self.Ctrlserial_, self.Ctrlstatus_, self.Cachestatus_, self.Batterystatus_)
}

// SprintAll is generic
func (self *Hpdata) SprintAll(_box string) string {
	if self == nil {
		return ""
	}
	ostr := ""
	ostr += fmt.Sprintf("hp.box,%s\n", Header())
	doneld := map[string]bool{}
	donectrl := map[string]bool{}
	for _, pd := range self.Pdmap_ {
		ostr += fmt.Sprintf("%s,", _box)
		ostr += fmt.Sprint(pd.Ctrl_.Csv())
		ostr += fmt.Sprint(",")
		ostr += fmt.Sprint(pd.Ld_.Csv())
		ostr += fmt.Sprint(",")
		ostr += fmt.Sprint(pd.Csv())
		if len(pd.Pdldnum_) > 0 {
			doneld[fmt.Sprintf("%s:%s", pd.Ctrl_.Slotnum_, pd.Pdldnum_)] = true
		}
		if pd.Ctrl_ != nil {
			donectrl[pd.Ctrl_.Slotnum_] = true
		}
		ostr += "\n"
	}
	for _, ld := range self.Ldmap_ {
		if _, ok := doneld[fmt.Sprintf("%s:%s", ld.Ctrl_.Slotnum_, ld.Ldnum_)]; ok {
			continue
		}
		ostr += fmt.Sprintf("%s,", _box)
		ostr += fmt.Sprint(ld.Ctrl_.Csv())
		ostr += fmt.Sprint(",")
		ostr += fmt.Sprint(ld.Csv())
		ostr += fmt.Sprint(",")
		ostr += fmt.Sprint(((*Hppd)(nil)).Csv())
		if ld.Ctrl_ != nil {
			donectrl[ld.Ctrl_.Slotnum_] = true
		}
		ostr += "\n"
	}
	for _, hpc := range self.Ctrls_ {
		if _, ok := donectrl[hpc.Slotnum_]; ok {
			continue
		}
		ostr += fmt.Sprintf("%s,", _box)
		ostr += fmt.Sprint(hpc.Csv())
		ostr += fmt.Sprint(",")
		ostr += fmt.Sprint(((*Hpld)(nil)).Csv())
		ostr += fmt.Sprint(",")
		ostr += fmt.Sprint(((*Hppd)(nil)).Csv())
		ostr += "\n"
	}
	return ostr
}

// Hp extracts HPacucli data
func Hp(_verbose bool) *Hpdata {
	hpdata := new(Hpdata)
	hpdata.Pdmap_ = make(map[string]*Hppd)
	hpdata.Ldmap_ = make(map[string]*Hpld)
	out := genutil.BashExecOrDie(_verbose, "/usr/bin/timeout 20 /usr/sbin/hpacucli ctrl all show config detail", ".")
	if _verbose {
		fmt.Println(out)
	}
	lines := genutil.CleanAndSplitOnSpaces(out, ",")
	var lasthpc *Hpctrl
	var lasthpl *Hpld
	var lasthpp *Hppd
	mode, modenum := "notstarted", ""
	for ii, line := range lines {
		if _verbose {
			fmt.Printf("LINE%d: %s\n", ii, line)
		}
		items := genutil.Resplit(line, ",", " ", ":")
		if _verbose {
			fmt.Printf("line%d: item0(%s) %s\n", ii, items[0], strings.Join(items, "#"))
		}
		if len(items) < 1 {
			continue
		}
		if strings.HasPrefix(line, "Smart,Array,") {
			lasthpl = nil
			lasthpp = nil
			lasthpc = new(Hpctrl)
			hpdata.Ctrls_ = append(hpdata.Ctrls_, lasthpc)
			lasthpc.Type_ = strings.Split(line, ",")[2]
			mode, modenum = "controller", ""
			if _verbose {
				fmt.Printf("Mode %s %s\n", mode, modenum)
			}
			continue
		}
		if strings.HasPrefix(line, "SEP,") {
			mode, modenum = "finished", ""
			lasthpl = nil
			lasthpp = nil
			lasthpc = nil
			if _verbose {
				fmt.Printf("Mode %s %s\n", mode, modenum)
			}
			continue
		}
		if strings.HasPrefix(line, "Array:,") {
			mode, modenum = "logical", items[1]
			lasthpp = nil
			lasthpl = new(Hpld)
			lasthpl.Ctrl_ = lasthpc
			if _verbose {
				fmt.Printf("Mode %s %s\n", mode, modenum)
			}
			continue
		}
		if strings.HasPrefix(line, "unassigned") {
			mode, modenum = "controller", ""
			lasthpp = nil
			lasthpl = nil
			if _verbose {
				fmt.Printf("Mode %s %s\n", mode, modenum)
			}
			continue
		}
		if strings.HasPrefix(line, "physicaldrive,") && !strings.Contains(line, "(") { // Only catch long format "physicaldrive" lines
			lasthpp = new(Hppd)
			if lasthpl != nil {
				lasthpp.Pdldnum_ = lasthpl.Ldnum_
				lasthpp.Ld_ = lasthpl
			}
			lasthpp.Ctrl_ = lasthpc
			tmpitems := strings.Split(line, ",")
			mode, modenum = "physical", tmpitems[1]
			pda := fmt.Sprintf("%s:%s", lasthpc.Slotnum_, modenum)
			hpdata.Pdmap_[pda] = lasthpp
			if lasthpl != nil {
				lasthpl.Pdas_ += ";" + pda
			}
			if _verbose {
				fmt.Printf("Mode %s %s\n", mode, modenum)
			}
			continue
		}

		switch mode {
		case "finished":
			continue
		case "notstarted":
			continue
		case "controller":
			switch items[0] {
			case "Slot":
				lasthpc.Slotnum_ = items[1]
			case "Serial Number":
				lasthpc.Ctrlserial_ = items[1]
			case "Controller Status":
				lasthpc.Ctrlstatus_ = items[1]
			case "Cache Status":
				lasthpc.Cachestatus_ = items[1]
			case "Battery/Capacitor Status":
				lasthpc.Batterystatus_ = items[1]
			}
		case "logical":
			lasthpp = nil
			switch items[0] {
			case "Size":
				lasthpl.Ldsize_ = items[1]
			case "Logical Drive":
				lasthpl.Ldnum_ = items[1]
				hpdata.Ldmap_[fmt.Sprintf("%s:%s", lasthpc.Slotnum_, lasthpl.Ldnum_)] = lasthpl
			case "Fault Tolerance":
				lasthpl.Raid_ = items[1]
			case "Disk Name":
				lasthpl.Lddev_ = items[1]
			case "Mount Points":
				lasthpl.Mountpts_ = items[1]
			case "Status":
				lasthpl.Ldstatus_ = items[1]
			}
		case "physical":
			switch items[0] {
			case "Port":
				lasthpp.Pdaddress_ = modenum
			case "Status":
				lasthpp.Pdstatus_ = items[1]
			case "Size":
				lasthpp.Pdsize_ = items[1]
			case "Model":
				lasthpp.Pdmodel_ = items[1]
			case "Serial Number":
				lasthpp.Pdserial_ = items[1]
			case "Rotational Speed":
				lasthpp.Pdspeed_ = items[1]
			case "Current Temperature (C)":
				lasthpp.PdtempC_ = items[1]
			case "Maximum Temperature (C)":
				lasthpp.PdtempmaxC_ = items[1]
			case "Firmware Revision":
				lasthpp.Pdfirmware_ = items[1]
			case "PHY Transfer Rate":
				lasthpp.Pdrate_ = strings.Replace(items[1], "  Unknown", "", -1)
			default:
				continue
			}
		default:
			fmt.Printf("line%d: lenitems=%d item0(%s) %s\n", ii, len(items), items[0], strings.Join(items, "#"))
		}
	}
	return hpdata
}
