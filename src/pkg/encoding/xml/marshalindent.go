package main

import (
	"encoding/xml"
	"fmt"
	"os"
)

type Address struct {
	City, State string
}

type Person struct {
	XMLName   xml.Name `xml:"person"`
	Id        int      `xml:"id,attr"`
	FirstName string   `xml:"name>first"`
	LastName  string   `xml:"name>last"`
	Age       int      `xml:"age"`
	Height    float32  `xml:"height,omitepty"`
	Married   bool
	Address
	Comment string `xml:",comment"`
}

func main() {

	v := &Person{Id: 13, FirstName: "John", LastName: "Doe", Age: 42}
	v.Comment = " Need more details. "
	v.Address = Address{"Hanga Roa", "Easter Island"}

	output, err := xml.MarshalIndent(v, "  ", "    ")
	if err != nil {
		fmt.Printf("error: %v", err)
	}

	os.Stdout.Write(output)
	fmt.Println()

}
