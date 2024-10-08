src = main.go comand_include.go command_exclude.go ../pkg/filedriver/filedriver.go ../pkg/filedriver/fastaoriginal/fasta.go ../pkg/filedriver/fastaoriginal/readnontarget.go ../pkg/filedriver/fastaoriginal/reattarget.go ../pkg/keccore/keccore.go ../pkg/keccore/kecoriginal/excluder.go ../pkg/keccore/kecoriginal/includer.go ../pkg/maps/maps.go ../pkg/maps/structs.go ../pkg/maps/kmapbasic/nontarget.go ../pkg/maps/kmapbasic/nontarget_test.go ../pkg/maps/kmapbasic/target.go ../pkg/maps/mmapbasic/maskmap.go ../pkg/utils/files.go ../pkg/utils/reversecomplement.go

all: kec
kec: $(src)
	go build -o kec
test: kec
	./kec e -t targets/ -n neighbors -o out.fasta
clean:
	rm -f kec
