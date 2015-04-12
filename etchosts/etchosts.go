// Package etchosts extracts useful info from that file on linux
//
// Csv output is particularly supported, so that a csvfile-based enterprise's ETL tools can also monitor its servers and desktops
package etchosts

import (
	"fmt"
	"github.com/LDCS/genutil"
	"sort"
	"strings"
)

// Hostsdata holds hosts info
type Hostsdata struct {
	Ip_      string
	Name_    string
	Aliases_ string
}

const (
	names     = "Ip,Name,Aliases"
	hdrprefix = ",xho."
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

// SortedKeys_String2PtrHostsdata is generic
func SortedKeys_String2PtrHostsdata(_mp *map[string]*Hostsdata) []string {
	keys := make([]string, len(*_mp))
	ii := 0
	for kk := range *_mp {
		keys[ii] = kk
		ii++
	}
	sort.Strings(keys)
	return keys
}

// Keys_String2PtrHostsdata is generic
func Keys_String2PtrHostsdata(_mp *map[string]*Hostsdata) []string {
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
func (self *Hostsdata) Csv() string {
	if self == nil {
		return commaString
	}
	return fmt.Sprintf(pctString, self.Ip_, self.Name_, self.Aliases_)
}

// Sprint is generic
func (self *Hostsdata) Sprint() string {
	if self == nil {
		return ""
	}
	return fmt.Sprintf(namePctString, self.Ip_, self.Name_, self.Aliases_)
}

// Print is generic
func (self *Hostsdata) Print() {
	if self == nil {
		return
	}
	fmt.Printf(fmt.Sprint())
}

// New is generic
func New() *Hostsdata { return new(Hostsdata) }

// Hosts extracts etc hosts info
func Hosts(_verbose bool) (smap map[string]*Hostsdata) {
	smap = make(map[string]*Hostsdata)
	out := genutil.BashExecOrDie(_verbose, "/bin/cat /etc/hosts", ".")
	if _verbose {
		fmt.Println(out)
	}
	lines := genutil.CleanAndSplitOnSpaces(out, ",")
	for ii, line := range lines {
		items := strings.Split(line, ",")
		num := len(items)
		if num < 2 {
			continue
		}
		lastxho := new(Hostsdata)
		lastxho.Ip_ = items[0]
		lastxho.Name_ = items[1]
		if num > 2 {
			lastxho.Aliases_ = strings.Join(items[2:], ":")
		}
		smap[lastxho.Ip_] = lastxho
		if _verbose {
			fmt.Printf("line%d: item0(%s) %s\n", ii, items[0], strings.Join(items, "#"))
		}
	}
	return smap
}
