// Package md extracts useful info from the mdadm command on linux
//
// Csv output is particularly supported, so that a csvfile-based enterprise's ETL tools can also monitor its servers and desktops
package md

import (
	"fmt"
	"github.com/LDCS/genutil"
	"sort"
	"strings"
)

// Mddata holds mdadm data
type Mddata struct {
	Name_             string // e.g. md0
	Status_           string // e.g. active
	Raidtype_         string // e.g. raid1
	Level_            string // e.g. 5
	Members_          string // e.g. sda2[0]/sdb2[1]
	Blocks_           string // e.g. 511988
	Chunk_            string // e.g. 65536Kb
	Bitmap_           string // e.g. 1/1
	Nearcopies_       string // e.g. 2
	Pages_            string // e.g. [4KB]
	Superversion_     string // e.g. 1.2
	Algo_             string // e.g. 2
	Numcomponents_    string // e.g. [5/5]	First num is ideal, second is actual
	Componentstatus_  string // e.g. [U_]
	Checkpct_         string // 40.2
	Checkminutesleft_ string // 6137.7
	Resync_           string // DELAYED
}

// SortedKeys_String2PtrMddata is generic
func SortedKeys_String2PtrMddata(_mp *map[string]*Mddata) []string {
	keys := make([]string, len(*_mp))
	ii := 0
	for kk := range *_mp {
		keys[ii] = kk
		ii++
	}
	sort.Strings(keys)
	return keys
}

// Keys_String2PtrMddata is generic
func Keys_String2PtrMddata(_mp *map[string]*Mddata) []string {
	keys := make([]string, len(*_mp))
	ii := 0
	for kk := range *_mp {
		keys[ii] = kk
		ii++
	}
	return keys
}

// Header is generic
func Header() string {
	return fmt.Sprintf("md.name,md.status,md.Raidtype,md.Level,md.Members,md.Blocks,md.Chunk,md.Bitmap,md.Nearcopies,md.Pages,md.Superversion,md.Algo,md.Numcomponents,md.Componentstatus,md.Checkpct,md.Checkminutesleft,md.Resync")
}

// Csv is generic
func (self *Mddata) Csv() string {
	if self == nil {
		return ",,,,,,,,,,,,,,,,"
	}
	return fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s",
		self.Name_, self.Status_, self.Raidtype_, self.Level_, self.Members_, self.Blocks_, self.Chunk_, self.Bitmap_, self.Nearcopies_, self.Pages_, self.Superversion_, self.Algo_, self.Numcomponents_, self.Componentstatus_, self.Checkpct_, self.Checkminutesleft_, self.Resync_)
}

// Print is generic
func (self *Mddata) Print() {
	if self == nil {
		return
	}
	fmt.Printf("name=%s status=%s Raidtype=%s Level=%s Members=%s Blocks=%s Chunk=%s Bitmap=%s Nearcopies=%s Pages=%s Superversion=%s Algo=%s Numcomponents=%s Componentstatus=%s Checkpct=%s Checkminutesleft=%s Resync=%s\n",
		self.Name_, self.Status_, self.Raidtype_, self.Level_, self.Members_, self.Blocks_, self.Chunk_, self.Bitmap_, self.Nearcopies_, self.Pages_, self.Superversion_, self.Algo_, self.Numcomponents_, self.Componentstatus_, self.Checkpct_, self.Checkminutesleft_, self.Resync_)
}

// Sprint is generic
func (self *Mddata) Sprint() string {
	if self == nil {
		return ""
	}
	return fmt.Sprintf("name=%s status=%s Raidtype=%s Level=%s Members=%s Blocks=%s Chunk=%s Bitmap=%s Nearcopies=%s Pages=%s Superversion=%s Algo=%s Numcomponents=%s Componentstatus=%s Checkpct=%s Checkminutesleft=%s Resync=%s",
		self.Name_, self.Status_, self.Raidtype_, self.Level_, self.Members_, self.Blocks_, self.Chunk_, self.Bitmap_, self.Nearcopies_, self.Pages_, self.Superversion_, self.Algo_, self.Numcomponents_, self.Componentstatus_, self.Checkpct_, self.Checkminutesleft_, self.Resync_)
}

