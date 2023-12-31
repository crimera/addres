package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
)

type Resources struct {
	XMLName   xml.Name   `xml:"resources"`
	Resources []Resource `xml:"public"`
}

type Resource struct {
	Type string `xml:"type,attr"`
	Name string `xml:"name,attr"`
	Id   string `xml:"id,attr"`
}

func main() {
	// get resources
	var input string

	if len(os.Args) < 2 {
		log.Print("No input file, using res' public.xml")

		input = "res/values/public.xml"
	} else {
		input = os.Args[1]
	}

	public, openError := os.Open(input)
	if openError != nil {
		log.Fatalln(openError)
	}

	defer func() {
		closeErr := public.Close()
		if closeErr != nil {
			log.Panicln("Something happened on close: ", closeErr)
		}
	}()

	byteValue, readErr := io.ReadAll(public)
	if readErr != nil {
		log.Fatalln(readErr)
	}

	var resources Resources

	xml.Unmarshal(byteValue, &resources)

	// get type
	fmt.Println("Resource type:")
	resTypes := map[int]string{
		1: "drawable",
		2: "color",
		3: "style",
		4: "id",
	}

	for i := 1; i <= len(resTypes); i++ {
		fmt.Printf("%d: %s\n", i, resTypes[i])
	}

	var choice int
	fmt.Print("\nchoice: ")
	fmt.Scanln(&choice)

	resType := resTypes[choice]

	// get name
	var name string
	fmt.Print("Resouce name: ")
	fmt.Scanln(&name)

	// get the highest id of type
	resourceSlice := resources.Resources
	id := "0x" + strconv.FormatInt(GetHighestIdOfType(resourceSlice, resType)+1, 16)

	resource := Resource{
		Name: name,
		Type: resType,
		Id:   id,
	}

	resourceSlice = append(resourceSlice, resource)

	publicStruct := Resources{
		Resources: resourceSlice,
	}

	// Build xml
	v, decodeErr := xml.MarshalIndent(publicStruct, "", "  ")
	if decodeErr != nil {
		log.Fatalln("Failed to encode xml: ", decodeErr)
	}

	publicXml := strings.ReplaceAll(string(v), "></public>", "/>")

	// Saving
	writeErr := os.WriteFile("out.xml", []byte(publicXml), 0666)
	if writeErr != nil {
		log.Fatalln("Failed saving file: ", writeErr)
	}
}

func GetHighestIdOfType(resources []Resource, resType string) int64 {
	var ids []int64

	for i := 0; i < len(resources); i++ {
		if resources[i].Type == resType {
			id, parseErr := strconv.ParseInt(resources[i].Id, 0, 0)
			if parseErr != nil {
				log.Println("Could not parse id: ", parseErr)
				continue
			}

			ids = append(ids, id)
		}
	}

	return slices.Max(ids)
}
