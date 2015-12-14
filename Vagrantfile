Vagrant.configure("2") do |config|
  config.vm.box = "micropcf/base"
  config.vm.box_version = "0"
  config.vm.define "micropcf" do |micropcf|
  end
  config.vm.synced_folder ".", "/vagrant", disabled: true

  provider_is_aws = (!ARGV.nil? && ARGV.join(' ').match(/provider(=|\s+)aws/))

  if Vagrant.has_plugin?("vagrant-proxyconf") && !provider_is_aws
    config.proxy.http = ENV["http_proxy"] || ENV["HTTP_PROXY"]
    config.proxy.https = ENV["https_proxy"] || ENV["HTTPS_PROXY"]
    config.proxy.no_proxy = [
      "localhost", "127.0.0.1",
      (ENV["MICROPCF_IP"] || "192.168.11.11"),
      (ENV["MICROPCF_DOMAIN"] || "local.micropcf.io")
    ].join(',')
  end

  resources = calculate_resource_allocation

  config.vm.provider "virtualbox" do |v|
    v.customize ["modifyvm", :id, "--memory", resources[:memory]]
    v.customize ["modifyvm", :id, "--cpus", resources[:cpus]]
    v.customize ["modifyvm", :id, "--ioapic", "on"]
  end

  config.vm.provider "vmware_fusion" do |v|
    v.vmx["memsize"] = resources[:memory]
    v.vmx["numvcpus"] = resources[:cpus]
  end

  config.vm.provider "vmware_workstation" do |v|
    v.vmx["memsize"] = resources[:memory]
    v.vmx["numvcpus"] = resources[:cpus]
  end

  config.vm.provider "aws" do |aws, override|
    aws.access_key_id = ENV["AWS_ACCESS_KEY_ID"]
    aws.secret_access_key = ENV["AWS_SECRET_ACCESS_KEY"]
    aws.keypair_name = ENV["AWS_SSH_PRIVATE_KEY_NAME"]
    aws.region = ENV["AWS_REGION"] || 'us-east-1'
    aws.instance_type = "m4.large"
    aws.ebs_optimized = true
    aws.tags = { "Name" => (ENV["AWS_INSTANCE_NAME"] || "micropcf") }
    aws.ami = ""

    override.ssh.username = "ubuntu"
    override.ssh.private_key_path = ENV["AWS_SSH_PRIVATE_KEY_PATH"]
  end

  if provider_is_aws
    network_config = <<-SCRIPT
      public_ip="$(curl http://169.254.169.254/latest/meta-data/public-ipv4)"
      domain="#{ENV["MICROPCF_DOMAIN"] || "${public_ip}.xip.io"}"
      garden_ip="$(ip route get 1 | awk '{print $NF;exit}')"
    SCRIPT
  else
    public_ip = ENV["MICROPCF_IP"] || "192.168.11.11"
    default_domain = (public_ip == "192.168.11.11") ? "local.micropcf.io" : "#{public_ip}.xip.io"

    network_config = <<-SCRIPT
      domain="#{ENV["MICROPCF_DOMAIN"] || default_domain}"
      garden_ip="#{public_ip}"
    SCRIPT

    config.vm.network "private_network", ip: public_ip
  end

  config.vm.provision "shell" do |s|
    s.inline = <<-SCRIPT
      set -e
      #{network_config}

      echo "GARDEN_IP=$garden_ip" >> /var/micropcf/setup
      echo "DOMAIN=$domain" >> /var/micropcf/setup
      echo 'HOST_ID=micropcf' >> /var/micropcf/setup

      /var/micropcf/run
    SCRIPT
  end
end

def calculate_resource_allocation
  cpus = ENV['VM_CORES'] ? ENV['VM_CORES'].to_i : nil
  memory = ENV['VM_MEMORY'] ? ENV['VM_MEMORY'].to_i : nil

  case RUBY_PLATFORM
  when /darwin/i
    cpus ||= `sysctl -n hw.ncpu`.to_i
    memory ||= `sysctl -n hw.memsize`.to_i / 1024 / 1024 / 4
  when /linux/i
    cpus ||= `nproc`.to_i
    memory ||= `grep 'MemTotal' /proc/meminfo | sed -e 's/MemTotal://' -e 's/ kB//'`.to_i / 1024 / 4
  when /cygwin|mswin|mingw|bccwin|wince|emx/i
    cpus ||= `wmic computersystem get numberoflogicalprocessors`.split("\n")[2].to_i
    memory ||= `wmic OS get TotalVisibleMemorySize`.split("\n")[2].to_i / 1024 / 4
  else
    cpus ||= 2
    memory ||= 4096
  end

  {memory: memory / 4 * 4, cpus: cpus}
end
