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