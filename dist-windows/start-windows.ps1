$minVagrantVersion = "1.8.0"
$minVirtualBoxVersion = "5.0.0"

Function Compare-Semver([string]$Version, [string]$MinVersion) {
  if (!$Version) {
    return 1
  }

  $firstMajor, $firstMinor, $firstPatch = $Version.Split("{.}")
  $secondMajor, $secondMinor, $secondPatch = $MinVersion.Split("{.}")

  if ($firstMajor -gt $secondMajor) { return 0 }
  if ($firstMajor -lt $secondMajor) { return 1 }
  if ($firstMinor -gt $secondMinor) { return 0 }
  if ($firstMinor -lt $secondMinor) { return 1 }
  if ($firstPatch -ge $secondPatch) { return 0 }

  return 1
}

if ($Env:VBOX_INSTALL_PATH) {
  $vBoxManagePath = "$Env:VBOX_INSTALL_PATH" + "VBoxManage.exe"
  $virtualBoxVersion = & $vBoxManagePath --version
}

if ($Env:VBOX_MSI_INSTALL_PATH) {
  $vBoxManagePath = "$Env:VBOX_MSI_INSTALL_PATH" + "VBoxManage.exe"
  $virtualBoxVersion = & $vBoxManagePath --version
}

$virtualBoxSemver = Select-String -pattern [\d]+.[\d]+.[\d]+ -InputObject $virtualBoxVersion | % { $_.Matches } | % { $_.Value }
if ($(Compare-Semver $virtualBoxSemver $minVirtualBoxVersion) -eq 1) {
  Write-Host "`nVirtualbox >= $minVirtualboxVersion must be installed to use PCF Dev."
  exit
}

$error.clear()
Get-Command -ErrorAction SilentlyContinue vagrant >$null 2>&1
if (!$error) {
  $vagrantVersion = vagrant -v
  $vagrantSemver = $vagrantVersion.Split("{ }")[1]
}

if ($(Compare-Semver $vagrantSemver $minVagrantVersion) -eq 1) {
  Write-Host "`nVagrant >= $MinVagrantVersion must be installed to use PCF Dev."
  exit
}

$visualRedistributableInstalled = (Get-ItemProperty -ErrorAction SilentlyContinue -Path 'HKLM:\SOFTWARE\Wow6432Node\Microsoft\VisualStudio\10.0\VC\VCRedist\x86' -Name Installed).Installed
if ($visualRedistributableInstalled -ne 1) {
  Write-Host "`nThe Visual C++ 2010 Redistributable Package (x86) is required to run Vagrant on Windows systems. Please follow this link to download and install:`n
  https://www.microsoft.com/en-us/download/details.aspx?id=5555`n"
  exit
}

$inets = Ipconfig
$vBoxInets = & $vBoxManagePath list hostonlyifs

for ($i = 1; $i -le 9; $i++) {
  $inet = "192.168.$i$i.1"
  $Env:PCFDEV_IP = "192.168.$i$i.11"

  if (!(Select-String -quiet -pattern $inet -InputObject $inets) -or ((Select-String -quiet -pattern $inet -InputObject $vBoxInets) -and !(Test-Connection -quiet -count 1 $Env:PCFDEV_IP))) {
    $Env:PCFDEV_DOMAIN = "local$i.pcfdev.io"
    break
  }
}

Trap {
  Remove-Item Env:\PCFDEV_IP
  Remove-Item Env:\PCFDEV_DOMAIN
}

$vagrantDir = Split-Path $MyInvocation.MyCommand.Path -Parent
pushd "$vagrantDir" >$null
  vagrant up --provider=virtualbox
popd >$null
