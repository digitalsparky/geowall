package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"strings"

	"github.com/coreos/go-iptables/iptables"
	"github.com/urfave/cli"
)

var V4URL string = "http://www.ipdeny.com/ipblocks/data/aggregated/%s-aggregated.zone"
var V6URL string = "http://www.ipdeny.com/ipv6/ipaddresses/aggregated/%s-aggregated.zone"

var fw *Firewall

type Firewall struct {
	Mode      string
	Table     string
	Chain     string
	Countries string
	IPTables  *iptables.IPTables
	IP6Tables *iptables.IPTables
	V4        bool
	V6        bool
}

// Setup the firewall
func (f *Firewall) Setup(chain string) {
	f.Table = "filter"
	f.Chain = "geowall"
	ip4t, err := iptables.NewWithProtocol(iptables.ProtocolIPv4)
	if err != nil {
		log.Fatal(err)
	}
	f.IPTables = ip4t
	ip6t, err := iptables.NewWithProtocol(iptables.ProtocolIPv6)
	if err != nil {
		log.Fatal(err)
	}
	f.IP6Tables = ip6t
}

// Add a V4 rule
func (f *Firewall) AddRuleV4(ipblock string) {
	if f.Mode == "allow" {
		f.check(f.IPTables.AppendUnique(f.Table, f.Chain, "-s", ipblock, "-j", "RETURN"))
	} else {
		f.check(f.IPTables.AppendUnique(f.Table, f.Chain, "-s", ipblock, "-j", "DROP"))
	}
}

// Add a V6 Rule
func (f *Firewall) AddRuleV6(ipblock string) {
	if f.Mode == "allow" {
		f.check(f.IP6Tables.AppendUnique(f.Table, f.Chain, "-s", ipblock, "-j", "RETURN"))
	} else {
		f.check(f.IP6Tables.AppendUnique(f.Table, f.Chain, "-s", ipblock, "-j", "DROP"))
	}
}

