package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/eiannone/keyboard"
	"github.com/fatih/color"
)

type Vagrantfile struct {
	Vmbox                  string     `toml:"box"`
	ForwardedPorts         [][]int    `toml:"forwarded_ports"`          // default: -
	InlineScript           []string   `toml:"inline_script"`            // default: -
	ShellScriptPath        string     `toml:"script"`                   // default: -
	BoxCheckUpdate         bool       `toml:"box_check_update"`         // default: true
	AllowFstabModification bool       `toml:"allow_fstab_modification"` // default: true
	AllowHostsModification bool       `toml:"allow_hosts_modification"` // default: true
	BootTimeout            int        `toml:"boot_timeout"`             // default: 300
	BoxDownloadInsecure    bool       `toml:"box_download_insecure"`    // default: false
	GracefulHaltTimeout    int        `toml:"graceful_halt_timeout"`    // default: 60
	Hostname               string     `toml:"hostname"`                 // default: -
	IgnoresBoxVagrantfile  bool       `toml:"ignores_box_vagrantfile"`  // default: false
	PublicNetwork          string     `toml:"public_network"`           // default: -
	PrivateNetwork         string     `toml:"private_network"`          // default: -
	PostUpMessage          string     `toml:"post_up_message"`          // default: -
	SyncedFolders          [][]string `toml:"synced_folder"`            // default: -
	UsablePortRange        [2]int     `toml:"usable_port_range"`        // default: 2200..2250

	VagrantFileLines []string
	Content          string
}

func (v *Vagrantfile) fillDefaults() {
	v.BoxCheckUpdate = true
	v.AllowFstabModification = true
	v.AllowHostsModification = true
	v.BootTimeout = 300
	v.BoxDownloadInsecure = false
	v.GracefulHaltTimeout = 60
	v.IgnoresBoxVagrantfile = false
	v.UsablePortRange = [2]int{2200, 2250}
}

func (v Vagrantfile) writeConfig() {
	file, err := os.Create("Vagrantfile")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = file.WriteString(v.Content)
	if err != nil {
		log.Fatal(err)
	}
}

func (v *Vagrantfile) addVmbox() {
	if v.Vmbox != "" {
		v.VagrantFileLines = append(v.VagrantFileLines, fmt.Sprintf("\tconfig.vm.box = \"%s\"", v.Vmbox))
	} else {
		color.Yellow("[-] box needs to be defined!")
	}
}

func (v *Vagrantfile) addPortForwarding() {
	for _, portPair := range v.ForwardedPorts {
		line := fmt.Sprintf("\tconfig.vm.network \"forwarded_port\", guest: %d, host: %d", portPair[0], portPair[1])
		v.VagrantFileLines = append(v.VagrantFileLines, line)
	}
}

func (v *Vagrantfile) addInlineScript() {
	if len(v.InlineScript) != 0 {
		v.VagrantFileLines = append(v.VagrantFileLines, "\tconfig.vm.provision \"shell\", inline: <<-SHELL")
		for _, cmd := range v.InlineScript {
			v.VagrantFileLines = append(v.VagrantFileLines, fmt.Sprintf("\t\t%s", cmd))
		}
		v.VagrantFileLines = append(v.VagrantFileLines, "\tSHELL")
	}
}

func (v *Vagrantfile) addShellScript() {
	if v.ShellScriptPath != "" {
		v.VagrantFileLines = append(v.VagrantFileLines, fmt.Sprintf("\tconfig.vm.provision \"shell\", path: \"%s\"", v.ShellScriptPath))
	}
}

func (v *Vagrantfile) addCheckUpdate() {
	if !v.BoxCheckUpdate {
		v.VagrantFileLines = append(v.VagrantFileLines, "\tconfig.vm.box_check_update = false")
	}
}

func (v *Vagrantfile) addAllowFstabModification() {
	if !v.AllowFstabModification {
		v.VagrantFileLines = append(v.VagrantFileLines, "\tconfig.vm.allow_fstab_modification = false")
	}
}

func (v *Vagrantfile) addAllowHostsModification() {
	if !v.AllowHostsModification {
		v.VagrantFileLines = append(v.VagrantFileLines, "\tconfig.vm.allow_hosts_modification = false")
	}
}

func (v *Vagrantfile) addBootTimeout() {
	if v.BootTimeout != 300 {
		v.VagrantFileLines = append(v.VagrantFileLines, fmt.Sprintf("\tconfig.vm.boot_timeout = %d", v.BootTimeout))
	}
}

