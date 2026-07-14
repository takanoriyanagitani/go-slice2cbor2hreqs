#!/bin/bash

export ENV_BASE_URL=http://127.0.0.1:8080

urlkey(){
	printf 68 | xxd -r -ps
	printf url_path
}

hdrkey(){
	printf 66 | xxd -r -ps
	printf header
}

bdykey(){
	printf 69 | xxd -r -ps
	printf post_body
}

urlval(){
	printf 6c | xxd -r -ps
	printf /api/v1/helo
}

hdrval(){
	printf a0 | xxd -r -ps
}

bdyval(){
	printf 43303132 | xxd -r -ps
}

input1(){
	printf a3 | xxd -r -ps
	urlkey; urlval
	hdrkey; hdrval
	bdykey; bdyval
}

input2(){
	input1
	input1
}

input4(){
	input2
	input2
}

input4w0(){
	input4
	printf '\0'
}

input4w0 | ./cmd/slice2cbor2hreqs/slice2cbor2hreqs
