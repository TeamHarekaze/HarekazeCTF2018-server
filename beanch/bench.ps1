# File name           : bench.ps1
# Author              : Hayato Doi
# Outline             : urlにアクセスしてレスポンス速度を計測するスクリプト
# license             : None
# Copyright (c) 2018, Hayato Doi

$url = 'http://27.133.153.12:5000/question'

for( $i=0; $i -ne 100; $i++ ){
    # track execution time:
    $timeTaken = Measure-Command -Expression {
        Invoke-WebRequest -Uri $url
    }
    $milliseconds = $timeTaken.TotalMilliseconds
    $milliseconds = [Math]::Round($milliseconds, 1)
    for(;;){
        Try{
            $milliseconds | Out-File -Append -FilePath out.txt
            break
        } Catch [System.IO.IOException] {
        }
    }
}
