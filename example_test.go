package xmldom

import (
	"fmt"
	"bytes"
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
	</testsuite>
</testsuites>`

	ExampleNamespaceXml = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE junit SYSTEM "junit-result.dtd">
<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
	<S:Body>
		<ns0:Content xmlns:ns0="namespace_0" xmlns:ns1="namespace_1">
			<ns1:Item>item1</ns1:Item>
			<ns1:Item>item2</ns1:Item>
		</ns0:Content>
		<ns2:Other xmlns:ns2="namespace_2" param="test_param_value"/>
	</S:Body>
</S:Envelope>`

	ExampleInheritNamespaceXML = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE junit SYSTEM "junit-result.dtd">
<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
	<S:Body>
		<ds:Signature xmlns:ds="http://www.w3.org/2000/09/xmldsig#">
			<ds:SignedInfo>
				<ds:DigestValue>KHDHFKSFH2IEFIDHJKSHFKJHSKDJFH==
				</ds:DigestValue>
			</ds:SignedInfo>
		</ds:Signature>
	</S:Body>
</S:Envelope>`
)

func ExampleParseXML() {
	node := Must(ParseXML(ExampleXml)).Root
	fmt.Printf("name = %v\n", node.Name.Local)
	fmt.Printf("attributes.len = %v\n", len(node.Attributes))
	fmt.Printf("children.len = %v\n", len(node.Children))
	fmt.Printf("root = %v", node == node.Root())
	// Output:
	// name = testsuites
	// attributes.len = 0
	// children.len = 1
	// root = true
}

func ExampleParseNamespacesXML() {
	node := Must(ParseXML(ExampleNamespaceXml)).Root
	fmt.Printf("name.Local = %v\n", node.Name.Local)
	fmt.Printf("name.Space = %v\n", node.Name.Space)
	fmt.Printf("attributes.len = %v\n", len(node.Attributes))
	fmt.Printf("children.len = %v\n", len(node.Children))
	fmt.Printf("root = %v", node == node.Root())
	// Output:
	// name.Local = Envelope
	// name.Space = http://schemas.xmlsoap.org/soap/envelope/
	// attributes.len = 1
	// children.len = 1
	// root = true
}

func ExampleEmptyElementTag() {
	doc := Must(ParseXML(ExampleNamespaceXml))
	doc.EmptyElementTag = true
	fmt.Println(doc.Root.XML())
	// Output:
	// <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/"><S:Body><ns0:Content xmlns:ns0="namespace_0"><ns1:Item>item1</ns1:Item><ns1:Item>item2</ns1:Item></ns0:Content><ns2:Other xmlns:ns2="namespace_2" param="test_param_value" /></S:Body></S:Envelope>
}

func ExampleStartTagEndTag() {
	doc := Must(ParseXML(ExampleNamespaceXml))
	doc.EmptyElementTag = false
	fmt.Println(doc.Root.XML())
	// Output:
	// <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/"><S:Body><ns0:Content xmlns:ns0="namespace_0"><ns1:Item>item1</ns1:Item><ns1:Item>item2</ns1:Item></ns0:Content><ns2:Other xmlns:ns2="namespace_2" param="test_param_value"></ns2:Other></S:Body></S:Envelope>
}

func ExampleNode_GetAttribute() {
	node := Must(ParseXML(ExampleXml)).Root
	attr := node.FirstChild().GetAttribute("name")
	fmt.Printf("%v = %v\n", attr.Name.Local, attr.Value)
	// Output:
	// name = github.com/subchen/go-xmldom
}

func ExampleNode_GetChildren() {
	node := Must(ParseXML(ExampleXml)).Root
	children := node.FirstChild().GetChildren("testcase")
	for _, c := range children {
		fmt.Printf("%v: id = %v\n", c.Name.Local, c.GetAttributeValue("id"))
	}
	// Output:
	// testcase: id = ExampleParseXML
	// testcase: id = ExampleParse
}

func ExampleNode_FindByID() {
	root := Must(ParseXML(ExampleXml)).Root
	node := root.FindByID("ExampleParseXML")
	fmt.Println(node.XML())
	// Output:
	// <testcase classname="go-xmldom" id="ExampleParseXML" time="0.004" />
}

func ExampleNode_FindOneByName() {
	root := Must(ParseXML(ExampleXml)).Root
	node := root.FindOneByName("property")
	fmt.Println(node.XML())
	// Output:
	// <property name="go.version">go1.8.1</property>
}

func ExampleNode_FindByName() {
	root := Must(ParseXML(ExampleXml)).Root
	nodes := root.FindByName("testcase")
	for _, node := range nodes {
		fmt.Println(node.XML())
	}
	// Output:
	// <testcase classname="go-xmldom" id="ExampleParseXML" time="0.004" />
	// <testcase classname="go-xmldom" id="ExampleParse" time="0.005" />
}

func ExampleNode_Query() {
	node := Must(ParseXML(ExampleXml)).Root
	// xpath expr: https://github.com/antchfx/xpath

	// find all children
	fmt.Printf("children = %v\n", len(node.Query("//*")))

	// find node matched tag name
	nodeList := node.Query("//testcase")
	for _, c := range nodeList {
		fmt.Printf("%v: id = %v\n", c.Name.Local, c.GetAttributeValue("id"))
	}
	// Output:
	// children = 5
	// testcase: id = ExampleParseXML
	// testcase: id = ExampleParse
}

func ExampleNode_QueryOne() {
	node := Must(ParseXML(ExampleXml)).Root
	// xpath expr: https://github.com/antchfx/xpath

	// find node matched attr name
	c := node.QueryOne("//testcase[@id='ExampleParseXML']")
	fmt.Printf("%v: id = %v\n", c.Name.Local, c.GetAttributeValue("id"))
	// Output:
	// testcase: id = ExampleParseXML
}

func ExampleDocument_XML() {
	doc := Must(ParseXML(ExampleXml))
	fmt.Println(doc.XML())
	// Output:
	// <?xml version="1.0" encoding="UTF-8"?><!DOCTYPE junit SYSTEM "junit-result.dtd"><testsuites><testsuite tests="2" failures="0" time="0.009" name="github.com/subchen/go-xmldom"><properties><property name="go.version">go1.8.1</property></properties><testcase classname="go-xmldom" id="ExampleParseXML" time="0.004" /><testcase classname="go-xmldom" id="ExampleParse" time="0.005" /></testsuite></testsuites>
}

func ExampleDocument_XMLPretty() {
	doc := Must(ParseXML(ExampleXml))
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
	doc := NewDocument("testsuites")

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

func ExampleInheritNamespace() {
	doc := NewDocument("")
	doc.EmptyElementTag = false
	doc.TextSafeMode = false
	err := doc.Parse(
		bytes.NewReader(
			[]byte(ExampleInheritNamespaceXML),
		),
	)
	fmt.Println(err)
	fmt.Println(doc.Root.XML())

	node := doc.Root.QueryOne("//Body/Signature/SignedInfo")
	fmt.Println(node.XML())
	// Output:
	// <nil>
	// <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/"><S:Body><ds:Signature xmlns:ds="http://www.w3.org/2000/09/xmldsig#"><ds:SignedInfo><ds:DigestValue>KHDHFKSFH2IEFIDHJKSHFKJHSKDJFH==
	// 				</ds:DigestValue></ds:SignedInfo></ds:Signature></S:Body></S:Envelope>
	// <ds:SignedInfo xmlns:ds="http://www.w3.org/2000/09/xmldsig#"><ds:DigestValue>KHDHFKSFH2IEFIDHJKSHFKJHSKDJFH==
	// 				</ds:DigestValue></ds:SignedInfo>
}