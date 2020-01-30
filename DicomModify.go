package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"

	"github.com/StarwalkerZYX/DicomModify/folderscanner"
	"github.com/suyashkumar/dicom/dicomlog"
)

var (
	anonymize           = flag.Bool("anonymize", true, "Anonymize the dicom folder")
	printMetadata       = flag.Bool("print-metadata", false, "Print image metadata")
	extractImages       = flag.Bool("extract-images", false, "Extract images into separate files")
	extractImagesStream = flag.Bool("extract-images-stream", false, "Extract images using frame streaming capability")
	verbose             = flag.Bool("verbose", false, "Activate high verbosity log operation")
)

func main() {

	// Update usage docs
	// flag.Usage = func() {
	// 	fmt.Fprintf(os.Stderr, "Usage of %s:\n%s <dicom file> [flags]\n", os.Args[0], os.Args[0])
	// 	flag.PrintDefaults()
	// }

	var rootPath string

	flag.Parse()
	if len(flag.Args()) == 0 {
		rootPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(rootPath)

		rootPath += "\\dicom\\"
		fmt.Println(rootPath)
	} else {
		rootPath = flag.Arg(0)
	}

	dicomFiles, err := folderscanner.ScanSibblingDICOMFolder(rootPath)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(dicomFiles)

	if *verbose {
		dicomlog.SetLevel(math.MaxInt32)
	}

}
