package xmldom_test

import (
	"fmt"

	"github.com/subchen/go-xmldom"
)

const (
	ExampleXml = `<?xml version="1.0" encoding="UTF-8"?>
<testsuites>
	<testsuite tests="2" failures="0" time="0.009" name="github.com/subchen/go-xmldom">
		<properties>
			<property name="go.version">go1.8.1</property>
		</properties>
		<testcase classname="go-xmldom" name="ExampleParseXML" time="0.004"></testcase>
		<testcase classname="go-xmldom" name="ExampleParse" time="0.005"></testcase>
	</testsuite>
</testsuites>
`
)

func ExampleParseXML() {
	node := xmldom.Must(xmldom.ParseXML(ExampleXml)).Root
	fmt.Printf("name = %v\n", node.Name)
	fmt.Printf("attributes.len = %v\n", len(node.Attributes))
	fmt.Printf("children.len = %v\n", len(node.Children))
	fmt.Printf("root = %v", node == node.Root())
	// Output:
	// name = testsuites
	// attributes.len = 0
	// children.len = 1
	// root = true
}

func ExampleNode_GetAttribute() {
	node := xmldom.Must(xmldom.ParseXML(ExampleXml)).Root
	attr := node.FirstChild().GetAttribute("name")
	fmt.Printf("%v = %v\n", attr.Name, attr.Value)
	// Output:
	// name = github.com/subchen/go-xmldom
}

func ExampleNode_GetChildren() {
	node := xmldom.Must(xmldom.ParseXML(ExampleXml)).Root
	children := node.FirstChild().GetChildren("testcase")
	for _, c := range children {
		fmt.Printf("%v: name = %v\n", c.Name, c.GetAttributeValue("name"))
	}
	// Output:
	// testcase: name = ExampleParseXML
	// testcase: name = ExampleParse
}

func ExampleDocument_XML() {
	doc := xmldom.Must(xmldom.ParseXML(ExampleXml))
	fmt.Println(doc.XML())
	// Output:
	// <?xml version="1.0" encoding="UTF-8"?><testsuites><testsuite tests="2" failures="0" time="0.009" name="github.com/subchen/go-xmldom"><properties><property name="go.version">go1.8.1</property></properties><testcase classname="go-xmldom" name="ExampleParseXML" time="0.004" /><testcase classname="go-xmldom" name="ExampleParse" time="0.005" /></testsuite></testsuites>
}

func ExampleDocument_XMLPretty() {
	doc := xmldom.Must(xmldom.ParseXML(ExampleXml))
	fmt.Println(doc.XMLPretty())
	// Output:
	// <?xml version="1.0" encoding="UTF-8"?>
	// <testsuites>
	//   <testsuite tests="2" failures="0" time="0.009" name="github.com/subchen/go-xmldom">
	//     <properties>
	//       <property name="go.version">go1.8.1</property>
	//     </properties>
	//     <testcase classname="go-xmldom" name="ExampleParseXML" time="0.004" />
	//     <testcase classname="go-xmldom" name="ExampleParse" time="0.005" />
	//   </testsuite>
	// </testsuites>
}