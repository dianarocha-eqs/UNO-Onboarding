#!/bin/bash
 
echo " - Navigate to root"
 
if ! command -v swag &> /dev/null; then
    echo " - Install swag"
    go install github.com/swaggo/swag/cmd/swag@v1.16.3
fi
 
echo " - Export gopath"
export PATH=$PATH:$HOME/go/bin
 
echo " - Generate code using swag"
swag init --parseDependency --overridesFile .swaggo
 
echo " - Done!"

 