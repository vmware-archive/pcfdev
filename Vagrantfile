Vagrant.configure("2") do |config|
  Vagrant.require_version ">= 1.8"

  config.vm.box = "pcfdev/oss"
  config.vm.box_version = "0"

  config.vm.synced_folder ".", "/vagrant", disabled: true

  vagrant_up = (!ARGV.nil? && ARGV.first == 'up')
  vagrant_up_aws = (vagrant_up && ARGV.join(' ').match(/provider(=|\s+)aws/))

  if Vagrant.has_plugin?("vagrant-proxyconf") && !vagrant_up_aws
    config.proxy.http = ENV["HTTP_PROXY"] || ENV["http_proxy"]
    config.proxy.https = ENV["HTTPS_PROXY"] || ENV["https_proxy"]
    ip = (ENV["PCFDEV_IP"] || "192.168.11.11")
    ip_octets = ip.split('.')
    ip_octets[3] = 1
    host_ip = ip_octets.join('.')

    config.proxy.no_proxy = [
      "localhost",
      "127.0.0.1",
      host_ip,
      ip,
      (ENV["PCFDEV_DOMAIN"] || "local.pcfdev.io"),
      (ENV["NO_PROXY"] || ENV["no_proxy"])
    ].compact.join(',')
  end

  resources = calculate_resource_allocation
  if resources[:memory] == 2048 && vagrant_up && !vagrant_up_aws
    puts "WARNING: PCF Dev has reserved 2 GBs out of #{resources[:max_memory] / 1024} GBs total system memory."
    puts "Performance may be impacted."
  end

  config.vm.provider "virtualbox" do |v|
    v.customize ["modifyvm", :id, "--memory", resources[:memory]]
    v.customize ["modifyvm", :id, "--cpus", resources[:cpus]]
    v.customize ["modifyvm", :id, "--ioapic", "on"]
    v.customize ["modifyvm", :id, "--paravirtprovider", "minimal"]
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
    aws.instance_type = "m4.xlarge"
    aws.block_device_mapping = [{
      'DeviceName' => '/dev/sda1',
      'Ebs.VolumeSize' => ENV["AWS_EBS_DISK_SIZE"] || 100,
      'Ebs.DeleteOnTermination' => true,
      'Ebs.VolumeType' => 'gp2',
    }]
    aws.ebs_optimized = true
    aws.tags = { "Name" => (ENV["AWS_INSTANCE_NAME"] || "pcfdev") }
    aws.ami = ""

    override.ssh.username = "ubuntu"
    override.ssh.private_key_path = ENV["AWS_SSH_PRIVATE_KEY_PATH"]
  end

  local_public_ip = ENV["PCFDEV_IP"] || "192.168.11.11"
  local_default_domain = (local_public_ip == "192.168.11.11") ? "local.pcfdev.io" : "#{local_public_ip}.xip.io"
  if !vagrant_up_aws
    config.vm.network "private_network", ip: local_public_ip
  end

  config.vm.provision "shell", run: "always" do |s|
    s.inline = <<-SCRIPT
      set -e
      if public_ip="$(curl -m 2 -s http://169.254.169.254/latest/meta-data/public-ipv4)" && [[ $public_ip =~ ^[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}$ ]]; then
        domain="#{ENV["PCFDEV_DOMAIN"] || "${public_ip}.xip.io"}"
      else
        domain="#{ENV["PCFDEV_DOMAIN"] || local_default_domain}"
        public_ip="#{local_public_ip}"
      fi
      CF_ONLY=#{ENV["CF_ONLY"]} /var/pcfdev/run "$domain" "$public_ip"
      #{cf_cli_present} || echo "Don't have the cf command line utility? Download it from https://github.com/cloudfoundry/cli/releases"
    SCRIPT
  end
end


def calculate_resource_allocation
  cpus = ENV['VM_CORES'] ? ENV['VM_CORES'].to_i : nil
  memory = ENV['VM_MEMORY'] ? ENV['VM_MEMORY'].to_i : nil

  case RUBY_PLATFORM
  when /darwin/i
    sysctl_path = `which sysctl || echo /usr/sbin/sysctl`.chomp
    cpus ||= `#{sysctl_path} -n hw.physicalcpu`.to_i
    max_memory = `#{sysctl_path} -n hw.memsize`.to_i / 1024 / 1024
  when /linux/i
    cpus ||= `grep 'core id' /proc/cpuinfo | uniq | wc -l`.to_i
    max_memory = `grep 'MemTotal' /proc/meminfo | sed -e 's/MemTotal://' -e 's/ kB//'`.to_i / 1024
  when /cygwin|mswin|mingw|bccwin|wince|emx/i
    cpus ||= `wmic computersystem get numberofprocessors`.split("\n")[2].to_i
    max_memory = `wmic OS get TotalVisibleMemorySize`.split("\n")[2].to_i / 1024
  else
    cpus ||= 2
    max_memory = 4096
  end

  memory ||= [[2048, max_memory / 2].max, 4096].min

  {memory: memory / 4 * 4, cpus: cpus, max_memory: max_memory}
end

def cf_cli_present
  case RUBY_PLATFORM
  when /darwin|linux/i
    system("which cf > /dev/null")
  when /cygwin|mswin|mingw|bccwin|wince|emx/i
    system("where /q cf.exe")
  else
    false
  end
end
