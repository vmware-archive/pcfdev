HTTP_PROXY = ENV["http_proxy"] || ENV["HTTP_PROXY"] unless defined? HTTP_PROXY
HTTPS_PROXY = ENV["https_proxy"] || ENV["HTTPS_PROXY"] unless defined? HTTPS_PROXY

Vagrant.configure("2") do |config|
  config.vm.box = "micropcf/colocated"
  config.vm.box_version = "0"

  config.vm.synced_folder ".", "/vagrant", disabled: true
  config.vm.provision "file", source: "micropcf.tgz", destination: "/tmp/micropcf.tgz"

  provider_is_aws = (!ARGV.nil? && ARGV.join(' ').match(/provider(=|\s+)aws/))

  if Vagrant.has_plugin?("vagrant-proxyconf") && !provider_is_aws
    config.proxy.http = HTTP_PROXY
    config.proxy.https = HTTPS_PROXY
    config.proxy.no_proxy = [
      "localhost", "127.0.0.1",
      (ENV["MICROPCF_IP"] || "192.168.11.11"),
      (ENV["MICROPCF_DOMAIN"] || "192.168.11.11.xip.io"),
      ".consul"
    ].join(',')

    config.vm.provision "shell", inline: "grep -i proxy /etc/environment > /var/micropcf/proxy || true"
  end

  cpus = ENV['VM_CORES'] ? ENV['VM_CORES'].to_i : nil
  mem = ENV['VM_MEMORY'] ? ENV['VM_MEMORY'].to_i : nil

  host = RbConfig::CONFIG['host_os']
  if host =~ /darwin/
    cpus ||= `sysctl -n hw.ncpu`.to_i
    mem ||= `sysctl -n hw.memsize`.to_i / 1024 / 1024 / 4
  elsif host =~ /linux/
    cpus ||= `nproc`.to_i
    mem ||= `grep 'MemTotal' /proc/meminfo | sed -e 's/MemTotal://' -e 's/ kB//'`.to_i / 1024 / 4
  elsif host =~ /cygwin|mswin|mingw|bccwin|wince|emx/
    cpus ||= `wmic computersystem get numberofprocessors`.split("\n")[2].to_i
    mem ||= `wmic OS get TotalVisibleMemorySize`.split("\n")[2].to_i / 1024 / 4
  else
    cpus ||= 2
    mem ||= 2048
  end

  mem = nearest_multiple_of_four(mem)

  # Source: https://stefanwrobel.com/how-to-make-vagrant-performance-not-suck
  config.vm.provider "virtualbox" do |v|
    v.customize ["modifyvm", :id, "--memory", mem]
    v.customize ["modifyvm", :id, "--cpus", cpus]
    v.customize ["modifyvm", :id, "--ioapic", "on"]
  end

  config.vm.provider "vmware_fusion" do |v|
    v.vmx["memsize"] = mem
    v.vmx["numvcpus"] = cpus
  end

  config.vm.provider "vmware_workstation" do |v|
    v.vmx["memsize"] = mem
    v.vmx["numvcpus"] = cpus
  end

  config.vm.provider :aws do |aws, override|
    aws.access_key_id = ENV["AWS_ACCESS_KEY_ID"]
    aws.secret_access_key = ENV["AWS_SECRET_ACCESS_KEY"]
    aws.keypair_name = ENV["AWS_SSH_PRIVATE_KEY_NAME"]
    aws.region = ENV["AWS_REGION"] || 'us-east-1'
    aws.instance_type = "m4.large"
    aws.ebs_optimized = true
    aws.tags = { "Name" => (ENV["AWS_INSTANCE_NAME"] || "vagrant") }
    aws.ami = ""
    aws.block_device_mapping = [
      {
        "DeviceName" => "/dev/sda1",
        "Ebs.VolumeType" => "gp2",
        "Ebs.VolumeSize" => 15,
        "Ebs.DeleteOnTermination" => true
      }
    ]

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
    default_domain = (public_ip == "192.168.11.11") ? "local.micropcf.cf" : "#{public_ip}.xip.io"

    network_config = <<-SCRIPT
      domain="#{ENV["MICROPCF_DOMAIN"] || default_domain}"
      garden_ip="#{public_ip}"
    SCRIPT

    config.vm.network "private_network", ip: public_ip
  end

  case RUBY_PLATFORM
  when /darwin/
    ltc_arch_path = 'osx/ltc'
  when /linux/i
    ltc_arch_path = 'linux/ltc'
  when /cygwin|mswin|mingw|bccwin|wince|emx/i
    ltc_arch_path = 'windows/ltc.exe'
  end

  config.vm.provision "shell" do |s|
    s.inline = <<-SCRIPT
      set -e
      #{network_config}

      echo "GARDEN_IP=$garden_ip" >> /var/micropcf/setup
      echo "DOMAIN=$domain" >> /var/micropcf/setup
      echo 'HOST_ID=micropcf-colocated-0' >> /var/micropcf/setup
      echo 'USERNAME=#{ENV["MICROPCF_USERNAME"]}' >> /var/micropcf/setup
      echo 'PASSWORD=#{ENV["MICROPCF_PASSWORD"]}' >> /var/micropcf/setup

      tar xzf /tmp/micropcf.tgz -C /tmp install
      /tmp/install/common
      /tmp/install/brain /tmp/micropcf.tgz
      /tmp/install/cell /tmp/micropcf.tgz
      /tmp/install/start

      echo "MicroPCF is now installed and running."
      echo "----------------------------------------------------------------"
      echo "To obtain a version of ltc that is compatible with your cluster:"
      if [ "#{ltc_arch_path}" = "windows/ltc.exe" ]; then
        echo "Download http://receptor.$domain/v1/sync/#{ltc_arch_path}"
      else
        echo "$ curl -O http://receptor.$domain/v1/sync/#{ltc_arch_path}"
        echo "$ chmod +x ltc"
      fi
      echo "Optionally, move ltc to a directory in your PATH."
      echo "You may then target your cluster using: ltc target $domain"
      echo "----------------------------------------------------------------"
    SCRIPT
  end
