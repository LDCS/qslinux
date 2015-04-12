// Package parted extracts useful info about iscsi targets provided by the linux server on which it runs.
//
// Csv output is particularly supported, so that a csvfile-based enterprise's ETL tools can also monitor its servers and desktops
package tgtd

import (
	"fmt"
	"github.com/LDCS/genutil"
	"net"
	"sort"
	"strings"
)

// Tgtddata holds tgtd data
type Tgtddata struct {
	Name_         string
	Targetpath_   string
	Targetstate_  string
	Tid_          string
	Initiator_    string
	Nexus_        string
	Connection_   string
	Ipaddress_    string
	LUN_          string
	Luntype_      string
	LunScsiID_    string
	LunScsiSN_    string
	Lunsize_      string
	Lunblocksize_ string
	Lunonline_    string
	Lunreadonly_  string
	Lunstore_     string
	Lunpath_      string
	Lunflags_     string
	ACL_          string
}

const (
	names     = "Name,Targetpath,Targetstate,Tid,Initiator,Nexus,Connection,Ipaddress,LUN,Luntype,LunScsiID,LunScsiSN,Lunsize,Lunblocksize,Lunonline,Lunreadonly,Lunstore,Lunpath,Lunflags,ACL"
	hdrprefix = ",tg."
	semi      = ";"
)

var (
	headerString  string
	commaString   string
	pctString     string
	namePctString string
)

// init  is generic
func init() {
	headerString = (hdrprefix + strings.Join(strings.Split(names, ","), hdrprefix))[1:]
	commaString = strings.Repeat(",", strings.Count(headerString, ","))
	pctString = strings.Repeat(",%s", 1+strings.Count(headerString, ","))[1:]
	namePctString = strings.Replace(names, ",", "=%s ", -1) + "=%s\n"
}

// SortedKeys_String2PtrTgtddata is generic
func SortedKeys_String2PtrTgtddata(_mp *map[string]*Tgtddata) []string {
	keys := make([]string, len(*_mp))
	ii := 0
	for kk := range *_mp {
		keys[ii] = kk
		ii++
	}
	sort.Strings(keys)
	return keys
}

