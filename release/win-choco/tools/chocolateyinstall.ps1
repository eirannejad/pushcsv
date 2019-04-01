$ErrorActionPreference = 'Stop';
$toolsDir   = "$(Split-Path -parent $MyInvocation.MyCommand.Definition)"
$url        = 'https://github.com/eirannejad/pushcsv/releases/download/v1.1/pushcsv-win64.zip'

$packageArgs = @{
  packageName   = $env:ChocolateyPackageName
  unzipLocation = $toolsDir
  fileType      = 'EXE'
  url           = $url
  softwareName  = 'pushcsv*'
  checksum      = 'A75A9E361A584582FC145181A793BC07F142D76C848B3C53CA9BE640A664B252'
  checksumType  = 'sha256'
}

Install-ChocolateyZipPackage @packageArgs