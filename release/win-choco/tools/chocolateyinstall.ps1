$ErrorActionPreference = 'Stop';
$toolsDir   = "$(Split-Path -parent $MyInvocation.MyCommand.Definition)"
$url        = 'https://github.com/eirannejad/pushcsv/releases/download/v1.4/pushcsv-win64.zip'

$packageArgs = @{
  packageName   = $env:ChocolateyPackageName
  unzipLocation = $toolsDir
  fileType      = 'EXE'
  url           = $url
  softwareName  = 'pushcsv*'
  checksum      = '83C7C528E5348353C15C2038AC248CE6CDC7A211BABDDEE8E08B88E4D532135E'
  checksumType  = 'sha256'
}

Install-ChocolateyZipPackage @packageArgs