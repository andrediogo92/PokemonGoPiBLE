param([string[]]$command = ("build", "-o", "./host", "./examples"))
$version = "1.13.5-alpine"
$base = "golang"
$build = "${base}:${version}"
$full_command = @()
$full_command += ("go")

foreach ($item in $command)
{
    $full_command += $item
}


# Ensure Go and the ARM toolchain are installed.
If (-Not (Get-Command docker -ErrorAction SilentlyContinue)) {
    'It looks like Docker is not installed on your system.' | Write-Error
    Exit
}

If (-Not (docker images | Select-String "$base(\s+)$version" -ErrorAction SilentlyContinue)) {
    "Image not detected, attempt to pull $build" | Write-Output
    docker pull $build
}

docker run --rm `
    -e GOARCH=$arch `
    -e GOOS=$os `
    -v ${PWD}:/usr/src/app -w /usr/src/app `
    ${build} `
    $full_command