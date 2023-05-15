package main

import (
	"bytes"
	"io"
	"os"
	"strings"

	jsoniter "github.com/json-iterator/go"
	reader "github.com/ledongthuc/pdf"
	"github.com/yaoapp/kun/grpc"
)

// PDF a simple pdf reader plugin
type PDF struct{ grpc.Plugin }

// Exec execute the plugin and return the result
func (pdf *PDF) Exec(method string, args ...interface{}) (*grpc.Response, error) {
	var v interface{}
	var err error

	if len(args) == 0 {
		bytes, err := jsoniter.Marshal(map[string]interface{}{"code": 400, "message": "missing file path"})
		if err != nil {
			return nil, err
		}
		return &grpc.Response{Bytes: bytes, Type: "map"}, nil
	}

	path, ok := args[0].(string)
	if !ok {
		bytes, err := jsoniter.Marshal(map[string]interface{}{"code": 400, "message": "invalid file path"})
		if err != nil {
			return nil, err
		}
		return &grpc.Response{Bytes: bytes, Type: "map"}, nil
	}

	switch strings.ToLower(method) {
	case "text":
		v, err = pdf.Text(path)
		if err != nil {
			bytes, err := jsoniter.Marshal(map[string]interface{}{"code": 404, "message": err.Error()})
			if err != nil {
				return nil, err
			}
			return &grpc.Response{Bytes: bytes, Type: "map"}, nil
		}

		bytes, err := jsoniter.Marshal(v)
		if err != nil {
			return nil, err
		}
		return &grpc.Response{Bytes: bytes, Type: "string"}, nil

	case "content":
		v, err = pdf.Content(path)
		if err != nil {
			bytes, err := jsoniter.Marshal(map[string]interface{}{"code": 404, "message": err.Error()})
			if err != nil {
				return nil, err
			}
			return &grpc.Response{Bytes: bytes, Type: "map"}, nil
		}

		bytes, err := jsoniter.Marshal(v)
		if err != nil {
			return nil, err
		}

		return &grpc.Response{Bytes: bytes, Type: "array"}, nil

	default:
		bytes, err := jsoniter.Marshal(map[string]interface{}{"code": 404, "message": "invalid method"})
		if err != nil {
			return nil, err
		}
		return &grpc.Response{Bytes: bytes, Type: "map"}, nil
	}

}

// Text get the plain text content of the pdf file
func (pdf *PDF) Text(path string) (string, error) {

	f, r, err := reader.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	var buf bytes.Buffer
	b, err := r.GetPlainText()
	if err != nil {
		return "", err
	}
	buf.ReadFrom(b)
	return buf.String(), nil
}

// Content get the content of the pdf file  (including all font and formatting information)
func (pdf *PDF) Content(path string) ([]string, error) {

	f, r, err := reader.Open(path)
	if err != nil {
		return nil, err
	}

	defer f.Close()
	totalPage := r.NumPage()
	result := []string{}
	for pageIndex := 1; pageIndex <= totalPage; pageIndex++ {
		p := r.Page(pageIndex)
		if p.V.IsNull() {
			continue
		}
		content := []byte{}
		rows, _ := p.GetTextByRow()
		for _, row := range rows {
			for _, word := range row.Content {
				content = append(content, word.S...)
			}
		}
		result = append(result, string(content))
	}
	return result, nil
}

func main() {
	var output io.Writer = os.Stdout
	plugin := &PDF{}
	plugin.SetLogger(output, grpc.Trace)
	grpc.Serve(plugin)
}
