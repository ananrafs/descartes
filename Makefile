
dumpFile = ${pwd}/dump
foldername = ./dump/${folder}
fact = fact
law = law
out = output
folder = rule_random

# use : make run folder=folder
run :
	@go run ./main.go -fact=${fact} -law=${law} -out=${out} -folder=${foldername}