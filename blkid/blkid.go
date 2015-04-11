// Package blkid extracts useful info from that command on linux
//
// Csv output is particularly supported, so that a csvfile-based enterprise's ETL tools can also monitor its servers and desktops
package blkid

import (
	"fmt"
	"github.com/LDCS/genutil"
	"sort"
	"strings"
)

// Blkiddata is a type
type Blkiddata struct {
	Devname_   string
	Uuid_      string
	Uuidsub_   string
	Type_      string
	Label_     string
	Parttype_  string
	Partuuid_  string
	Partlabel_ string
}

const (
	names     = "Devname,Uuid,Uuidsub,Type,Label,Parttype,Partuuid,Partlabel"
	hdrprefix = ",bi."
	semi      = ";"
)

var (
	headerString  string
	commaString   string
	pctString     string
	namePctString string
)

func init() {
	headerString = (hdrprefix + strings.Join(strings.Split(names, ","), hdrprefix))[1:]
	commaString = strings.Repeat(",", strings.Count(headerString, ","))
	pctString = strings.Repeat(",%s", 1+strings.Count(headerString, ","))[1:]
	namePctString = strings.Replace(names, ",", "=%s ", -1) + "=%s\n"
}

// SortedKeys_String2PtrBlkiddata is generic
func SortedKeys_String2PtrBlkiddata(_mp *map[string]*Blkiddata) []string {
	keys := make([]string, len(*_mp))
	ii := 0
	for kk := range *_mp {
		keys[ii] = kk
		ii++
	}
	sort.Strings(keys)
	return keys
}

// Keys_String2PtrBlkiddata is generic
func Keys_String2PtrBlkiddata(_mp *map[string]*Blkiddata) []string {
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
func (self *Blkiddata) Csv() string {
	if self == nil {
		return commaString
	}
	return fmt.Sprintf(pctString, self.Devname_, self.Uuid_, self.Uuidsub_, self.Type_, self.Label_, self.Parttype_, self.Partuuid_, self.Partlabel_)
}

// Sprint is generic
func (self *Blkiddata) Sprint() string {
	if self == nil {
		return ""
	}
	return fmt.Sprintf(namePctString, self.Devname_, self.Uuid_, self.Uuidsub_, self.Type_, self.Label_, self.Parttype_, self.Partuuid_, self.Partlabel_)
}

// Print is generic
func (self *Blkiddata) Print() {
	if self == nil {
		return
	}
	fmt.Printf(fmt.Sprint())
}

// New is generic
func New() *Blkiddata { return new(Blkiddata) }

// DoListBlkiddata obtains information from blkid
func DoListBlkiddata(_verbose bool) (smap map[string]*Blkiddata) {
	smap = make(map[string]*Blkiddata)
	out := genutil.BashExecOrDie(_verbose, "/usr/bin/timeout 10 /sbin/blkid | sed -e 's/\" /\"|/g' | sed -s 's/ /|/' | sed -e 's/ /_/g' | sed -e 's/|/ /g'", ".") // turn sep spaces into pipe, then embedded spaces into underscore, then revert sep pipe to sep space
	if _verbose {
		fmt.Println(out)
	}
	lines := genutil.CleanAndSplitOnSpaces(out, ",")
	var lastblkidd *Blkiddata
	for ii, line := range lines {
		items := strings.Split(line, ",")
		if len(items) < 1 {
			continue
		}
		switch {
		case len(items) > 1:
			lastblkidd = new(Blkiddata)
			lastblkidd.Devname_ = genutil.ChompStr(items[0], ":")
			for _, item := range items[1:] {
				kk, vv := genutil.EqualsSplit2Trimmed(item)
				switch kk {
				case "UUID":
					lastblkidd.Uuid_ = genutil.ChompQuotes(vv, true)
				case "UUID_SUB":
					lastblkidd.Uuidsub_ = genutil.ChompQuotes(vv, true)
				case "TYPE":
					lastblkidd.Type_ = genutil.ChompQuotes(vv, true)
				case "SEC_TYPE": // secondary type is not interesting
				case "LABEL":
					lastblkidd.Label_ = genutil.ChompQuotes(vv, true)
				case "PTTYPE":
					lastblkidd.Parttype_ = genutil.ChompQuotes(vv, true)
				case "PARTUUID":
					lastblkidd.Partuuid_ = genutil.ChompQuotes(vv, true)
				case "PARTLABEL":
					lastblkidd.Partlabel_ = genutil.ChompQuotes(vv, true)
				default:
					fmt.Printf("line%d: unknown key: kk=%s vv=%s\n", kk, vv)
				}
			}
			if len(lastblkidd.Devname_) < 1 {
				fmt.Printf("DoListBlkiddata: Skipping empty Devname: line%d:  %s\n", ii, strings.Join(items, "#"))
				continue
			}
			smap[lastblkidd.Devname_] = lastblkidd
			if _verbose {
				fmt.Printf("line%d: item0(%s) %s\n", ii, items[0], strings.Join(items, "#"))
			}
		default:
			fmt.Printf("line%d: lenitems=%d item0(%s) %s\n", ii, len(items), items[0], strings.Join(items, "#"))
		}
	}
	return smap
}
