# Vagranter
Make your Vagrantfile simpler by using a toml file instead of ruby!

# Usage
Create a `vconfig.toml` file and type in your configuration (these are all the current features supported):
```toml
vmbox = "centos/7"
forwarded_ports = [[8080, 8081], [22, 2221]]
inline_script = ["echo 'Hello World'", "whoami"]
script = "sample.sh"
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
  config.vm.box= "centos/7"
  config.vm.network "forwarded_port", guest: 8080, host: 8081
  config.vm.network "forwarded_port", guest: 22, host: 2221
  config.vm.provision "shell", inline: <<-SHELL
    echo 'Hello World'
    whoami
  SHELL
  config.vm.provision "shell", path: "sample.sh"
end
```

You can now run the virtual machine using `vagrant up`!

# future featues
I am currently trying to add as much as possible from the `config.vm` so that we can also start using toml
for bigger / more complicated Vagrantfile configurations.
