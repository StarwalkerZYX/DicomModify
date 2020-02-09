package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/StarwalkerZYX/DicomModify/folderutils"
	"github.com/suyashkumar/dicom"
	"github.com/suyashkumar/dicom/dicomlog"
	"github.com/suyashkumar/dicom/dicomtag"
	"github.com/suyashkumar/dicom/element"
	"github.com/suyashkumar/dicom/write"
	//"github.com/suyashkumar/dicom/element"
)

var (
	anonymize      = flag.Bool("anonymize", true, "Anonymize the dicom folder")
	dicomFolder    = flag.String("dicom-folder", "", "The full path of the DICOM folder to modify")
	modifiedFolder = flag.String("modified-folder", "", "The full path where the modified DICOM files will be written to")
	verbose        = flag.Bool("verbose", false, "Activate high verbosity log operation")
)

func setByName(elems []*element.Element, elementName string, newVal string) error {

	t, err := dicomtag.FindByName(elementName)
	if err != nil {
		return err
	}
	for idx, elem := range elems {
		if elem.Tag == t.Tag {
			var newEle *element.Element
			newEle, err = element.NewElement(t.Tag, newVal)

			elems[idx] = newEle
			return nil
		}
	}
	return fmt.Errorf("Could not find element named '%s' in dicom file", elementName)

}

func SetByName(ds *element.DataSet, elementName string, newVal string) error {
	return setByName(ds.Elements, elementName, newVal)
}

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

		parsedData, err := p.Parse(dicom.ParseOptions{DropPixelData: false})
		if parsedData == nil || err != nil {
			log.Panicf("Error reading %s: %v. Maybe not a valid DICOM File. Ignore it.", f, err)
			continue
		}

		var newF string

		if strings.HasPrefix(f, *dicomFolder) {
			newF = strings.TrimPrefix(f, *dicomFolder)
			newF = *modifiedFolder + newF
		}

		modifiedDicomFiles = append(modifiedDicomFiles, newF)
		fmt.Println(newF)

		k := filepath.ToSlash(newF)
		println(k)

		d := path.Dir(k)
		println(d)

		err = os.MkdirAll(d, os.ModePerm)
		if err != nil {
			fmt.Println(err)
		}

		dataSetModified := *parsedData

		var pn *element.Element
		pn, err = dataSetModified.FindElementByName("PatientName")

		v := pn.Value[0]
		fmt.Println(reflect.ValueOf(v))

		SetByName(&dataSetModified, "PatientName", "NewName")

		pn, err = dataSetModified.FindElementByName("PatientName")

		err = write.DataSetToFile(newF, &dataSetModified)
		if err != nil {
			fmt.Println(err)
		}

		v = pn.Value[0]
		fmt.Println(reflect.ValueOf(v))

		if err != nil {
			println(err)
		} else {
			patientName, _ := pn.GetString()
			patientName = "NewName"
			println(patientName)
		}
	}

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(filePaths)

	if *verbose {
		dicomlog.SetLevel(math.MaxInt32)
	}

}
