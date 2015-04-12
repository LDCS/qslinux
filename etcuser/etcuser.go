// Package etcuser extracts useful info from that file on linux
//
// Csv output is particularly supported, so that a csvfile-based enterprise's ETL tools can also monitor its servers and desktops
package etcuser

import (
	"fmt"
	"github.com/LDCS/genutil"
	"sort"
	"strings"
)

// Userdata holds etc passwd data
type Userdata struct {
	Username_   string
	Uid_        string
	Gid_        string
	Home_       string
	Groups_     string
	Shell_      string
	Pwinfo_     string
	Lastlogin_  string
	HomeFS_     string
	HomeUsedGB_ string
}

const (
	names     = "Username,Uid,Gid,Home,Groups,Shell,Pwinfo,Lastlogin,HomeFS,HomeUsedGB"
	hdrprefix = ",xus."
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

// SortedKeys_String2PtrUserdata is generic
func SortedKeys_String2PtrUserdata(_mp *map[string]*Userdata) []string {
	keys := make([]string, len(*_mp))
	ii := 0
	for kk := range *_mp {
		keys[ii] = kk
		ii++
	}
	sort.Strings(keys)
	return keys
}

// Keys_String2PtrUserdata is generic
func Keys_String2PtrUserdata(_mp *map[string]*Userdata) []string {
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
func (self *Userdata) Csv() string {
	if self == nil {
		return commaString
	}
	return fmt.Sprintf(pctString, self.Username_, self.Uid_, self.Gid_, self.Home_, self.Groups_, self.Shell_, self.Pwinfo_, self.Lastlogin_, self.HomeFS_, self.HomeUsedGB_)
}

// Sprint is generic
func (self *Userdata) Sprint() string {
	if self == nil {
		return ""
	}
	return fmt.Sprintf(namePctString, self.Username_, self.Uid_, self.Gid_, self.Home_, self.Groups_, self.Shell_, self.Pwinfo_, self.Lastlogin_, self.HomeFS_, self.HomeUsedGB_)
}

// Print is generic
func (self *Userdata) Print() {
	if self == nil {
		return
	}
	fmt.Printf(fmt.Sprint())
}

// User extracts user info from etc passwd
func User(_verbose bool) (smap map[string]*Userdata) {
	smap = make(map[string]*Userdata)
	out := genutil.BashExecOrDie(_verbose, "/bin/cat /etc/passwd | sed -e 's/[ ]/_/g'", ".")
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
		lastxus := new(Userdata)
		lastxus.Username_ = items[0]
		lastxus.Pwinfo_ = genutil.StrTernary(len(items[1]) < 2, items[1], fmt.Sprintf("len=%d", len(items[1])))
		lastxus.Uid_ = items[2]
		lastxus.Gid_ = items[3]
		lastxus.Home_ = items[5]
		lastxus.Groups_ = ""
		lastxus.Shell_ = items[6]
		lastxus.Lastlogin_ = ""
		lastxus.HomeFS_ = ""
		lastxus.HomeUsedGB_ = ""
		smap[lastxus.Username_] = lastxus
		if _verbose {
			fmt.Printf("line%d: item0(%s) %s\n", ii, items[0], strings.Join(items, "#"))
		}
	}
	return smap
}
