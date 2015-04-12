// Package parted extracts useful info from the lsscsi on linux
//
// Csv output is particularly supported, so that a csvfile-based enterprise's ETL tools can also monitor its servers and desktops
package scsi

import (
	"fmt"
	"github.com/LDCS/genutil"
	"sort"
	"strings"
)

// Scsdidata holds scsi data
type Scsidata struct {
	Device_         string // e.g. /dev/sda
	Generic_        string // e.g. /dev/sg0
	Host_           string // e.g. scsi0
	Channel_        string // e.g. 0
	Target_         string // e.g. 01
	LUN_            string // e.g. 0
	Devicetype_     string // e.g. disk
	Vendor_         string // e.g. ATA
	Model_          string // e.g. ST4000DM000-1F21
	Revision_       string // e.g. CC52
	Device_blocked_ string // e.g. 0
	Iocounterbits_  string // e.g. 32
	Iodone_cnt_     string // e.g. 0x2816e7
	Ioerr_cnt_      string // e.g. 0x13a74
	Iorequest_cnt_  string // e.g. 0x284812
	Queue_depth_    string // e.g. 31
	Queue_type_     string // e.g. simple
	Scsi_level_     string // e.g. 6
	State_          string // e.g. running
	Timeout_        string // e.g. 30
	Type_           string // e.g. 0
	Transport_      string // e.g. ISCSI or sata
	Targetname_     string
}

