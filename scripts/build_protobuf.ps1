Param(
    [string] $projRoot
)

$go_module = 'github.com/shenrytech/shdb'

[System.IO.Directory]::CreateDirectory($ts_out) | Out-Null
[System.IO.Directory]::CreateDirectory($python_out) | Out-Null

$protoc_inc_path = (get-item (Get-Command protoc).Path).Directory.Parent.FullName

$google_path = Join-Path $protoc_inc_path include -Resolve

function build_pb_go([string] $file) {
    Write-Host -NoNewline -ForegroundColor Cyan "golang"
    & protoc `
        -I $projRoot `
        -I $google_path `
        --go_opt=module=$go_module `
        --go_out=$projRoot `
        --go-grpc_out=$projRoot `
        --go-grpc_opt=module=$go_module `
        $file
}

function Convert-Protobuf([string] $file) {
    $fname = (get-item $file).Name
    Write-Host -NoNewline -ForegroundColor Cyan "Generating bindings for '$fname' ["
    build_pb_go $file
    Write-Host -ForegroundColor Cyan "]"
}

Get-ChildItem $projRoot -Filter *.proto |
Foreach-Object {
    Convert-Protobuf $_
}
 