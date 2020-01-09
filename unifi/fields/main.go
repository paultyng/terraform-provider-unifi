package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/iancoleman/strcase"
)

type replacement struct {
	Old string
	New string
}

var fieldReps = []replacement{
	{"Ipv6", "IPV6"},
	{"Pd", "PD"},
	{"Leasetime", "LeaseTime"},
	{"Networkgroup", "NetworkGroup"},
	{"Openvpn", "OpenVPN"},
	{"Nat", "NAT"},
	{"Wlangroup", "WLANGroup"},
	{"Dns", "DNS"},
	{"Tftp", "TFTP"},
	{"Dhcpdv6", "DHCPDV6"},
	{"Dhcpd", "DHCPD"},
	{"Dhcp", "DHCP"},
	{"Wlan", "WLAN"},
	{"Idx", "IDX"},
	{"Id", "ID"},
	{"Mac", "MAC"},
	{"Vlan", "VLAN"},
	{"Wep", "WEP"},
	{"Dtim", "DTIM"},
	{"Wpa", "WPA"},
	{"Usergroup", "UserGroup"},
	{"Mcastenhance", "MulticastEnhance"},
	{"Ssid", "SSID"},
	{"Dpi", "DPI"},
	{"Bc", "Broadcast"},
	{"Radius", "RADIUS"},
	{"Fixedip", "FixedIP"},
	{"Ipsec", "IPSec"},
	{"Ip", "IP"},
	{"Pmf", "PMF"},
	{"Minrssi", "MinRSSI"},
	{"Qos", "QOS"},
	{"Wan", "WAN"},
	{"Vpn", "VPN"},
	{"Igmp", "IGMP"},
	{"Tcp", "TCP"},
	{"Udp", "UDP"},
}

var fileReps = []replacement{
	{"WlanConf", "WLAN"},
	{"Dhcp", "DHCP"},
	{"Wlan", "WLAN"},
	{"NetworkConf", "Network"},
}

func cleanName(name string, reps []replacement) string {
	for _, rep := range reps {
		name = strings.ReplaceAll(name, rep.Old, rep.New)
	}

	return name
}

func main() {
	version := os.Args[1]
	out := os.Args[2]
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	fieldsDir := filepath.Join(wd, version)
	outDir := filepath.Join(wd, out)

	fieldsFiles, err := ioutil.ReadDir(fieldsDir)
	if err != nil {
		panic(err)
	}

	for _, fieldsFile := range fieldsFiles {
		name := fieldsFile.Name()

		// if name != "WlanConf.json" {
		// 	continue
		// }

		ext := filepath.Ext(name)

		if filepath.Ext(name) != ".json" {
			continue
		}

		if name == "Setting.json" {
			continue
		}

		name = name[:len(name)-len(ext)]

		structName := cleanName(name, fileReps)

		goFile := strcase.ToSnake(structName) + ".generated.go"
		code, err := generateCode(filepath.Join(fieldsDir, fieldsFile.Name()), structName)
		if err != nil {
			fmt.Printf("skipping file %s: %s", fieldsFile.Name(), err)
			continue
			// panic(err)
		}

		_ = os.Remove(filepath.Join(outDir, goFile))
		ioutil.WriteFile(filepath.Join(outDir, goFile), ([]byte)(code), 0644)
	}

	fmt.Printf("%s\n", outDir)
}

