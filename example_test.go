package xmldom_test

import (
	"fmt"
	"testing"

	"github.com/subchen/go-xmldom"
)

const (
	ExampleXml = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE junit SYSTEM "junit-result.dtd">
<testsuites>
	<testsuite tests="2" failures="0" time="0.009" name="github.com/subchen/go-xmldom">
		<properties>
			<property name="go.version">go1.8.1</property>
		</properties>
		<testcase classname="go-xmldom" id="ExampleParseXML" time="0.004"></testcase>
		<testcase classname="go-xmldom" id="ExampleParse" time="0.005"></testcase>
    <testcase xmlns:test="mock" id="AttrNamespace"></testcase>
	</testsuite>
</testsuites>`
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
		fmt.Printf("%v: id = %v\n", c.Name, c.GetAttributeValue("id"))
	}
	// Output:
	// testcase: id = ExampleParseXML
	// testcase: id = ExampleParse
}

func ExampleNode_FindByID() {
	root := xmldom.Must(xmldom.ParseXML(ExampleXml)).Root
	node := root.FindByID("ExampleParseXML")
	fmt.Println(node.XML())
	// Output:
	// <testcase classname="go-xmldom" id="ExampleParseXML" time="0.004" />
}

func ExampleNode_FindOneByName() {
	root := xmldom.Must(xmldom.ParseXML(ExampleXml)).Root
	node := root.FindOneByName("property")
	fmt.Println(node.XML())
	// Output:
	// <property name="go.version">go1.8.1</property>
}

func ExampleNode_FindByName() {
	root := xmldom.Must(xmldom.ParseXML(ExampleXml)).Root
	nodes := root.FindByName("testcase")
	for _, node := range nodes {
		fmt.Println(node.XML())
	}
	// Output:
	// <testcase classname="go-xmldom" id="ExampleParseXML" time="0.004" />
	// <testcase classname="go-xmldom" id="ExampleParse" time="0.005" />
}

func ExampleNode_Query() {
	node := xmldom.Must(xmldom.ParseXML(ExampleXml)).Root
	// xpath expr: https://github.com/antchfx/xpath

	// find all children
	fmt.Printf("children = %v\n", len(node.Query("//*")))

	// find node matched tag name
	nodeList := node.Query("//testcase")
	for _, c := range nodeList {
		fmt.Printf("%v: id = %v\n", c.Name, c.GetAttributeValue("id"))
	}
	// Output:
	// children = 5
	// testcase: id = ExampleParseXML
	// testcase: id = ExampleParse
}

func ExampleNode_QueryOne() {
	node := xmldom.Must(xmldom.ParseXML(ExampleXml)).Root
	// xpath expr: https://github.com/antchfx/xpath

	// find node matched attr name
	c := node.QueryOne("//testcase[@id='ExampleParseXML']")
	fmt.Printf("%v: id = %v\n", c.Name, c.GetAttributeValue("id"))
	// Output:
	// testcase: id = ExampleParseXML
}

func ExampleDocument_XML() {
	doc := xmldom.Must(xmldom.ParseXML(ExampleXml))
	fmt.Println(doc.XML())
	// Output:
	// <?xml version="1.0" encoding="UTF-8"?><!DOCTYPE junit SYSTEM "junit-result.dtd"><testsuites><testsuite tests="2" failures="0" time="0.009" name="github.com/subchen/go-xmldom"><properties><property name="go.version">go1.8.1</property></properties><testcase classname="go-xmldom" id="ExampleParseXML" time="0.004" /><testcase classname="go-xmldom" id="ExampleParse" time="0.005" /></testsuite></testsuites>
}

func ExampleDocument_XMLPretty() {
	doc := xmldom.Must(xmldom.ParseXML(ExampleXml))
	fmt.Println(doc.XMLPretty())
	// Output:
	// <?xml version="1.0" encoding="UTF-8"?>
	// <!DOCTYPE junit SYSTEM "junit-result.dtd">
	// <testsuites>
	//   <testsuite tests="2" failures="0" time="0.009" name="github.com/subchen/go-xmldom">
	//     <properties>
	//       <property name="go.version">go1.8.1</property>
	//     </properties>
	//     <testcase classname="go-xmldom" id="ExampleParseXML" time="0.004" />
	//     <testcase classname="go-xmldom" id="ExampleParse" time="0.005" />
	//   </testsuite>
	// </testsuites>
}

func ExampleNewDocument() {
	doc := xmldom.NewDocument("testsuites")

	testsuiteNode := doc.Root.CreateNode("testsuite").SetAttributeValue("name", "github.com/subchen/go-xmldom")
	testsuiteNode.CreateNode("testcase").SetAttributeValue("name", "case 1").Text = "PASS"
	testsuiteNode.CreateNode("testcase").SetAttributeValue("name", "case 2").Text = "FAIL"

	fmt.Println(doc.XMLPretty())
	// Output:
	// <?xml version="1.0" encoding="UTF-8"?>
	// <testsuites>
	//   <testsuite name="github.com/subchen/go-xmldom">
	//     <testcase name="case 1">PASS</testcase>
	//     <testcase name="case 2">FAIL</testcase>
	//   </testsuite>
	// </testsuites>
}

func TestAttrNamespace(t *testing.T) {
	root := xmldom.Must(xmldom.ParseXML(ExampleXml)).Root
	node := root.FindByID("AttrNamespace")

	if node.Attributes[0].Name != "xmlns:test" {
		t.Fatalf("Expected attribute name to be xmlns:test, got=%s", node.Attributes[0].Name)
	}
}
