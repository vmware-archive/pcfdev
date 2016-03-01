$vagrantDir = Split-Path $MyInvocation.MyCommand.Path -Parent
pushd "$vagrantDir" >$null
  vagrant halt
popd >$null
