// Package df extracts useful info from that command on linux
//
// Csv output is particularly supported, so that a csvfile-based enterprise's ETL tools can also monitor its servers and desktops
package df

import (
	"fmt"
	"github.com/LDCS/genutil"
	"sort"
	"strings"
)

// Dfdata holds df data
type Dfdata struct {
	Name_       string // e.g. /dev/md1p1
	DevName_    string // e.g. /dev/md1 (inferred)
	Type_       string // e.g. local, network, etc
	Mountpoint_ string // e.g. /
	Sizegb_     string // e.g. 25G
	Usedgb_     string // e.g. 10G
	Availgb_    string // e.g. 15G
	Usepct_     string // e.g. 43%
}

const (
	names     = "name,type,mountpoint,sizegb,usedgb,availgb,usepct"
	hdrprefix = ",df."
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

// SortedKeys_String2PtrDfdata is generic
func SortedKeys_String2PtrDfdata(_mp *map[string][]*Dfdata) []string {
	keys := make([]string, len(*_mp))
	ii := 0
	for kk := range *_mp {
		keys[ii] = kk
		ii++
	}
	sort.Strings(keys)
	return keys
}

// Keys_String2PtrDfdata is generic
func Keys_String2PtrDfdata(_mp *map[string][]*Dfdata) []string {
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
func (self *Dfdata) Csv() string {
	if self == nil {
		return commaString
	}
	return fmt.Sprintf(pctString, self.Name_, self.Type_, self.Mountpoint_, self.Sizegb_, self.Usedgb_, self.Availgb_, self.Usepct_)
}

// Sprint is generic
func (self *Dfdata) Sprint() string {
	if self == nil {
		return ""
	}
	return fmt.Sprintf(namePctString, self.Name_, self.Type_, self.Mountpoint_, self.Sizegb_, self.Usedgb_, self.Availgb_, self.Usepct_)
}

// Print is generic
func (self *Dfdata) Print() {
	if self == nil {
		return
	}
	fmt.Printf(fmt.Sprint())
}

// New is generic
func New() *Dfdata { return new(Dfdata) }

// inferDevName is generic
func inferDevName(_name string) string {
	switch {
	case strings.Contains(_name, "/md"):
		parts := strings.Split(_name, "/md")
		if len(parts) != 2 {
			return _name
		}
		pre, post := parts[0], parts[1]
		partitionBeg := strings.Index(post, "p")
		if partitionBeg < 1 {
			return _name
		}
		return pre + "/md" + post[:partitionBeg]
	case strings.Contains(_name, "/sd"):
		parts := strings.Split(_name, "/sd")
		if len(parts) != 2 {
			return _name
		}
		pre, post := parts[0], parts[1]
		partitionBeg := strings.IndexAny(post, "0123456789")
		if partitionBeg < 1 {
			return _name
		}
		return pre + "/sd" + post[:partitionBeg]
	}
	return _name
}

// Df collects df data
func Df(_localOnly, _verbose bool) (smap map[string][]*Dfdata) {
	smap = make(map[string][]*Dfdata)
	localmap := map[string]bool{}
	out := genutil.BashExecOrDie(_verbose, genutil.StrTernary(_localOnly, "/usr/bin/timeout 10 /bin/df -klPT", "/bin/df -klPT; /usr/bin/timeout 10 /bin/df -kPT"), ".")
	if _verbose {
		fmt.Println(out)
	}
	lines := genutil.CleanAndSplitOnSpaces(out, ",")
	var lastdfd *Dfdata
	passNo := -1
	for ii, line := range lines {
		items := strings.Split(line, ",")
		if len(items) < 1 {
			continue
		}
		switch {
		case items[0] == "Filesystem":
			passNo++
			continue
		case len(items) == 7:
			lastdfd = new(Dfdata)
			lastdfd.Name_ = items[0]
			lastdfd.Type_ = items[1]
			lastdfd.Sizegb_ = genutil.KB2GB(items[2])
			lastdfd.Usedgb_ = genutil.KB2GB(items[3])
			lastdfd.Availgb_ = genutil.KB2GB(items[4])
			lastdfd.Usepct_ = strings.Replace(items[5], "%", "", -1)
			lastdfd.Mountpoint_ = strings.Join(items[6:], " ")
			lastdfd.DevName_ = inferDevName(lastdfd.Name_)
			if passNo == 0 {
				localmap[lastdfd.DevName_] = true                                // save the local filesystem names during the local pass
				smap[lastdfd.DevName_] = append(smap[lastdfd.DevName_], lastdfd) // was	smap[lastdfd.Name_]	= lastdfd
				if _verbose {
					fmt.Printf("line%d: item0(%s) %s\n", ii, items[0], strings.Join(items, "#"))
				}
			} else if !localmap[lastdfd.DevName_] { // skip local filesystem names during the nonlocal pass
				smap[lastdfd.DevName_] = append(smap[lastdfd.DevName_], lastdfd) // was	smap[lastdfd.Name_]	= lastdfd
				if _verbose {
					fmt.Printf("line%d: item0(%s) %s\n", ii, items[0], strings.Join(items, "#"))
				}
			}
		default:
			fmt.Printf("line%d: lenitems=%d item0(%s) %s\n", ii, len(items), items[0], strings.Join(items, "#"))
		}
	}
	return smap
}
