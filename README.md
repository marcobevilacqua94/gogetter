# gogetter

build with

go build couchbase.go (make sure OS is set to linux)
docker build . --tag gogetter

docker tag gogetter marcobevilacqua94/gogetter
docker push marcobevilacqua94/gogetter

run with 
./couchbase  Administrator password 127.0.0.1 beer-sample _default _default 16
