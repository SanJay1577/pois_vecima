#Install the swagger in Our Project
#first check the swagger is already installed 
check_install:
	which swagger || GO111MODULE=off go get -u github.com/go-swagger/go-swagger/cmd/swagger

#Generate the swagger file 
swagger: check_install
	GO111MODULE=off swagger generate spec -o ./swagger.yaml --scan-models

swagger_xml: check_install
	GO111MODULE=off swagger generate spec -o ./swagger_xml.yaml --scan-models
