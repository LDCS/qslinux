// Package etcshadow extracts useful info from that file on linux
//
// Csv output is particularly supported, so that a csvfile-based enterprise's ETL tools can also monitor its servers and desktops
package etcshadow

import (
	"crypto/md5"
	"fmt"
	"github.com/LDCS/genutil"
	"io"
	"sort"
	"strings"
)

// Shadowdata holds etc shadow info
type Shadowdata struct {
	Shadowname_  string
	PwMd5sum_    string
	Nlastchange_ string
	Ncanchanges_ string
	Nmustchange_ string
	Nwarn_       string
	Nexpire_     string
	Nexpired_    string
	Nreserved_   string
}

const (
	names     = "Shadowname,PwMd5sum,Nlastchange,Ncanchanges,Nmustchange,Nwarn,Nexpire,Nexpired,Nreserved"
	hdrprefix = ",xsh."
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

// SortedKeys_String2PtrShadowdata is generic
func SortedKeys_String2PtrShadowdata(_mp *map[string]*Shadowdata) []string {
	keys := make([]string, len(*_mp))
	ii := 0
	for kk := range *_mp {
		keys[ii] = kk
		ii++
	}
	sort.Strings(keys)
	return keys
}

// Keys_String2PtrShadowdata is generic
func Keys_String2PtrShadowdata(_mp *map[string]*Shadowdata) []string {
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
func (self *Shadowdata) Csv() string {
	if self == nil {
		return commaString
	}
	return fmt.Sprintf(pctString, self.Shadowname_, self.PwMd5sum_, self.Nlastchange_, self.Ncanchanges_, self.Nmustchange_, self.Nwarn_, self.Nexpire_, self.Nexpired_, self.Nreserved_)
}

// Sprint is generic
func (self *Shadowdata) Sprint() string {
	if self == nil {
		return ""
	}
	return fmt.Sprintf(namePctString, self.Shadowname_, self.PwMd5sum_, self.Nlastchange_, self.Ncanchanges_, self.Nmustchange_, self.Nwarn_, self.Nexpire_, self.Nexpired_, self.Nreserved_)
}

// SprintShort is generic
func (self *Shadowdata) SprintShort() string {
	if self == nil {
		return ""
	}
	return fmt.Sprintf("Nlastchange=%s:Ncanchanges=%s:Nmustchange=%s:Nwarn=%s:Nexpire=%s:Nexpired=%s:Nreserved=%s",
		self.Nlastchange_, self.Ncanchanges_, self.Nmustchange_, self.Nwarn_, self.Nexpire_, self.Nexpired_, self.Nreserved_)
}

// Print is generic
func (self *Shadowdata) Print() {
	if self == nil {
		return
	}
	fmt.Printf(fmt.Sprint())
}

// DoListShadowdata extracts info from etc shadow
func DoListShadowdata(_verbose bool) (smap map[string]*Shadowdata) {
	smap = make(map[string]*Shadowdata)
	out := genutil.BashExecOrDie(_verbose, "/bin/cat /etc/shadow| sed -e 's/[ ]/_/g'", ".")
	if _verbose {
		fmt.Println(out)
	}
	lines := genutil.CleanAndSplitOnSeparator(out, ":", ",")
	for ii, line := range lines {
		items := strings.Split(line, ",")
		num := len(items)
		if _verbose {
			fmt.Printf("line%d: item0(%s) %s\n", ii, items[0], strings.Join(items, "#"))
		}
		if num < 6 {
			continue
		}
		lastxsh := new(Shadowdata)
		lastxsh.Shadowname_ = items[0]
		smap[lastxsh.Shadowname_] = lastxsh
		switch items[1] {
		case "", "*":
			lastxsh.PwMd5sum_ = items[1]
		default:
			h := md5.New()
			io.WriteString(h, items[1])
			lastxsh.PwMd5sum_ = fmt.Sprintf("%x", h.Sum(nil))
		}
		lastxsh.Nlastchange_ = items[2]
		lastxsh.Ncanchanges_ = items[3]
		lastxsh.Nmustchange_ = items[4]
		lastxsh.Nwarn_ = items[5]
		if num > 6 {
			lastxsh.Nexpire_ = items[6]
		}
		if num > 7 {
			lastxsh.Nexpired_ = items[7]
		}
		if num > 8 {
			lastxsh.Nreserved_ = items[8]
		}
		if _verbose {
			fmt.Printf("line%d: item0(%s) %s\n", ii, items[0], strings.Join(items, "#"))
		}
	}
	return smap
}
