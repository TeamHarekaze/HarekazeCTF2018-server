# File name           : exec.ps1
# Author              : Hayato Doi
# Outline             : コマンドを同時起動するスクリプト
# license             : None
# Copyright (c) 2018, Hayato Doi

Workflow Test-Parallel()
{
    $max = 5
    foreach -parallel ($i in 1..$max) {
        Start-Process powershell -ArgumentList ".\bench.ps1 $i"
    }
}

Test-Parallel
