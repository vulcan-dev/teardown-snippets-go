package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

type SInput struct {
	Name string `xml:"name,attr"`
	Type string `xml:"type,attr"`
	Optional string `xml:"optional,attr"`
	Desc string `xml:"desc,attr"`
}

type SOutput struct {
	Name string `xml:"name,attr"`
	Type string `xml:"type,attr"`
	Desc string `xml:"desc,attr"`
}

type SFunctions struct {
	Name string `xml:"name,attr"`
	Output SOutput `xml:"output"`
	Input []SInput `xml:"input"`
}

type SAPI struct {
	XMLName xml.Name `xml:"api"`
	Function []SFunctions `xml:"function"`
}

func main() {
	/* Get current path */
	currentPath, err := os.Getwd(); if err != nil {
		log.Fatalln("Failed getting current path: ", err)
	}

	/* Find File */
	var files[]string
	filepath.Walk(currentPath, func(path string, f os.FileInfo, _ error) error {
		if !f.IsDir() {
			r, err := regexp.MatchString(".xml", f.Name()); if err == nil && r {
				files = append(files, f.Name())
			}
		}

		return nil
	})

	/* Open File */
	file, err := os.Open(files[0]); if err != nil {
		log.Fatalln(fmt.Sprintf("Failed opening file \"%s\": %s", files[0], err))
	}

	defer file.Close()

	/* Read data from file */
	byteValue, _ := ioutil.ReadAll(file)
	var data SAPI

	if err := xml.Unmarshal(byteValue, &data); err != nil {
		log.Fatalln("Unmarshal Failed")
	}

	out, err := os.OpenFile("log.txt", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644); if err != nil {
		log.Println(err)
	}
	defer out.Close()

	valid := false

	for i := 0; i < len(data.Function); i++ {
		
		for j, _ := range data.Function[i].Input {
			out.WriteString(fmt.Sprintf("[Input] Name: %s\n", data.Function[i].Input[j].Name))
			out.WriteString(fmt.Sprintf("[Input] Type: %s\n", data.Function[i].Input[j].Type))
			out.WriteString(fmt.Sprintf("[Input] Optional: %s\n", data.Function[i].Input[j].Optional))
			out.WriteString(fmt.Sprintf("[Input] Description: %s\n\n", data.Function[i].Input[j].Desc))
			valid = true
		}

		if data.Function[i].Output.Name != "" {
			out.WriteString(fmt.Sprintf("[Out] Name: %s\n", data.Function[i].Output.Name))
			out.WriteString(fmt.Sprintf("[Out] Type: %s\n", data.Function[i].Output.Type))
			out.WriteString(fmt.Sprintf("[Out] Desc: %s\n", data.Function[i].Output.Desc))
			valid = true
		}

		if valid {
			out.WriteString(fmt.Sprintf("===== %s =====\n", data.Function[i].Name))
			out.WriteString("\n")
		} else {
			out.WriteString(fmt.Sprintf("\n== INVALID ==\n%s\n\n", data.Function[i].Name))
		}

		valid = false
	}
}