func generateCode(fieldsFile string, structName string) (string, error) {
	b, err := ioutil.ReadFile(fieldsFile)
	if err != nil {
		return "", err
	}

	var fields map[string]interface{}
	err = json.Unmarshal(b, &fields)
	if err != nil {
		return "", err
	}

	code := fmt.Sprintf(`// Code generated from ace.jar fields *.json files
// DO NOT EDIT.

package unifi

type %s struct {
	ID     string `+"`json:\"_id,omitempty\"`"+`
	SiteID string `+"`json:\"site_id,omitempty\"`"+`

	Hidden   bool   `+"`json:\"attr_hidden,omitempty\"`"+`
	HiddenID string `+"`json:\"attr_hidden_id,omitempty\"`"+`
	NoDelete bool   `+"`json:\"attr_no_delete,omitempty\"`"+`
	NoEdit   bool   `+"`json:\"attr_no_edit,omitempty\"`"+`

`, structName)

	fieldNames := []string{}
	for name := range fields {
		fieldNames = append(fieldNames, name)
	}

	// TODO: sort by normalized name, not this name
	sort.Strings(fieldNames)

	for _, name := range fieldNames {
		switch {
		case structName == "User" && name == "blocked":
			code += "\tBlocked bool `json:\"blocked,omitempty\"`\n"
			continue
		case structName == "User" && name == "last_seen":
			code += "\tLastSeen int `json:\"last_seen,omitempty\"`\n"
			continue
		case structName == "SettingUsg" && strings.HasSuffix(name, "_timeout"):
			field := strcase.ToCamel(name)
			field = cleanName(field, fieldReps)
			code += fmt.Sprintf("\t%s int `json:\"%s,omitempty\"`\n", field, name)
			continue
		}

		validation := fields[name]
		fieldCode, err := generateField(name, validation)
		if err != nil {
			return "", err
		}
		code += fieldCode + "\n"
	}

	code = code + "}\n"

	return code, nil
}

func normalizeValidation(re string) string {
	re = strings.ReplaceAll(re, "\\d", "[0-9]")
	re = strings.ReplaceAll(re, "[-+]?", "")
	re = strings.ReplaceAll(re, "[+-]?", "")
	re = strings.ReplaceAll(re, "[-]?", "")
	re = strings.ReplaceAll(re, "\\.", ".")
	re = strings.ReplaceAll(re, "[.]?", ".")

	quants := regexp.MustCompile(`\{\d*,?\d*\}|\*|\+|\?`)
	re = quants.ReplaceAllString(re, "")

	control := regexp.MustCompile(`[\(\[\]\)\|\-\$\^]`)
	re = control.ReplaceAllString(re, "")

	re = strings.TrimPrefix(re, "^")
	re = strings.TrimSuffix(re, "$")

	return re
}

func typeFromValidation(validation interface{}) (string, string, bool, error) {
	switch validation := validation.(type) {
	case []interface{}:
		if len(validation) == 0 {
			return "[]string", "", false, nil
		}
		if len(validation) > 1 {
			return "", "", false, fmt.Errorf("unknown validation %#v", validation)
		}
		elementType, elementComment, _, err := typeFromValidation(validation[0])
		if err != nil {
			return "", "", false, err
		}
		return fmt.Sprintf("[]%s", elementType), elementComment, true, nil
	case string:
		comment := validation
		normalized := normalizeValidation(validation)
		allowEmpty := strings.HasSuffix(validation, "|^$") || strings.HasPrefix(validation, "^$|")
		switch {
		case normalized == "falsetrue" || normalized == "truefalse":
			return "bool", "", false, nil
		default:
			if _, err := strconv.ParseFloat(normalized, 64); err == nil {
				if normalized == "09" || normalized == "09.09" {
					comment = ""
				}

				if strings.Contains(normalized, ".") {
					if strings.Contains(validation, "\\.){3}") {
						break
					}

					return "float64", comment, true, nil
				}

				return "int", comment, true, nil
			}
		}
		if validation != "" && normalized != "" {
			fmt.Printf("normalize %q to %q\n", validation, normalized)
		}
		return "string", validation, !allowEmpty, nil
	}
	return "", "", false, fmt.Errorf("unable to determine type from validation %q", validation)
}

func generateField(name string, validation interface{}) (string, error) {
	field := strcase.ToCamel(name)
	field = cleanName(field, fieldReps)
	fieldType, comment, omitempty, err := typeFromValidation(validation)
	if err != nil {
		return "", err
	}

	comment = strings.TrimSpace(fmt.Sprintf("// %s", comment))
	if comment == "//" {
		comment = ""
	}

	if fieldType == "string" && strings.HasSuffix(field, "ID") {
		omitempty = false
	}

	omitemptyCode := ""
	if omitempty {
		omitemptyCode = ",omitempty"
	}

	return fmt.Sprintf("\t%s %s `json:\"%s%s\"` %s", field, fieldType, name, omitemptyCode, comment), nil
}
