BASEDIR=$(dirname $0)
echo "cleanup : $BASEDIR/../out"
rm $BASEDIR/../out/*
echo "build:go"
go build  -o $BASEDIR/../out/shell-butler $BASEDIR/../main.go
echo "cp: binary"
#mkdir -p $BASEDIR/../package/debian/shell-butler/usr/local/bin
mkdir -p $BASEDIR/../package/debian/shell-butler/opt/shell-butler/bin
cp $BASEDIR/../run $BASEDIR/../package/debian/shell-butler/opt/shell-butler/bin/shell-butler-run
cp $BASEDIR/../out/shell-butler $BASEDIR/../package/debian/shell-butler/opt/shell-butler/bin
echo "run: dpkg-deb"
dpkg-deb --build $BASEDIR/../package/debian/shell-butler  $BASEDIR/../out/shell-butler.deb
echo "done"