func (f *Firewall) check(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

// Initiate the firewall
func (f *Firewall) InitFirewall() {
	if f.V4 {
		f.check(f.IPTables.NewChain(f.Table, f.Chain))
		if a, _ := f.IPTables.Exists(f.Table, "INPUT", "-m", "state", "--state", "NEW", "-j", f.Chain); !a {
			f.check(f.IPTables.Insert(f.Table, "INPUT", 1, "-m", "state", "--state", "NEW", "-j", f.Chain))
		}
		if f.Mode == "allow" {
			f.check(f.IPTables.Append(f.Table, f.Chain, "-j", "DROP"))
		} else {
			f.check(f.IPTables.Append(f.Table, f.Chain, "-j", "RETURN"))
		}
	}
	if f.V6 {
		f.check(f.IP6Tables.NewChain(f.Table, f.Chain))
		if a, _ := f.IP6Tables.Exists(f.Table, "INPUT", "-m", "state", "--state", "NEW", "-j", f.Chain); !a {
			f.check(f.IP6Tables.Insert(f.Table, "INPUT", 1, "-m", "state", "--state", "NEW", "-j", f.Chain))
		}
		if f.Mode == "allow" {
			f.check(f.IP6Tables.Append(f.Table, f.Chain, "-j", "DROP"))
		} else {
			f.check(f.IP6Tables.Append(f.Table, f.Chain, "-j", "RETURN"))
		}
	}
}

// Clear the existing chains if they exist
func (f *Firewall) ClearFirewall() {
	if f.V4 {
		f.check(f.IPTables.ClearChain(f.Table, f.Chain))
	}
	if f.V6 {
		f.check(f.IP6Tables.ClearChain(f.Table, f.Chain))
	}
}

// Unload the firewall
func (f *Firewall) UnloadFirewall() {
	if f.V4 {
		// Remove the IPv4 pre-process rule
		fmt.Println("Unloading IPv4 Rules")
		if a, _ := f.IPTables.Exists(f.Table, "INPUT", "-m", "state", "--state", "NEW", "-j", f.Chain); a {
			f.check(f.IPTables.Delete(f.Table, "INPUT", "-m", "state", "--state", "NEW", "-j", f.Chain))
		}
	}

	if f.V6 {
		// Remove the IPv6 pre-process rule
		fmt.Println("Unloading IPv6 Rules")
		if a, _ := f.IP6Tables.Exists(f.Table, "INPUT", "-m", "state", "--state", "NEW", "-j", f.Chain); a {
			f.check(f.IP6Tables.Delete(f.Table, "INPUT", "-m", "state", "--state", "NEW", "-j", f.Chain))
		}
	}
	// Drop the chains
	f.ClearFirewall()

	if f.V4 {
		f.check(f.IPTables.DeleteChain(f.Table, f.Chain))
	}
	if f.V6 {
		f.check(f.IP6Tables.DeleteChain(f.Table, f.Chain))
	}
}

// Reset the firewall
func (f *Firewall) ResetFirewall() {
	f.UnloadFirewall()
	f.InitFirewall()
}

func (f *Firewall) getRules(url string) []string {
	var list []string
	countryArr := strings.Split(f.Countries, ",")
	for _, c := range countryArr {
		data := f.downloadFile(fmt.Sprintf(url, strings.ToLower(strings.TrimSpace(c))))
		dataArr := strings.Split(data, "\n")
		for _, d := range dataArr {
			dd := strings.TrimSpace(d)
			if dd == "" {
				continue
			}
			list = append(list, dd)
		}
	}

	return list
}

func (f *Firewall) GetV4Rules() []string {
	return f.getRules(V4URL)
}

func (f *Firewall) GetV6Rules() []string {
	return f.getRules(V6URL)
}

func (f *Firewall) downloadFile(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("bad status: %s", resp.Status)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return string(data)
}

func (f *Firewall) ProcessRules() {
	var rcount int
	var count int
	if f.V4 {
		fmt.Println("Processing IPv4 Rules")
		v4rules := f.GetV4Rules()
		rcount = len(v4rules)
		fmt.Printf("Got %d rules to add...\n", rcount)
		count = 0
		for _, rule := range v4rules {
			count = count + 1
			f.AddRuleV4(rule)
			fmt.Printf("%d of %d\r", count, rcount)
		}
		fmt.Printf("Imported %d rules\n", rcount)
	}

	if f.V6 {
		fmt.Println("Processing IPv6 Rules")
		v6rules := f.GetV6Rules()
		rcount = len(v6rules)
		count = 0
		for _, rule := range v6rules {
			count = count + 1
			f.AddRuleV6(rule)
			fmt.Printf("%d of %d\r", count, rcount)
		}
		fmt.Printf("Imported %d rules\n", rcount)
	}
}

func runCMD(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}

func which(binary string) string {
	response, err := runCMD("which", binary)
	if err != nil {
		log.Fatal(err)
	}
	return response
}

func init() {
	currUser, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	if currUser.Username != "root" {
		fmt.Println("This tool manipulates the firewall and must be run as root")
		os.Exit(1)
	}
	fw = new(Firewall)
	fw.Setup("geowall")
}

func main() {
	app := cli.NewApp()
	app.Name = "geowall"
	app.EnableBashCompletion = true
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Matt Spurrier",
			Email: "matthew@spurrier.com.au",
		},
	}
	app.Copyright = "(c) 2018 Matt Spurrier"
	app.Usage = "GeoIP Based Firewall"

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "v4, 4, ipv4",
			Usage:       "Include IPv4 Rules",
			Destination: &fw.V4,
		},
		cli.BoolFlag{
			Name:        "v6, 6, ipv6",
			Usage:       "Include IPv6 Rules",
			Destination: &fw.V6,
		},
	}

	app.Commands = []cli.Command{
		cli.Command{
			Name: "start",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "countries, c",
					Usage:       "Comma delimited list of countries the mode will run actions on",
					EnvVar:      "GEOWALL_COUNTRIES",
					Destination: &fw.Countries,
				},
				cli.StringFlag{
					Name:        "mode, m",
					Usage:       "Mode (allow or deny)",
					EnvVar:      "GEOWALL_MODE",
					Destination: &fw.Mode,
				},
			},
			Action: func(c *cli.Context) error {
				if fw.Countries == "" {
					return errors.New("countries must be listed")
				}
				if fw.Mode != "allow" && fw.Mode != "deny" {
					return errors.New("unrecognised mode, needs to be allow or deny")
				}
				if !fw.V4 && !fw.V6 {
					return errors.New("Both V4 and V6 disabled, nothing to do")
				}

				fmt.Println("Initiating Firewall")
				fw.InitFirewall()

				fw.ProcessRules()
				fmt.Println("Update complete")
				fmt.Println("Make sure you save your IPtables Rules with iptables-save if you wish to keep them.")
				return nil
			},
		},

		cli.Command{
			Name: "update",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "countries, c",
					Usage:       "Comma delimited list of countries the mode will run actions on",
					EnvVar:      "GEOWALL_COUNTRIES",
					Destination: &fw.Countries,
				},
			},
			Action: func(c *cli.Context) error {
				if fw.Countries == "" {
					return errors.New("countries must be listed")
				}
				fmt.Println("Clearing old rules")
				fw.ClearFirewall()
				fw.ProcessRules()
				fmt.Println("Update complete")
				fmt.Println("Make sure you save your IPtables Rules with iptables-save if you wish to keep them.")
				return nil
			},
		},

		cli.Command{
			Name: "stop",
			Action: func(c *cli.Context) error {
				if !fw.V4 && !fw.V6 {
					return errors.New("Both V4 and V6 disabled, nothing to do")
				}

				fw.UnloadFirewall()
				fmt.Println("Firewall Unloaded")
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
