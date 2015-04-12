// Package parted extracts useful info from that command on linux
//
// Csv output is particularly supported, so that a csvfile-based enterprise's ETL tools can also monitor its servers and desktops
package parted

import (
	"fmt"
	"github.com/LDCS/genutil"
	"sort"
	"strings"
)

// Parteddata holds parted data
type Parteddata struct {
	// "path":"size":"transport-type":"logical-sector-size":"physical-sector-size":"partition-table-type":"model-name";
	Unit_ string // BYT
	Path_ string
	//Type_			string
	Devsize_            string // e.g. 0
	DevPath_            string // e.g. /dev/sda
	Transporttype_      string
	Logicalsectorsize_  string
	Physicalsectorsize_ string
	Partitiontabletype_ string
	Modelname_          string
	// "number":"begin":"end":"size":"filesystem-type":"partition-name":"flags-set";
	Partnumber_ string // Partition_number
	Partbegin_  string
	Partend_    string
	Partsize_   string
	Partfstype_ string
	Partname_   string
	Flagset_    string
	Skip_       bool
}

const (
	names     = "Unit,Path,Devsize,DevPath,Transporttype,Logicalsectorsize,Physicalsectorsize,Partitiontabletype,Modelname,Partnumber,Partbegin,Partend,Partsize,Partfstype,Partname,Flagset"
	hdrprefix = ",pd."
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

// SortedKeys_String2PtrParteddata is generic
func SortedKeys_String2PtrParteddata(_mp *map[string][]*Parteddata) []string {
	keys := make([]string, len(*_mp))
	ii := 0
	for kk := range *_mp {
		keys[ii] = kk
		ii++
	}
	sort.Strings(keys)
	return keys
}

// Keys_String2PtrParteddata is generic
func Keys_String2PtrParteddata(_mp *map[string][]*Parteddata) []string {
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
func (self *Parteddata) Csv() string {
	if self == nil {
		return commaString
	}
	return fmt.Sprintf(pctString, self.Unit_, self.Path_, self.Devsize_, self.DevPath_, self.Transporttype_, self.Logicalsectorsize_, self.Physicalsectorsize_, self.Partitiontabletype_, self.Modelname_, self.Partnumber_, self.Partbegin_, self.Partend_, self.Partsize_, self.Partfstype_, self.Partname_, self.Flagset_)
}

// Sprint is generic
func (self *Parteddata) Sprint() string {
	if self == nil {
		return ""
	}
	return fmt.Sprintf(namePctString, self.Unit_, self.Path_, self.Devsize_, self.DevPath_, self.Transporttype_, self.Logicalsectorsize_, self.Physicalsectorsize_, self.Partitiontabletype_, self.Modelname_, self.Partnumber_, self.Partbegin_, self.Partend_, self.Partsize_, self.Partfstype_, self.Partname_, self.Flagset_)
}

// Print is generic
func (self *Parteddata) Print() {
	if self == nil {
		return
	}
	fmt.Printf(fmt.Sprint())
}

// Parted extracts parted data
func Parted(_verbose bool) (smap map[string][]*Parteddata) {
	smap = make(map[string][]*Parteddata)
	out := genutil.BashExecOrDie(_verbose, "/usr/bin/timeout 10 /sbin/parted -lms", ".")
	if _verbose {
		fmt.Println(out)
	}
	lines := genutil.CleanAndSplitOnSpaces(out, ";")
	var lastpd *Parteddata
	for ii, lineraw := range lines {
		line := strings.Replace(strings.TrimSpace(lineraw), ",", ";", -1)
		if _verbose {
			fmt.Printf("LINE%d=%s\n", ii, line)
		}
		items := strings.Split(line, ":")
		if len(items) < 1 {
			continue
		}
		switch {
		case strings.HasPrefix(items[0], "BYT;"):
			fallthrough
		case strings.HasPrefix(items[0], "CYL"):
			fallthrough
		case strings.HasPrefix(items[0], "CHS;"):
			if _verbose {
				fmt.Printf("line%d: %s\n", ii, strings.Join(items, "!"))
			}
			lastpd = new(Parteddata)
			// lastpd.Type_	= "rawdevice"
			lastpd.Unit_ = items[0][:3]
			if _verbose {
				fmt.Printf("part X0 lastpd Type=%s Unit=%s\n", "obsolete" /*lastpd.Type_*/, lastpd.Unit_)
			}
			// smap[lastpd.Path_]	= lastpd
			// if _verbose { fmt.Printf("sd: %s\n", lastpd.Csv()) }
		case genutil.StrIsInt(items[0]):
			// Our setup :
			//	Centos 6.x going to Centos 7.x
			//	No LVM
			//	Some hardware raid (hp) but mostly software raid (mdadm)
			//	Application:
			//		(1) Mostly local storage
			//			Some of the storate is exported NFS (nfs4 exclusively) and rarely Samba
			//		(2) A few instances of mdadm disks/partitions exported iscsi, to other boxes that do not have sufficient local storage
			//	Naively mapping parted output to rows would result in rows for disks and partitions.
			//	The goal is to identify cases where partition rows can be collapsed into the parent disk rows.
			//	For example, this is the case if there is a single partition on the disk.
			//	Or if there is a single partition that is locally mounted (other partitions exported as iscsi will retain their rows in the output)
			if _verbose {
				fmt.Printf("part X1 line%d: %s\n", ii, strings.Join(items, "!"))
			}
			// "number":"begin":"end":"size":"filesystem-type":"partition-name":"flags-set";
			partpd := *lastpd
			partpd.Partnumber_, partpd.Partbegin_, partpd.Partend_, partpd.Partsize_, partpd.Partfstype_, partpd.Partname_, partpd.Flagset_ = items[0], items[1], items[2], items[3], items[4], items[5], items[6]
			partpd.Flagset_ = genutil.ShrinkSep(partpd.Flagset_, ';')
			partpd.DevPath_ = partpd.Path_
			switch {
			case strings.Contains(partpd.Path_, "/md"):
				if true {
					partpd.Path_ += "p" + partpd.Partnumber_
					// partpd.Type_	= "softpart"
					if _verbose {
						fmt.Printf("part X2keep? Type=%s path=%s num=%s beg=%s end=%s size=%s fstype=%s name=%s flagset=%s\n", "obsolete" /*partpd.Type_*/, partpd.Path_, partpd.Partnumber_, partpd.Partbegin_, partpd.Partend_, partpd.Partsize_, partpd.Partfstype_, partpd.Partname_, partpd.Flagset_)
					}
					smap[partpd.DevPath_] = append(smap[partpd.DevPath_], &partpd)
				} else {
					// Softraids have single partition so assign values to the rawdevice, and no need to add to smap
					// lastpd.Type_	= "softraid"
					lastpd.Partnumber_, lastpd.Partbegin_, lastpd.Partend_, lastpd.Partsize_, lastpd.Partfstype_, lastpd.Partname_, lastpd.Flagset_ = items[0], items[1], items[2], items[3], items[4], items[5], items[6]
					lastpd.Flagset_ = genutil.ShrinkSep(lastpd.Flagset_, ';')
					if _verbose {
						fmt.Printf("part X2skip Type=%s path=%s num=%s beg=%s end=%s size=%s fstype=%s name=%s flagset=%s\n", "obsolte" /*lastpd.Type_*/, lastpd.Path_, lastpd.Partnumber_, lastpd.Partbegin_, lastpd.Partend_, lastpd.Partsize_, lastpd.Partfstype_, lastpd.Partname_, lastpd.Flagset_)
					}
				}
			default:
				partpd.Path_ += partpd.Partnumber_
				// partpd.Type_	= "partition"
				if _verbose {
					fmt.Printf("part X2keep Type=%s path=%s num=%s beg=%s end=%s size=%s fstype=%s name=%s flagset=%s\n", "obsolete" /*partpd.Type_*/, partpd.Path_, partpd.Partnumber_, partpd.Partbegin_, partpd.Partend_, partpd.Partsize_, partpd.Partfstype_, partpd.Partname_, partpd.Flagset_)
				}
				smap[partpd.DevPath_] = append(smap[partpd.DevPath_], &partpd)

			}
		case strings.HasPrefix(items[0], "Error"):
			if _verbose {
				fmt.Printf("line%d: item0(%s) %s\n", ii, items[0], strings.Join(items, "~"))
			}
			if (len(items) > 1) && strings.Contains(items[1], "/dev/md") {
				lastpd = new(Parteddata)
				lastpd.Path_ = items[1][strings.Index(items[1], "/dev/md"):]
				lastpd.DevPath_ = lastpd.Path_
				// lastpd.Type_	= "error"
				smap[lastpd.DevPath_] = append(smap[lastpd.DevPath_], lastpd)
			}
		case strings.HasPrefix(items[0], "Warning"):
			if _verbose {
				fmt.Printf("line%d: item0(%s) %s\n", ii, items[0], strings.Join(items, "~"))
			}
		case len(items[0]) > 0:
			// if _verbose { fmt.Printf("full line%d: %s\n", ii, strings.Join(items, "!")) }
			lastpd.Path_, lastpd.Devsize_, lastpd.Transporttype_, lastpd.Logicalsectorsize_, lastpd.Physicalsectorsize_, lastpd.Partitiontabletype_, lastpd.Modelname_ = items[0], items[1], items[2], items[3], items[4], items[5], items[6]
			lastpd.DevPath_ = lastpd.Path_
			lastpd.Modelname_ = genutil.ShrinkSep(lastpd.Modelname_, ';')
			smap[lastpd.DevPath_] = append(smap[lastpd.DevPath_], lastpd)
		default:
			if _verbose {
				fmt.Printf("line%d: item0(%s) %s\n", ii, items[0], strings.Join(items, "~"))
			}
		}
	}
	for _, kk := range SortedKeys_String2PtrParteddata(&smap) {
		if len(smap[kk]) < 2 {
			continue
		}
		for _, row := range smap[kk] {
			if row.Path_ == row.DevPath_ {
				row.Skip_ = true
			}
			if _verbose {
				fmt.Println("KEYS devpath=", kk, "path=", row.Path_, "Skip=", row.Skip_)
			}
		}
	}
	return smap
}
