$ErrorActionPreference = 'Stop';
$toolsDir   = "$(Split-Path -parent $MyInvocation.MyCommand.Definition)"
$url        = 'https://github.com/eirannejad/pushcsv/releases/download/v1.6/pushcsv-win64.zip'

$packageArgs = @{
  packageName   = $env:ChocolateyPackageName
  unzipLocation = $toolsDir
  fileType      = 'EXE'
  url           = $url
  softwareName  = 'pushcsv*'
  checksum      = '518C7CDFAA76D8E5D8DE14639528848BA9D943B63632E91020F1D5EA72DD2804'
  checksumType  = 'sha256'
}

Install-ChocolateyZipPackage @packageArgs