const (
	names     = "Device,Generic,Host,Channel,Target,LUN,Devicetype,Vendor,Model,Revision,Device_blocked,Iocounterbits,Iodone_cnt,Ioerr_cnt,Iorequest_cnt,Queue_depth,Queue_type,Scsi_level,State,Timeout,Type,Transport,Targetname"
	hdrprefix = ",sd."
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

// SortedKeys_String2PtrScsidata is generic
func SortedKeys_String2PtrScsidata(_mp *map[string]*Scsidata) []string {
	keys := make([]string, len(*_mp))
	ii := 0
	for kk, _ := range *_mp {
		keys[ii] = kk
		ii += 1
	}
	sort.Strings(keys)
	return keys
}

// Keys_String2PtrScsidata is generic
func Keys_String2PtrScsidata(_mp *map[string]*Scsidata) []string {
	keys := make([]string, len(*_mp))
	ii := 0
	for kk, _ := range *_mp {
		keys[ii] = kk
		ii += 1
	}
	return keys
}

// Header is generic
func Header() string { return headerString }

// Csv is generic
func (self *Scsidata) Csv() string {
	if self == nil {
		return commaString
	}
	return fmt.Sprintf(pctString, self.Device_, self.Generic_, self.Host_, self.Channel_, self.Target_, self.LUN_, self.Devicetype_, self.Vendor_, self.Model_, self.Revision_, self.Device_blocked_, self.Iocounterbits_, self.Iodone_cnt_, self.Ioerr_cnt_, self.Iorequest_cnt_, self.Queue_depth_, self.Queue_type_, self.Scsi_level_, self.State_, self.Timeout_, self.Type_, self.Transport_, self.Targetname_)
}

// Sprint is generic
func (self *Scsidata) Sprint() string {
	if self == nil {
		return ""
	}
	return fmt.Sprintf(namePctString, self.Device_, self.Generic_, self.Host_, self.Channel_, self.Target_, self.LUN_, self.Devicetype_, self.Vendor_, self.Model_, self.Revision_, self.Device_blocked_, self.Iocounterbits_, self.Iodone_cnt_, self.Ioerr_cnt_, self.Iorequest_cnt_, self.Queue_depth_, self.Queue_type_, self.Scsi_level_, self.State_, self.Timeout_, self.Type_, self.Transport_, self.Targetname_)
}

// Print is generic
func (self *Scsidata) Print() {
	if self == nil {
		return
	}
	fmt.Printf(fmt.Sprint())
}

// Scsi extracts lsscsi data
func Scsi(_verbose bool) (smap map[string]*Scsidata) {
	smap = make(map[string]*Scsidata)
	out := genutil.BashExecOrDie(_verbose, "/usr/bin/timeout 10 /usr/bin/lsscsi -Lg; /usr/bin/timeout 10 /usr/bin/lsscsi -Lgt;", ".")
	if _verbose {
		fmt.Println(out)
	}
	lines := genutil.CleanAndSplitOnSpaces(out, ",")
	var lastsd *Scsidata
	for ii, lineraw := range lines {
		line := strings.TrimSpace(lineraw)
		// if _verbose { fmt.Printf("LINE%d=%s\n", ii, line) }
		items := strings.Split(line, ",")
		if len(items) < 1 {
			continue
		}
		switch {
		case strings.HasPrefix(items[0], "["):
			if _verbose {
				fmt.Printf("line%d: %s\n", ii, strings.Join(items, ","))
			}
			num := len(items)
			devicetype, device, generic := items[1], items[num-2], items[num-1]
			if (devicetype == "storage") && (device == "-") {
				device = "-" + generic
			}
			lastsd = smap[device]
			if lastsd == nil {
				lastsd = new(Scsidata)
				// fmt.Printf("   lastsd name=%s status=%s\n", lastsd.Name_, lastsd.Status_)
				hctl := ""
				hctl, lastsd.Devicetype_, lastsd.Vendor_ = items[0], items[1], items[2]
				lastsd.Model_ = strings.Join(items[3:num-3], "-") // Combines "Logical Volume" into single item
				lastsd.Revision_, lastsd.Device_, lastsd.Generic_ = items[num-3], items[num-2], items[num-1]
				// if _verbose { fmt.Printf("foo %s,%s,%s,%s,%s,%s,%s\n", hctl,lastsd.Devicetype_,lastsd.Vendor_,lastsd.Model_,lastsd.Revision_,lastsd.Device_,lastsd.Generic_) }
				lastsd.Host_, lastsd.Channel_, lastsd.Target_, lastsd.LUN_ = genutil.ColonSplit4(hctl[1:(len(hctl) - 1)])
				if (lastsd.Devicetype_ == "storage") && (lastsd.Device_ == "-") {
					lastsd.Device_ = "-" + lastsd.Generic_
				}
				smap[lastsd.Device_] = lastsd
				// if _verbose { fmt.Printf("sd: %s\n", lastsd.Csv()) }
			}
		case strings.HasPrefix(items[0], "device_blocked="):
			_, lastsd.Device_blocked_ = genutil.EqualsSplit2Trimmed(items[0])
		case strings.HasPrefix(items[0], "iocounterbits="):
			_, lastsd.Iocounterbits_ = genutil.EqualsSplit2Trimmed(items[0])
		case strings.HasPrefix(items[0], "iodone_cnt="):
			_, lastsd.Iodone_cnt_ = genutil.EqualsSplit2Trimmed(items[0])
		case strings.HasPrefix(items[0], "ioerr_cnt="):
			_, lastsd.Ioerr_cnt_ = genutil.EqualsSplit2Trimmed(items[0])
		case strings.HasPrefix(items[0], "iorequest_cnt="):
			_, lastsd.Iorequest_cnt_ = genutil.EqualsSplit2Trimmed(items[0])
		case strings.HasPrefix(items[0], "queue_depth="):
			_, lastsd.Queue_depth_ = genutil.EqualsSplit2Trimmed(items[0])
		case strings.HasPrefix(items[0], "queue_type="):
			_, lastsd.Queue_type_ = genutil.EqualsSplit2Trimmed(items[0])
		case strings.HasPrefix(items[0], "scsi_level="):
			_, lastsd.Scsi_level_ = genutil.EqualsSplit2Trimmed(items[0])
		case strings.HasPrefix(items[0], "state="):
			_, lastsd.State_ = genutil.EqualsSplit2Trimmed(items[0])
		case strings.HasPrefix(items[0], "timeout="):
			_, lastsd.Timeout_ = genutil.EqualsSplit2Trimmed(items[0])
		case strings.HasPrefix(items[0], "type="):
			_, lastsd.Type_ = genutil.EqualsSplit2Trimmed(items[0])
		case strings.HasPrefix(items[0], "transport="):
			_, lastsd.Transport_ = genutil.EqualsSplit2Trimmed(items[0])
		case strings.HasPrefix(items[0], "targetname="):
			_, lastsd.Targetname_ = genutil.EqualsSplit2Trimmed(items[0])
		default:
			if _verbose {
				fmt.Printf("line%d: item0(%s) %s\n", ii, items[0], strings.Join(items, "~"))
			}
		}
	}
	return smap
}
