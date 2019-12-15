package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

var pathPrefix string = ""

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}

var prefix string

func dirTree(out io.Writer , path string,  printFiles bool) error {

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	if !printFiles {
		var ofi []os.FileInfo

		for _, file := range files {
			if file.Mode().IsDir() {
				ofi = append(ofi, file)
			}
		}

		files = ofi
	}

	for i, f := range files {

		isLast := i == len(files)-1


		if f.IsDir() {

			out.Write([]byte(prefix))
			if isLast {
				out.Write([]byte("└───"))
			} else {
				out.Write([]byte("├───"))
			}


			out.Write([]byte(f.Name()))
			out.Write([]byte("\n"))
			sPrefix := prefix
			if isLast {
				prefix += "	"
			} else {
				prefix += "│	"
			}
			dirTree(out,path + "/" + f.Name(), printFiles)
			prefix = sPrefix
		} else if printFiles {
			out.Write([]byte(prefix))
			if isLast {
				out.Write([]byte("└───"))
			} else {
				out.Write([]byte("├───"))
			}


			out.Write([]byte(f.Name()))
			if f.Size() > 0 {
				out.Write([]byte(fmt.Sprintf(" (%db)", f.Size()) ))
			} else {
				out.Write([]byte(" (empty)"))
			}
			out.Write([]byte("\n"))
		}
	}

	return nil
}

func dirTree1(out io.Writer , path string,  printFiles bool) error {

	o := out
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}


	if !printFiles {
		var ofi []os.FileInfo

		for _, file := range files {
			if file.Mode().IsDir() {
				ofi = append(ofi, file)
			}
		}

		files = ofi
	}



	for i, file := range files {

		if  !file.Mode().IsDir() && !printFiles{
			continue
		}

		fullpath := path + string(os.PathSeparator) + file.Name()

		o.Write ([]byte(pathPrefix))
		if i + 1 < len(files) {
			o.Write ([]byte("├───"))
		} else {
			o.Write ([]byte("└───"))
		}
		o.Write ([]byte(file.Name()))


		if  !file.Mode().IsDir() {
			if file.Size() == 0 {
				o.Write ([]byte(" (empty)"))
			} else {
				o.Write ([]byte(fmt.Sprintf(" (%db)", file.Size())))
			}
		}
		o.Write ([]byte("\n"))

		if file.Mode().IsDir() {

			if i+1 < len(files) {
				pathPrefix = pathPrefix + "│\t"
			} else {
				pathPrefix = pathPrefix + "\t"
			}

			err = dirTree1(out, fullpath, printFiles)

			if i+1 < len(files) {
				pathPrefix = strings.Replace(pathPrefix, "│\t", "", 1)
			} else {
				pathPrefix = strings.Replace(pathPrefix, "\t", "", 1)
			}
		}
	}

	return err
}