// Keys_String2PtrTgtddata is generic
func Keys_String2PtrTgtddata(_mp *map[string]*Tgtddata) []string {
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
func (self *Tgtddata) Csv() string {
	if self == nil {
		return commaString
	}
	return fmt.Sprintf(pctString, self.Name_, self.Targetpath_, self.Targetstate_, self.Tid_, self.Initiator_, self.Nexus_, self.Connection_, self.Ipaddress_, self.LUN_, self.Luntype_, self.LunScsiID_, self.LunScsiSN_, self.Lunsize_, self.Lunblocksize_, self.Lunonline_, self.Lunreadonly_, self.Lunstore_, self.Lunpath_, self.Lunflags_, self.ACL_)
}

// Sprint is generic
func (self *Tgtddata) Sprint() string {
	if self == nil {
		return ""
	}
	return fmt.Sprintf(namePctString, self.Name_, self.Targetpath_, self.Targetstate_, self.Tid_, self.Initiator_, self.Nexus_, self.Connection_, self.Ipaddress_, self.LUN_, self.Luntype_, self.LunScsiID_, self.LunScsiSN_, self.Lunsize_, self.Lunblocksize_, self.Lunonline_, self.Lunreadonly_, self.Lunstore_, self.Lunpath_, self.Lunflags_, self.ACL_)
}

// Print is generic
func (self *Tgtddata) Print() {
	if self == nil {
		return
	}
	fmt.Printf(fmt.Sprint())
}

// Tgtd extracts tgtd data
func Tgtd(_verbose bool) (smap map[string]*Tgtddata) {
	smap = make(map[string]*Tgtddata)
	out := genutil.BashExecOrDie(_verbose, "which tgtadm && /usr/bin/timeout 10 tgtadm --lld iscsi --op show --mode target", ".")
	if _verbose {
		fmt.Println(out)
	}
	lines := genutil.CleanAndSplitOnSpaces(out, ",")
	inSystem, inNexusInfo, inLUN, inACL := false, false, false, false
	var lasttg *Tgtddata
	for ii, lineraw := range lines {
		line := strings.TrimSpace(lineraw)
		if _verbose {
			fmt.Printf("LINE%d=%s\n", ii, line)
		}
		items := strings.Split(line, ",")
		if len(items) < 1 {
			continue
		}
		switch {
		case items[0] == "Target":
			inSystem, inNexusInfo, inLUN, inACL = false, false, false, false
			if _verbose {
				fmt.Printf("line%d: %s\n", ii, strings.Join(items, ","))
			}
			lasttg = new(Tgtddata)
			lasttg.Tid_ = genutil.ChompStr(items[1], ":")
			lasttg.Name_ = items[2]
		case (items[0] == "System") && (items[1] == "information"):
			inSystem, inNexusInfo, inLUN, inACL = true, false, false, false
			continue
		case (items[0] == "Account") && (items[1] == "information"):
			inSystem, inNexusInfo, inLUN, inACL = false, false, false, false
			continue
		case (items[0] == "ACL") && (items[1] == "information:"):
			inSystem, inNexusInfo, inLUN, inACL = false, false, false, true
			continue
		case items[0] == "Driver:":
			continue
		case items[0] == "State:":
			lasttg.Targetstate_ = items[1]
		case (items[0] == "I_T") && (items[1] == "nexus") && (items[2] == "information:"):
			inSystem, inNexusInfo, inLUN, inACL = false, true, false, false
			continue
		case (items[0] == "I_T") && (items[1] == "nexus:"):
			lasttg.Nexus_ = items[2]
		case items[0] == "Initiator:":
			lasttg.Initiator_ = items[1]
		case items[0] == "Connection:":
			lasttg.Connection_ += items[1] + semi
		case (items[0] == "IP") && (items[1] == "Address:"):
			lasttg.Ipaddress_ += items[2] + semi
		case (items[0] == "LUN") && (items[1] == "information:"):
			inSystem, inNexusInfo, inLUN, inACL = false, false, true, false
			continue
		case (items[0] == "LUN:"):
			lasttg.LUN_ += items[1] + semi
		case (items[0] == "Type:") && inLUN:
			lasttg.Luntype_ += items[1] + semi
		case (items[0] == "SCSI") && (items[1] == "ID:") && inLUN:
			lasttg.LunScsiID_ += strings.Join(items[2:], "=") + semi
		case (items[0] == "SCSI") && (items[1] == "SN:") && inLUN:
			lasttg.LunScsiSN_ += strings.Join(items[2:], "=") + semi
		case (items[0] == "Size:") && inLUN:
			lasttg.Lunsize_ += items[1] + " " + items[2] + semi
			lasttg.Lunblocksize_ += items[6] + semi
		case (items[0] == "Online:") && inLUN:
			lasttg.Lunonline_ += items[1] + semi
		case (items[0] == "Removable") && (items[1] == "media:") && inLUN:
			continue
		case (items[0] == "Prevent") && (items[1] == "removal:") && inLUN:
			continue
		case (items[0] == "Readonly:") && inLUN:
			lasttg.Lunreadonly_ += items[1] + semi
		case (items[0] == "Backing") && (items[1] == "store") && (items[2] == "type:") && inLUN:
			lasttg.Lunstore_ += items[3] + semi
		case (items[0] == "Backing") && (items[1] == "store") && (items[2] == "path:") && inLUN:
			lasttg.Lunpath_ += items[3] + semi
			if items[3] != "None" {
				lasttg.Targetpath_ = items[3]
				smap[lasttg.Targetpath_] = lasttg
			}
		case (items[0] == "Backing") && (items[1] == "store") && (items[2] == "flags:") && inLUN:
			lasttg.Lunflags_ += strings.Join(items[3:], "=") + semi
		case inACL && (net.ParseIP(items[0]) != nil):
			lasttg.ACL_ += items[0] + semi
		default:
			if _verbose {
				fmt.Printf("line%d: item0(%s) %s\n", ii, items[0], strings.Join(items, "~"))
			}
		}
	}
	if false {
		fmt.Println("inSystem=", inSystem)
	}
	if false {
		fmt.Println("inNexusInfo=", inNexusInfo)
	}
	for _, vv := range smap {
		vv.Connection_ = genutil.ChompStr(vv.Connection_, semi)
		vv.Ipaddress_ = genutil.ChompStr(vv.Ipaddress_, semi)
		vv.LUN_ = genutil.ChompStr(vv.LUN_, semi)
		vv.Luntype_ = genutil.ChompStr(vv.Luntype_, semi)
		vv.LunScsiID_ = genutil.ChompStr(vv.LunScsiID_, semi)
		vv.LunScsiSN_ = genutil.ChompStr(vv.LunScsiSN_, semi)
		vv.Lunsize_ = genutil.ChompStr(vv.Lunsize_, semi)
		vv.Lunblocksize_ = genutil.ChompStr(vv.Lunblocksize_, semi)
		vv.Lunonline_ = genutil.ChompStr(vv.Lunonline_, semi)
		vv.Lunreadonly_ = genutil.ChompStr(vv.Lunreadonly_, semi)
		vv.Lunstore_ = genutil.ChompStr(vv.Lunstore_, semi)
		vv.Lunpath_ = genutil.ChompStr(vv.Lunpath_, semi)
		vv.Lunflags_ = genutil.ChompStr(vv.Lunflags_, semi)
		vv.ACL_ = genutil.ChompStr(vv.ACL_, semi)
	}
	return smap
}
