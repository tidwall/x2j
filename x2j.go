package x2j

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io"
	"strings"
)

type xnode struct {
	parent   *xnode
	xelem    xml.StartElement
	xdata    string
	children []*xnode
}

func appendJSONString(dst []byte, str string) []byte {
	data, _ := json.Marshal(str)
	return append(dst, data...)
}

func appendNode(dst []byte, node *xnode) []byte {
	if node.xelem.Name.Local == "" {
		return appendJSONString(dst, node.xdata)
	}
	dst = append(dst, '{')
	dst = appendJSONString(dst, "name")
	dst = append(dst, ':')
	dst = appendJSONString(dst, node.xelem.Name.Local)
	if len(node.xelem.Attr) > 0 {
		dst = append(dst, ',')
		dst = appendJSONString(dst, "attrs")
		dst = append(dst, ':')
		dst = append(dst, '{')
		for i, attr := range node.xelem.Attr {
			if i > 0 {
				dst = append(dst, ',')
			}
			dst = appendJSONString(dst, attr.Name.Local)
			dst = append(dst, ':')
			dst = appendJSONString(dst, attr.Value)
		}
		dst = append(dst, '}')
	}
	if len(node.children) > 0 {
		dst = append(dst, ',')
		dst = appendJSONString(dst, "children")
		dst = append(dst, ':')
		dst = append(dst, '[')
		for i, child := range node.children {
			if i > 0 {
				dst = append(dst, ',')
			}
			dst = appendNode(dst, child)
		}
		dst = append(dst, ']')
	}
	return append(dst, '}')
}

func Convert(xmldata []byte) ([]byte, error) {
	dec := xml.NewDecoder(bytes.NewReader(xmldata))
	var node *xnode
	for {
		tok, err := dec.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		switch tok := tok.(type) {
		case xml.StartElement:
			node = &xnode{parent: node, xelem: tok}
		case xml.CharData:
			cdata := string(xml.CharData(tok))
			nlws := len(cdata) > 0 && (cdata[0] == '\r' || cdata[0] == '\n') &&
				strings.TrimSpace(cdata) == ""
			if !nlws && cdata != "" && node != nil {
				node.children = append(node.children, &xnode{
					parent: node,
					xdata:  cdata,
				})
			}
		case xml.EndElement:
			if node.parent != nil {
				node.parent.children = append(node.parent.children, node)
				node = node.parent
			}
		}
	}
	return appendNode(nil, node), nil
}