func (v *Vagrantfile) addBoxDownloadInsecure() {
	if v.BoxDownloadInsecure {
		v.VagrantFileLines = append(v.VagrantFileLines, "\tconfig.vm.box_download_insecure = true")
	}
}

func (v *Vagrantfile) addGracefulHaltTimeout() {
	if v.GracefulHaltTimeout != 60 {
		v.VagrantFileLines = append(v.VagrantFileLines, fmt.Sprintf("\tconfig.vm.graceful_halt_timeout = %d", v.GracefulHaltTimeout))
	}
}

func (v *Vagrantfile) addHostname() {
	if v.Hostname != "" {
		v.VagrantFileLines = append(v.VagrantFileLines, fmt.Sprintf("\tconfig.vm.hostname = \"%s\"", v.Hostname))
	}
}

func (v *Vagrantfile) addIgnoresBoxVagrantfile() {
	if v.IgnoresBoxVagrantfile {
		v.VagrantFileLines = append(v.VagrantFileLines, "\tconfig.vm.ignores_box_vagrantfile = true")
	}
}

func (v *Vagrantfile) addPublicNetwork() {
	if v.PublicNetwork != "" {
		v.VagrantFileLines = append(v.VagrantFileLines, fmt.Sprintf("\tconfig.vm.network \"public_network\", ip: \"%s\"", v.PublicNetwork))
	}
}

func (v *Vagrantfile) addPrivateNetwork() {
	if v.PublicNetwork != "" {
		v.VagrantFileLines = append(v.VagrantFileLines, fmt.Sprintf("\tconfig.vm.network \"private_network\", ip: \"%s\"", v.PrivateNetwork))
	}
}

func (v *Vagrantfile) addPostUpMessage() {
	if v.PostUpMessage != "" {
		v.VagrantFileLines = append(v.VagrantFileLines, fmt.Sprintf("\tconfig.vm.post_up_message = \"%s\"", v.PostUpMessage))
	}
}

func (v *Vagrantfile) addSyncedFolder() {
	for _, folderPair := range v.SyncedFolders {
		v.VagrantFileLines = append(v.VagrantFileLines, fmt.Sprintf("\tconfig.vm.synced_folder \"%s\", \"%s\"", folderPair[0], folderPair[1]))
	}
}

func (v *Vagrantfile) addUsablePortRange() {
	if v.UsablePortRange != [2]int{2200, 2250} {
		v.VagrantFileLines = append(v.VagrantFileLines, fmt.Sprintf("\tconfig.vm.usable_port_range \"%s\"..\"%s\"", v.UsablePortRange[0], v.UsablePortRange[1]))
	}
}

func (v *Vagrantfile) build() {
	v.VagrantFileLines = append(v.VagrantFileLines, "Vagrant.configure(\"2\") do |config|")

	v.addVmbox()
	v.addPortForwarding()
	v.addInlineScript()
	v.addShellScript()
	v.addCheckUpdate()
	v.addAllowFstabModification()
	v.addAllowHostsModification()
	v.addBootTimeout()
	v.addBoxDownloadInsecure()
	v.addGracefulHaltTimeout()
	v.addHostname()
	v.addIgnoresBoxVagrantfile()
	v.addPublicNetwork()
	v.addPrivateNetwork()
	v.addPostUpMessage()
	v.addSyncedFolder()
	v.addUsablePortRange()

	v.VagrantFileLines = append(v.VagrantFileLines, "end")
	v.Content = strings.Join(v.VagrantFileLines, "\n")

	v.writeConfig()
}

func main() {
	color.Green("[+] reading vconfig.toml...")
	tomlData, err := ioutil.ReadFile("vconfig.toml")
	if err != nil {
		log.Fatal(err)
	}

	var vagrantFile Vagrantfile
	vagrantFile.fillDefaults()
	_, err = toml.Decode(string(tomlData), &vagrantFile)

	color.Green("[+] building Vagrantfile...")
	vagrantFile.build()

	err = keyboard.Open()
	if err != nil {
		log.Fatal(err)
	}

	color.Yellow("[?] Would you like to look at the Vagrantfile? [Y/N]: ")
	for {
		char, _, err := keyboard.GetKey()
		if err != nil {
			log.Fatal(err)
		}

		if char == 'y' || char == 'Y' {
			fmt.Printf("%s", vagrantFile.Content)
			break
		}

		if char == 'n' || char == 'N' {
			color.Green("\n[!] Okay, bye! :)")
			break
		}
	}
}
