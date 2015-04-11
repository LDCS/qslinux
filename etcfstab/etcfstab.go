// Package etcfstab extracts useful info from that file on linux
//
// Csv output is particularly supported, so that a csvfile-based enterprise's ETL tools can also monitor its servers and desktops
package etcfstab

import (
	"fmt"
	"github.com/LDCS/genutil"
	"sort"
	"strings"
)

// Fstabdata holds fstab data
type Fstabdata struct {
	Spec_    string
	File_    string
	Vfstype_ string
	Mntops_  string
	Freq_    string
	Passno_  string
}

const (
	names     = "Spec,File,Vfstype,Mntops,Freq,Passno"
	hdrprefix = ",xfs."
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

// SortedKeys_String2PtrFstabdata is generic
func SortedKeys_String2PtrFstabdata(_mp *map[string]*Fstabdata) []string {
	keys := make([]string, len(*_mp))
	ii := 0
	for kk := range *_mp {
		keys[ii] = kk
		ii++
	}
	sort.Strings(keys)
	return keys
}

// Keys_String2PtrFstabdata is generic
func Keys_String2PtrFstabdata(_mp *map[string]*Fstabdata) []string {
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
func (self *Fstabdata) Csv() string {
	if self == nil {
		return commaString
	}
	return fmt.Sprintf(pctString, self.Spec_, self.File_, self.Vfstype_, self.Mntops_, self.Freq_, self.Passno_)
}

// Sprint is generic
func (self *Fstabdata) Sprint() string {
	if self == nil {
		return ""
	}
	return fmt.Sprintf(namePctString, self.Spec_, self.File_, self.Vfstype_, self.Mntops_, self.Freq_, self.Passno_)
}

// Print is generic
func (self *Fstabdata) Print() {
	if self == nil {
		return
	}
	fmt.Printf(fmt.Sprint())
}

// New is generic
func New() *Fstabdata { return new(Fstabdata) }

// DoListFstabdata extracts fstab data
func DoListFstabdata(_verbose bool) (smap map[string]*Fstabdata) {
	smap = make(map[string]*Fstabdata)
	out := genutil.BashExecOrDie(_verbose, "/bin/cat /etc/fstab| sed -e 's/,/|/g'", ".")
	if _verbose {
		fmt.Println(out)
	}
	lines := genutil.CleanAndSplitOnSpaces(out, ",")
	for ii, line := range lines {
		items := strings.Split(line, ",")
		num := len(items)
		if num < 6 {
			continue
		}
		if items[0] == "#" {
			continue
		}
		lastxfs := new(Fstabdata)
		lastxfs.Spec_ = items[0]
		lastxfs.File_ = items[1]
		lastxfs.Vfstype_ = items[2]
		lastxfs.Mntops_ = items[3]
		lastxfs.Freq_ = items[4]
		lastxfs.Passno_ = items[5]
		smap[lastxfs.Spec_] = lastxfs
		if _verbose {
			fmt.Printf("line%d: item0(%s) %s\n", ii, items[0], strings.Join(items, "#"))
		}
	}
	return smap
}
