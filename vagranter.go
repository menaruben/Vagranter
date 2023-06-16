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
	Vmbox           string   `toml:"vmbox"`
	ForwardedPorts  [][]int  `toml:"forwarded_ports"`
	InlineScript    []string `toml:"inline_script"`
	ShellScriptPath string   `toml:"script"`
	Content         string
}

func writeConfig(text string) {
	file, err := os.Create("Vagrantfile")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = file.WriteString(text)
	if err != nil {
		log.Fatal(err)
	}
}

func (v *Vagrantfile) build() {
	var vagrantFileLines []string
	vagrantFileLines = append(vagrantFileLines, "Vagrant.configure(\"2\") do |config|")

	if v.Vmbox != "" {
		vagrantFileLines = append(vagrantFileLines, fmt.Sprintf("\tconfig.vm.box= \"%s\"", v.Vmbox))
	} else {
		color.Yellow("[-] vm.box needs to be defined!")
	}

	for _, portPair := range v.ForwardedPorts {
		line := fmt.Sprintf("\tconfig.vm.network \"forwarded_port\", guest: %d, host: %d", portPair[0], portPair[1])
		vagrantFileLines = append(vagrantFileLines, line)
	}

	if len(v.InlineScript) != 0 {
		vagrantFileLines = append(vagrantFileLines, "\tconfig.vm.provision \"shell\", inline: <<-SHELL")
		for _, cmd := range v.InlineScript {
			vagrantFileLines = append(vagrantFileLines, fmt.Sprintf("\t\t%s", cmd))
		}
		vagrantFileLines = append(vagrantFileLines, "\tSHELL")
	}

	if v.ShellScriptPath != "" {
		vagrantFileLines = append(vagrantFileLines, fmt.Sprintf("\tconfig.vm.provision \"shell\", path: \"%s\"", v.ShellScriptPath))
	}

	vagrantFileLines = append(vagrantFileLines, "end")
	vagrantFileContent := strings.Join(vagrantFileLines, "\n")
	v.Content = vagrantFileContent

	writeConfig(vagrantFileContent)
}

func main() {
	color.Green("[+] reading vconfig.toml...")
	tomlData, err := ioutil.ReadFile("vconfig.toml")
	if err != nil {
		log.Fatal(err)
	}

	var vagrantFile Vagrantfile
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
