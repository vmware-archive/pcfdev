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
  $virtualBoxVersion = & ("$Env:VBOX_INSTALL_PATH" + "VBoxManage.exe") --version
}

if ($Env:VBOX_MSI_INSTALL_PATH) {
  $virtualBoxVersion = & ("$Env:VBOX_MSI_INSTALL_PATH" + "VBoxManage.exe") --version
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

$vagrantDir = Split-Path $MyInvocation.MyCommand.Path -Parent
pushd "$vagrantDir" >$null
  vagrant up --provider=virtualbox
popd >$null
