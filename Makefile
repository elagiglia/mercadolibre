test:
	cd mockapi ; npm install
	node mockapi/app.js &
	export GOPATH=export PWD=`pwd`
	go test -v src/github.com/mercadolibre/sdk/*
		
	kill `cat /tmp/mockapi.pid`

deploy:
	export GOPATH=export PWD=`pwd`
	go build -v github.com/mercadolibre/sdk/
	#cd mockapi ; npm install
	#node mockapi/app.js &
	#mvn -DaltDeploymentRepository=snapshot-repo::default::file:../java-sdk-repo/snapshots clean deploy
	#kill `cat /tmp/mockapi.pid`

.PHONY: test
