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