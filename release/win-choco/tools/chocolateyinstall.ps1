$ErrorActionPreference = 'Stop';
$toolsDir   = "$(Split-Path -parent $MyInvocation.MyCommand.Definition)"
$url        = 'https://github.com/eirannejad/pushcsv/releases/download/v1.0/pushcsv.zip'

$packageArgs = @{
  packageName   = $env:ChocolateyPackageName
  unzipLocation = $toolsDir
  fileType      = 'EXE'
  url           = $url
  softwareName  = 'pushcsv*'
  checksum      = 'C6B929B666B0098B00D6DFF22A03E918A302580D6BE236AD516869790C775D90'
  checksumType  = 'sha256'
}

Install-ChocolateyZipPackage @packageArgs