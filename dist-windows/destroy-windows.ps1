$vagrantDir = Split-Path $MyInvocation.MyCommand.Path -Parent
pushd "$vagrantDir" >$null
  vagrant destroy -f
popd >$null
