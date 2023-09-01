package generate

import (
	"flag"
	"io"
	"os"
	"strings"

	"go-micro.dev/v4/util/log"
)

func main() {
	pack := flag.String("package", "main", "Package name of the generated file")
	fileSuffix := flag.String("fileSuffix", ".cql", "File suffix to read into string")
	dir := flag.String("dir", "main", "directory to read into string")
	fileName := flag.String("fileName", "main", "Name of the generated file")

	flag.Parse()

	prefix := *dir + "/"

	fs, err := os.ReadDir(*dir)
	if err != nil {
		log.Fatal("Can't find query directory", err)
	}
	out, err := os.Create(*fileName)
	if err != nil {
		log.Fatal("Can't create output file", err)
	}
	_, err = out.Write([]byte("package " + *pack + "\n\nconst (\n"))
	if err != nil {
		log.Fatal("Error writing output file", err)
	}

	constantNames := "const (\n"

	for _, f := range fs {
		if strings.HasSuffix(f.Name(), *fileSuffix) {
			name := strings.TrimSuffix(f.Name(), *fileSuffix)
			_, err = out.Write([]byte(name + " = `"))
			if err != nil {
				log.Fatal("Error writing output file stream", err)
			}
			f, err := os.Open(prefix + f.Name())
			if err != nil {
				log.Fatal("Error opening a query file", err)
			}
			_, err = io.Copy(out, f)
			if err != nil {
				log.Fatal("Error writing to output file stream", err)
			}
			_, err = out.Write([]byte("`\n\n"))
			if err != nil {
				log.Fatal("Error writing to output file stream", err)
			}
			constantNames += name + "Name = \"" + name + "\" \n"
		}
	}

	_, err = out.Write([]byte(")\n\n"))
	if err != nil {
		log.Fatal("Error writing to output file stream", err)
	}

	constantNames += ")\n"
	_, err = out.Write([]byte(constantNames))
	if err != nil {
		log.Fatal("Error writing to output file stream", err)
	}
}
