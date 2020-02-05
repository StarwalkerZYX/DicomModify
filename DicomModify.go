package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/StarwalkerZYX/DicomModify/folderutils"
	"github.com/suyashkumar/dicom"
	"github.com/suyashkumar/dicom/dicomlog"
)

var (
	anonymize      = flag.Bool("anonymize", true, "Anonymize the dicom folder")
	dicomFolder    = flag.String("dicom-folder", "", "The full path of the DICOM folder to modify")
	modifiedFolder = flag.String("modified-folder", "", "The full path where the modified DICOM files will be written to")
	verbose        = flag.Bool("verbose", false, "Activate high verbosity log operation")
)

func main() {

	//Update usage docs
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n%s <dicom file> [flags]\n", os.Args[0], os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	var err error
	if *dicomFolder == string("") {
		*dicomFolder, err = filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(*dicomFolder)

		*dicomFolder += string("\\dicom\\")
		fmt.Println(*dicomFolder)
	}

	if *modifiedFolder == string("") {
		*modifiedFolder, err = filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(modifiedFolder)

		*modifiedFolder += string("\\modified_dicom\\")
		fmt.Println(modifiedFolder)
	}

	filePaths, err := folderutils.ScanFolder(*dicomFolder)

	var modifiedDicomFiles []string

	for _, f := range filePaths {

		p, err := dicom.NewParserFromFile(f, nil)
		if err != nil {
			log.Panic("error creating new parser", err)
		}

		parsedData, err := p.Parse(dicom.ParseOptions{DropPixelData: true})
		if parsedData == nil || err != nil {
			log.Panicf("Error reading %s: %v", f, err)
		}

		var newF string

		if strings.HasPrefix(f, *dicomFolder) {
			newF = strings.TrimPrefix(f, *dicomFolder)
			newF = *modifiedFolder + string("\\") + newF
		}

		modifiedDicomFiles = append(modifiedDicomFiles, newF)
		fmt.Println(newF)
	}

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(filePaths)

	if *verbose {
		dicomlog.SetLevel(math.MaxInt32)
	}

}
