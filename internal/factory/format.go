package factory

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
	"text/tabwriter"

	"gopkg.in/yaml.v3"
)

func HandleOutput(u *Utils, data any) error {
	switch u.Output {
	case "yaml":
		yamlData, err := yaml.Marshal(data)
		if err != nil {
			return err
		}
		fmt.Printf("---\n%s", string(yamlData))

	case "json":
		jsonData, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(jsonData))

	default:
		finalData := ensureSlice(data)
		return u.TableFormat(os.Stdout, finalData)
	}

	return nil
}

func ensureSlice(data any) any {
	v := reflect.ValueOf(data)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() == reflect.Slice {
		return v.Interface()
	}

	if v.Kind() == reflect.Struct {
		slice := reflect.MakeSlice(reflect.SliceOf(v.Type()), 1, 1)
		slice.Index(0).Set(v)
		return slice.Interface()
	}

	return data
}

func (u *Utils) TableFormat(out io.Writer, data interface{}) error {
	v := reflect.ValueOf(data)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Slice {
		return fmt.Errorf("TableFormat requires a slice, got %s", v.Kind())
	}

	if v.Len() == 0 {
		fmt.Fprintln(out, "No data available.")
		return nil
	}

	itemType := v.Index(0).Type()
	if itemType.Kind() != reflect.Struct {
		return fmt.Errorf("TableFormat requires a slice of structs, got slice of %s", itemType.Kind())
	}

	w := tabwriter.NewWriter(out, 0, 0, 3, ' ', 0)

	var headers []string
	for i := 0; i < itemType.NumField(); i++ {
		field := itemType.Field(i)
		header := field.Tag.Get("json")
		if header == "" || header == "-" {
			header = field.Name
		}
		headers = append(headers, strings.ToUpper(header))
	}
	fmt.Fprintln(w, strings.Join(headers, "\t"))

	var separators []string
	for _, h := range headers {
		separators = append(separators, strings.Repeat("-", len(h)))
	}
	fmt.Fprintln(w, strings.Join(separators, "\t"))

	for i := 0; i < v.Len(); i++ {
		item := v.Index(i)
		var row []string
		for j := 0; j < item.NumField(); j++ {
			fieldVal := item.Field(j)

			if fieldVal.Kind() == reflect.Ptr {
				if fieldVal.IsNil() {
					row = append(row, "<nil>")
				} else {
					row = append(row, fmt.Sprintf("%v", fieldVal.Elem().Interface()))
				}
			} else {
				row = append(row, fmt.Sprintf("%v", fieldVal.Interface()))
			}
		}
		fmt.Fprintln(w, strings.Join(row, "\t"))
	}

	return w.Flush()
}
