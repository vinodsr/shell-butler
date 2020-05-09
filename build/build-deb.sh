#!/bin/bash
BASEDIR=$(dirname $0)
echo -e "\U1F9F9 : $(realpath $BASEDIR/../out)"
rm $BASEDIR/../out/*
echo -e "\U1F527 : go build"
go build  -o $BASEDIR/../out/shell-butler $BASEDIR/../main.go
echo -e "\U1F69A : moving binary to deb"
#mkdir -p $BASEDIR/../package/debian/shell-butler/usr/local/bin
mkdir -p $BASEDIR/../package/debian/shell-butler/opt/shell-butler/bin
cp $BASEDIR/../run $BASEDIR/../package/debian/shell-butler/opt/shell-butler/bin/shell-butler-run
cp $BASEDIR/../out/shell-butler $BASEDIR/../package/debian/shell-butler/opt/shell-butler/bin
echo -e "\U1F3C3 : running dpkg-deb"
dpkg-deb --build $BASEDIR/../package/debian/shell-butler  $BASEDIR/../out/shell-butler.deb
echo -e "\U1F4E6 : package available in out directory"
echo -e "\U1F603 : done"