// DoListMddata extracts mdadm data
func DoListMddata(_verbose bool) (smap map[string]*Mddata) {
	smap = make(map[string]*Mddata)
	out := genutil.BashExecOrDie(_verbose, "/usr/bin/timeout 10 /bin/cat /proc/mdstat", ".")
	if _verbose {
		fmt.Println(out)
	}
	lines := genutil.CleanAndSplitOnSpaces(out, ",")
	var lastmdd *Mddata
	for ii, line := range lines {
		items := strings.Split(line, ",")
		if len(items) < 1 {
			continue
		}
		switch {
		case items[0] == "Personalities":
			continue
		case items[0] == "unused":
			continue
		case strings.HasPrefix(items[0], "md"):
			if _verbose {
				fmt.Printf("line%d: %s\n", ii, strings.Join(items, "/"))
			}
			lastmdd = new(Mddata)
			lastmdd.Name_ = "/dev/" + items[0]
			lastmdd.Status_ = items[2]
			lastmdd.Raidtype_ = items[3]
			lastmdd.Members_ = strings.Join(items[4:], "/")
			smap[lastmdd.Name_] = lastmdd
			// fmt.Printf("   lastmdd name=%s status=%s\n", lastmdd.Name_, lastmdd.Status_)
		case strings.HasPrefix(items[0], "bitmap") || ((len(items) > 1) && (items[1] == "blocks")):
			for jj := 0; jj < len(items); jj++ {
				switch items[jj] {
				case "pages":
					lastmdd.Pages_ = items[jj+1]
				case "near-copies":
					lastmdd.Nearcopies_ = items[jj-1]
				case "chunk", "chunks":
					lastmdd.Chunk_ = items[jj-1]
				case "blocks":
					lastmdd.Blocks_ = items[jj-1]
				case "bitmap:":
					lastmdd.Bitmap_ = items[jj+1]
				case "super":
					lastmdd.Superversion_ = items[jj+1]
				case "level":
					lastmdd.Level_ = items[jj+1]
				case "algorithm":
					lastmdd.Algo_ = items[jj-1]
				}
			}
			if len(items) > 2 {
				lastitem := items[len(items)-1]
				if strings.HasPrefix(lastitem, "[") && strings.HasSuffix(lastitem, "]") {
					lastmdd.Numcomponents_ = items[len(items)-2]
					lastmdd.Componentstatus_ = items[len(items)-1]
				}
			}
			if _verbose {
				fmt.Printf("line%d: item0(%s) %s\n", ii, items[0], strings.Join(items, "/"))
			}
		case (len(items) > 1) && strings.HasPrefix(items[1], "check"):
			for jj := 0; jj < len(items); jj++ {
				if (items[jj] == "check") && (len(items) >= jj+2) {
					lastmdd.Checkpct_ = items[jj+2]
				}
				if strings.HasPrefix(items[jj], "finish") {
					str2 := items[jj][7:]
					lastmdd.Checkminutesleft_ = str2[:(len(str2) - 3)]
				}
			}
			if _verbose {
				fmt.Printf("line%d: item0(%s) %s\n", ii, items[0], strings.Join(items, "/"))
			}
		case strings.HasPrefix(items[0], "resync"):
			lastmdd.Resync_ = items[0][7:]
			if _verbose {
				fmt.Printf("line%d: item0(%s) %s\n", ii, items[0], strings.Join(items, "/"))
			}
		default:
			if _verbose {
				fmt.Printf("line%d: item0(%s) %s\n", ii, items[0], strings.Join(items, "~"))
			}
		}
	}
	return smap
}
