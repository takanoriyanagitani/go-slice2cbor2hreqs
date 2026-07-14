#!/bin/sh

bname="slice2cbor2hreqs"
bdir="./cmd/${bname}"
oname="${bdir}/${bname}"

mkdir -p "${bdir}"

go \
	build \
	-v \
	./...

go \
	build \
	-v \
	-o "${oname}" \
	"${bdir}"
