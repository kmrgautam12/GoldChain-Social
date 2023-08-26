cd ../..
cd src
set GOARCH=amd64
set GOOS=linux 
go build -o main .
zip function.zip main
aws lambda update-function-code \
    --function-name  gautam_api \
    --zip-file fileb://function.zip
rm function.zip
