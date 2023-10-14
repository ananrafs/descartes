
dumpFile = ${pwd}/dump
foldername = ./dump/${folder}

run :
	@go run ./main.go -fact=${fact} -law=${law} -out=${out} -folder=${foldername}