package main

import (
	"io"
	"sort"
	"strconv"
	"strings"

	//"bytes"
	//"fmt"
	//"io"
	"os"

	//"path/filepath"
	//"strings"
	"github.com/labstack/gommon/log"
)

func dirTree(out io.Writer,path string,printFiles bool) error{
	file, err := os.Open(path) // For read access.
	if err != nil {
		log.Fatal(err)
	}

	// Read dir`s items names
	names,err := file.Readdirnames(-1)
	if err != nil {
		return err
	}
	if !printFiles{
		names = getOnlyDir(names)
	}

	isLast := false
	len := len(names)

	sort.Strings(names)
	for index,item := range names{
		// If item is last, next call readDir with isLast = true
		if len == index +1{
			isLast = true
		}

		// Call readDir
		err := printAndRead(out,path +"/" + item,"",isLast,printFiles)
		if err != nil{
			return err
		}
	}

	return nil
}

func getOnlyDir(names []string) []string{
	// Separating directories from files
	for i := 0; i < len(names); i++{
		if strings.Contains(names[i],"." ){
			copy(names[i:], names[i+1:])
			names[len(names)-1] = ""
			names = names[:len(names)-1]
			i--
		}
	}
	return names
	
}

func printAndRead(out io.Writer,path,preStr string,isLast,printFiles bool)error{
	// For read access
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}


	// Read stat of file.
	fInfo,err := file.Stat()
	if err != nil{
		log.Error(err)
		return err
	}

	// if you don't need to print files
	if !fInfo.IsDir() && !printFiles{
		return nil
	}


	// check size of file
	printSize := ""
	if !fInfo.IsDir() {
		sizeStr := strconv.FormatInt(fInfo.Size(), 10)
		if sizeStr == "0"{
			printSize = " (empty)"
		}else{
			printSize = " (" +sizeStr+ "b)"
		}
	}


	// Print and put the desired character
	if isLast {
		if _,err := out.Write([]byte(preStr + "└───" + fInfo.Name() + printSize + "\n")); err != nil{
			return err
		}
		//fmt.Println(preStr + "└───" + fInfo.Name())
	}else{
		if _,err := out.Write([]byte(preStr + "├───" + fInfo.Name()+ printSize + "\n")); err != nil{
			return err
		}
		//fmt.Println(preStr + "├───" + fInfo.Name())
	}

	// If it dir call readDir

	if fInfo.IsDir() == true{
		err = readDir(file,out,path,preStr,isLast,printFiles)
		if err != nil{
			log.Error(err)
			return err
		}
	}

	return nil
}

func readDir(file *os.File,out io.Writer,path,preStr string,isLast,printFiles bool) error{
	// Read dir`s items names
	names,err := file.Readdirnames(0)
	if err != nil{
		return err
	}

	// if you don't need to print files
	if !printFiles {
		names = getOnlyDir(names)
	}

	//  Put the desired character
	if isLast{
		preStr += "\t"
	}else{
		preStr += "│	"
	}

	// Set isLast false by default
	isLast = false

	// Check len
	len := len(names)

	// Sort
	sort.Strings(names)
	// Start cycle by names
	for index,item := range names{
		// If item is last, next call readDir with isLast = true

			if len == index +1{
				isLast = true
			}

		// Call readDir
		err := printAndRead(out,path + "/" + item,preStr,isLast,printFiles)
		if err != nil{
			log.Error(err)
			return err
		}
	}

	return nil
}
func main() {
	out := os.Stdout

	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "[-f]"

	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}

}