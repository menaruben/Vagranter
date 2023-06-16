# Vagranter
Make your Vagrantfile simpler by using a toml file instead of ruby!

# Usage
Create a `vconfig.toml` file and type in your configuration (these are all the current features supported):
```toml
box = "centos/7"
forwarded_ports = [[8080, 8081], [22, 2221]]        # 0: guest, 1: host
inline_script = ["echo 'Hello World'", "whoami"]
script = "sample.sh"
box_check_update = true                     # default
allow_fstab_modification = true             # default
allow_hosts_modification = true             # default
boot_timeout = 300                          # default
box_download_insecure = false               # default
graceful_halt_timeout = 60                  # default
hostname = "vagranter"
ignores_box_vagrantfile = false             # default
public_network = "192.168.42.10"
private_network = "192.168.43.10"
post_up_message = "VM is now ready to use!"
synced_folder = [["./", "/var/www"]]        # multiple synced folders can be defined; 0: host, 1: guest
usable_port_range = [2200, 2250]            # default
```

Now you can run the following command in your terminal to build your Vagrantfile:
```
go run vagranter.go
```

You should get the following output:
```
[+] reading vconfig.toml...
[+] building Vagrantfile...
[?] Would you like to look at the Vagrantfile? [Y/N]:
```

You can enter 'y' to take a look at the Vagrantfile that has been created:
```ruby
Vagrant.configure("2") do |config|
        config.vm.box = "centos/7"
        config.vm.network "forwarded_port", guest: 8080, host: 8081
        config.vm.network "forwarded_port", guest: 22, host: 2221
        config.vm.provision "shell", inline: <<-SHELL
                echo 'Hello World'
                whoami
        SHELL
        config.vm.provision "shell", path: "sample.sh"
        config.vm.hostname = "vagranter"
        config.vm.network "public_network", ip: "192.168.42.10"
        config.vm.network "private_network", ip: "192.168.43.10"
        config.vm.post_up_message = "VM is now ready to use!"
        config.vm.synced_folder "./", "/var/www"
end
```

You can now run the virtual machine using `vagrant up`!

# future featues
I am currently trying to add as much as possible from the [config.vm](https://developer.hashicorp.com/vagrant/docs/vagrantfile/machine_settings) so that we can also start using toml
for bigger / more complicated Vagrantfile configurations.

Here are the currently supported configurations (need to contain the "toml:"...""):
```go
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
```
The networking part of vagrant is **not** finished **yet** though.
