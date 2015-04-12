// Package nmap extracts mainly up/down status of network devices, as seen from linux.
//
// Csv output is particularly supported, so that a csvfile-based enterprise's ETL tools can also monitor its servers and desktops
package nmap

import (
	"fmt"
	"github.com/LDCS/genutil"
	"sort"
	"strings"
)

// Nmapdata holds nmap data
type Nmapdata struct {
	Subname_    string // ld
	Subnet_     string
	Ip_         string // e.g. 10.10.1.250
	Hostname_   string // e.g. ldfoo
	Status_     string // e.g. up or down
	MacAddress_ string // e.g. 00:26:52:0D:16:C3
	MacName_    string // e.g. Cisco Systems
}

const (
	names     = "Subname,Subnet,Ip,Hostname,Status,MacAddress,MacName"
	hdrprefix = ",nm."
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

// SortedKeys_String2PtrNmapdata is generic
func SortedKeys_String2PtrNmapdata(_mp *map[string]*Nmapdata) []string {
	keys := make([]string, len(*_mp))
	ii := 0
	for kk := range *_mp {
		keys[ii] = kk
		ii++
	}
	sort.Strings(keys)
	return keys
}

// Keys_String2PtrNmapdata is generic
func Keys_String2PtrNmapdata(_mp *map[string]*Nmapdata) []string {
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
func (self *Nmapdata) Csv() string {
	if self == nil {
		return commaString
	}
	return fmt.Sprintf(pctString, self.Subname_, self.Subnet_, self.Ip_, self.Hostname_, self.Status_, self.MacAddress_, self.MacName_)
}

// Sprint is generic
func (self *Nmapdata) Sprint() string {
	if self == nil {
		return ""
	}
	return fmt.Sprintf(namePctString, self.Subname_, self.Subnet_, self.Ip_, self.Hostname_, self.Status_, self.MacAddress_, self.MacName_)
}

// Print is generic
func (self *Nmapdata) Print() {
	if self == nil {
		return
	}
	fmt.Printf(fmt.Sprint())
}

// Nmap extracts nmap data
func Nmap(_subnets map[string][]string, _verbose bool) (smap map[string]*Nmapdata) {
	smap = make(map[string]*Nmapdata)
	out := ""
	for subname, subnetinfo := range _subnets {
		if _verbose {
			fmt.Println("subname=", subname, "subnet=", subnetinfo[0], "nmap.options=", subnetinfo[1])
		}
		out += fmt.Sprintf("\nsubname %s\n", subname)
		out += fmt.Sprintf("\nsubnet %s\n", subnetinfo[0])
		out += genutil.BashExecOrDie(_verbose, fmt.Sprintf("/usr/bin/timeout 20 nmap -v %s %s %s;", subnetinfo[1], subnetinfo[2], subnetinfo[0]), ".") // args are : -sP -sT subnet
	}
	if _verbose {
		fmt.Println(out)
	}
	lines := genutil.CleanAndSplitOnSpaces(out, ",")
	var lastnm *Nmapdata
	subname := ""
	subnet := ""
	status := ""
	for ii, lineraw := range lines {
		line := strings.TrimSpace(lineraw)
		// if _verbose { fmt.Printf("LINE%d=%s\n", ii, line) }
		items := strings.Split(line, ",")
		num := len(items)
		switch {
		case num < 2:
			continue
		case items[0] == "subname":
			subname = items[1]
		case items[0] == "subnet":
			subnet = items[1]
		case num < 4:
			continue
		case (items[0] == "Nmap") && (items[1] == "scan") && (items[2] == "report") && (items[3] == "for"):
			if _verbose {
				fmt.Printf("line%d: %s\n", ii, strings.Join(items, ","))
			}
			lastnm = new(Nmapdata)
			lastnm.Subname_ = subname
			lastnm.Subnet_ = subnet
			// fmt.Printf("   lastnm name=%s status=%s\n", lastnm.Name_, lastnm.Status_)
			switch num {
			case 5:
				lastnm.Ip_ = items[4]
			case 7:
				lastnm.Ip_, lastnm.Hostname_, status = items[5], items[4], items[6]
				if (lastnm.Ip_ == "[host") && (status == "down]") {
					continue
				}
			default:
				lastnm.Ip_, lastnm.Hostname_ = items[5], items[4]
			}
			lastnm.Ip_ = genutil.ChompParens(lastnm.Ip_, true)
			smap[lastnm.Ip_] = lastnm
			if _verbose {
				fmt.Printf("nm: %s\n", lastnm.Csv())
			}
		case (items[0] == "Host") && (items[1] == "is"):
			lastnm.Status_ = items[2]
		case (items[0] == "MAC") && (items[1] == "Address:"):
			lastnm.MacAddress_ = items[2]
			lastnm.MacName_ = genutil.ChompParens(strings.Join(items[3:], " "), true)
		default:
			if _verbose {
				fmt.Printf("line%d: item0(%s) %s\n", ii, items[0], strings.Join(items, "~"))
			}
		}
	}
	return smap
}