end

require 'net/http'
require 'rubygems/package'
require 'uri'
require 'zlib'

provision_required = (!ARGV.nil? && ['up', 'provision', 'reload'].include?(ARGV[0]))
micropcf_tgz = File.join(File.dirname(__FILE__), "micropcf.tgz")
micropcf_tgz_url = defined?(MICROPCF_TGZ_URL) && MICROPCF_TGZ_URL

def download_micropcf_tgz(url)
  uri = URI(url)

  http_args = [uri.host, uri.port]
  proxy_url = (uri.scheme=='https' ? HTTPS_PROXY : HTTP_PROXY)
  if proxy_url
    proxy_uri = URI(proxy_url)
    http_args << proxy_uri.host
    http_args << proxy_uri.port
  end

  Net::HTTP.start(*http_args, use_ssl: uri.scheme == 'https') do |http|
    http.request_get(uri) do |response|
      open('micropcf.tgz', 'wb') do |file|
        response.read_body do |chunk|
          file.write(chunk)
          sleep(0.005)
        end
      end
    end
  end
end

def extract_version(tgz_file, version_file)
  Zlib::GzipReader.open(tgz_file) do |gz|
    tr = Gem::Package::TarReader.new(gz)
    tr.seek "brain/var/micropcf/versions/#{version_file}" do |entry|
      return entry.read.chomp
    end
  end
  return nil
end

if provision_required && File.exists?(micropcf_tgz)
  if micropcf_tgz_url
    tgz_version = extract_version(micropcf_tgz, 'MICROPCF_TGZ')
    url_version = micropcf_tgz_url.match(/\/micropcf-v([^\/]+)\.tgz$/)[1]
    if tgz_version != url_version
      puts "Warning: micropcf.tgz file version (#{tgz_version}) does not match Vagrantfile version (#{url_version})."
      puts 'Re-downloading and replacing local micropcf.tgz...'
      download_micropcf_tgz(micropcf_tgz_url)
    end
  else
    tgz_version = extract_version(micropcf_tgz, 'MICROPCF')
    repo_version = `git rev-parse HEAD`.chomp
    if tgz_version != repo_version && ENV['IGNORE_VERSION_MISMATCH'] != "true"
      puts <<-EOM.gsub(/^ +/, '')
      *******************************************************************************
      Error: micropcf.tgz #{tgz_version[0..6]} != current commit #{repo_version[0..6]}

      The micropcf.tgz file was built using a different commit than the current one.
      To ignore this error, set IGNORE_VERSION_MISMATCH=true in your environment.

      NOTE: As of v0.4.0, the process for deploying MicroPCF via Vagrant has changed.
      Please use the process documented here:
      \thttp://micropcf.cf/docs/vagrant/
      *******************************************************************************
      EOM
      exit(1)
    end
  end
end

if provision_required && !File.exists?(micropcf_tgz)
  if micropcf_tgz_url
    puts 'Local micropcf.tgz not found, downloading...'
    download_micropcf_tgz(micropcf_tgz_url)
  else
    puts <<-EOM.gsub(/^ +/, '')
    *******************************************************************************
    Could not determine MicroPCF version, and no local micropcf.tgz present.

    NOTE: As of v0.4.0, the process for deploying MicroPCF via Vagrant has changed.
    Please use the process documented here:
    \thttp://micropcf.cf/docs/vagrant/
    *******************************************************************************
    EOM
    exit(1)
  end
end

def nearest_multiple_of_four(n)
  (n / 4.0).ceil * 4
end
