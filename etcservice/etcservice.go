// Package etcservices extracts useful info from the service command on linux
//
// Csv output is particularly supported, so that a csvfile-based enterprise's ETL tools can also monitor its servers and desktops
package etcservice

import (
	"fmt"
	"github.com/LDCS/genutil"
	"sort"
	"strings"
)

// Servicedata hold service data
type Servicedata struct {
	Service_ string
	Status_  string
	Pid_     string
}

const (
	names     = "Service,Status,Pid"
	hdrprefix = ",svc."
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

// SortedKeys_String2PtrServicedata is generic
func SortedKeys_String2PtrServicedata(_mp *map[string]*Servicedata) []string {
	keys := make([]string, len(*_mp))
	ii := 0
	for kk := range *_mp {
		keys[ii] = kk
		ii++
	}
	sort.Strings(keys)
	return keys
}

// Keys_String2PtrServicedata is generic
func Keys_String2PtrServicedata(_mp *map[string]*Servicedata) []string {
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
func (self *Servicedata) Csv() string {
	if self == nil {
		return commaString
	}
	return fmt.Sprintf(pctString, self.Service_, self.Status_, self.Pid_)
}

// Sprint is generic
func (self *Servicedata) Sprint() string {
	if self == nil {
		return ""
	}
	return fmt.Sprintf(namePctString, self.Service_, self.Status_, self.Pid_)
}

// Print is generic
func (self *Servicedata) Print() {
	if self == nil {
		return
	}
	fmt.Printf(fmt.Sprint())
}

// New is generic
func New() *Servicedata { return new(Servicedata) }

// Service holds chkconfig info
func Service(_verbose bool) (smap map[string]*Servicedata) {
	smap = make(map[string]*Servicedata)
	out := genutil.BashExecOrDie(_verbose, "/sbin/service --status-all | sed -s 's/,/|/g'", ".")

	statuslistA := [...]string{"stopped", "running"}
	statuslistB := [...]string{"running..."}

	if _verbose {
		fmt.Println(out)
	}
	lines := genutil.CleanAndSplitOnSpaces(out, ",")
	for ii, line := range lines {
		items := strings.Split(line, ",")
		num := len(items)
		switch {
		case num < 2:
			continue
		case items[0] == "#":
			continue
		case (num == 3) && (items[1] == "is") && genutil.SliceContainsStr(statuslistA[:], items[2]):
			lastsvc := new(Servicedata)
			lastsvc.Service_ = items[0]
			lastsvc.Status_ = items[2]
			smap[lastsvc.Service_] = lastsvc
		case (num == 4) && (items[0] == "Process") && (items[2] == "is") && (items[3] == "disabled"):
			lastsvc := new(Servicedata)
			lastsvc.Service_ = items[1]
			lastsvc.Status_ = items[3]
			smap[lastsvc.Service_] = lastsvc
		case (num == 6) && (items[2] == "enabled") && (items[3] == "using"):
			lastsvc := new(Servicedata)
			lastsvc.Service_ = items[0] + "-" + items[1]
			lastsvc.Status_ = "enabled"
			lastsvc.Pid_ = items[4] + "-" + items[5]
			smap[lastsvc.Service_] = lastsvc
		case (num == 4) && (items[1] == "is") && (items[2] == "not") && genutil.SliceContainsStr(statuslistA[:], items[3]):
			lastsvc := new(Servicedata)
			lastsvc.Service_ = items[0]
			lastsvc.Status_ = "not-" + genutil.ChompStr(items[3], "...")
			smap[lastsvc.Service_] = lastsvc
		case (num == 4) && (items[1] == "service") && (items[2] == "not") && (items[3] == "started"):
			lastsvc := new(Servicedata)
			lastsvc.Service_ = genutil.ChompStr(items[0], ":")
			lastsvc.Status_ = "not-started"
			smap[lastsvc.Service_] = lastsvc
		case (num == 5) && (items[1] == "(pid") && genutil.StrIsInt(genutil.ChompStr(items[2], ")")) && (items[3] == "is") && genutil.SliceContainsStr(statuslistB[:], items[4]):
			lastsvc := new(Servicedata)
			lastsvc.Service_ = items[0]
			lastsvc.Pid_ = genutil.ChompStr(items[2], ")")
			lastsvc.Status_ = genutil.ChompStr(items[4], "...")
			smap[lastsvc.Service_] = lastsvc
		case (num >= 5) && (items[1] == "(pid") && (items[num-2] == "is") && genutil.SliceContainsStr(statuslistB[:], items[num-1]):
			lastsvc := new(Servicedata)
			lastsvc.Service_ = items[0]
			lastsvc.Pid_ = genutil.ChompStr(strings.Join(items[2:(num-3)], ":"), ")")
			lastsvc.Status_ = genutil.ChompStr(items[num-1], "...")
			smap[lastsvc.Service_] = lastsvc
		default:
			if _verbose {
				fmt.Printf("line%d: item0(%s) %s\n", ii, items[0], strings.Join(items, "#"))
			}
		}
	}
	return smap
}
