package tools

/*
 *  Generic functions to load/save in-memory structures to JSON file
 */

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"os"

	"github.com/gookit/ini/v2"

	//	"github.com/ZhengHe-MD/properties"
	"github.com/ghodss/yaml"
	"github.com/pelletier/go-toml/v2"
)

/*
 *  Load/Save in-memory database (interface) to JSON file
 *  - file: the full path of the JSON file
 *  - data: the generic structure to save
 */
func ReadFromJSONStream(s io.Reader, data interface{}) error {
	txt, err := ioutil.ReadAll(s)
	if err != nil {
		return err
	}
	return json.Unmarshal(txt, data)
}
func ReadFromJSONFile(file string, data interface{}) error {
	var txt []byte
	var err error
	if file == "-" {
		return ReadFromJSONStream(os.Stdin, data)
	}
	txt, err = ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	return json.Unmarshal(txt, data)
}
func ReadFromJSON(txt []byte, data interface{}) error {
	return json.Unmarshal(txt, data)
}
func WriteToJSON(data interface{}) ([]byte, error) {
	return json.Marshal(data)
}
func WriteToJSONFile(file string, data interface{}) error {
	txt, err := json.Marshal(data)
	if err != nil {
		return errors.New("Can not marshal datas")
	}
	if file == "-" {
		_, err := os.Stdout.Write(txt)
		return err
	}
	return ioutil.WriteFile(file, txt, 0644)
}
func WriteToJSONStream(file *os.File, data interface{}) (int, error) {
	txt, err := json.Marshal(data)
	if err != nil {
		return 0, errors.New("Can not marshal datas")
	}
	return file.Write(txt)
}

/*
 *  Load/Save in-memory database (insterface) to YAML file
 *  - file: the full path of the YAML file
 *  - data: the generic structure to save
 */
func ReadFromYAMLStream(s io.Reader, data interface{}) error {
	txt, err := ioutil.ReadAll(s)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(txt, data)
}
func ReadFromYAMLFile(file string, data interface{}) error {
	var txt []byte
	var err error
	if file == "-" {
		return ReadFromYAMLStream(os.Stdin, data)
	}
	txt, err = ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(txt, data)
}
func ReadFromYAML(txt []byte, data interface{}) error {
	return yaml.Unmarshal(txt, data)
}
func WriteToYAML(data interface{}) ([]byte, error) {
	return yaml.Marshal(data)
}
func WriteToYAMLFile(file string, data interface{}) error {
	txt, err := yaml.Marshal(data)
	if err != nil {
		return errors.New("Can not marshal datas")
	}
	if file == "-" {
		_, err := os.Stdout.Write(txt)
		return err
	}
	return ioutil.WriteFile(file, txt, 0644)
}
func WriteToYAMLStream(file *os.File, data interface{}) (int, error) {
	txt, err := yaml.Marshal(data)
	if err != nil {
		return 0, errors.New("Can not marshal datas")
	}
	return file.Write(txt)
}

/*
 *  Load/Save in-memory database (insterface) to TOML file
 *  - file: the full path of the TOML file
 *  - data: the generic structure to save
 */
func ReadFromTOMLStream(s io.Reader, data interface{}) error {
	txt, err := ioutil.ReadAll(s)
	if err != nil {
		return err
	}
	return toml.Unmarshal(txt, data)
}
func ReadFromTOMLFile(file string, data interface{}) error {
	if file == "-" {
		return ReadFromTOMLStream(os.Stdin, data)
	}
	var txt []byte
	var err error
	txt, err = ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	return toml.Unmarshal(txt, data)
}
func ReadFromTOML(txt []byte, data interface{}) error {
	return toml.Unmarshal(txt, data)
}
func WriteToTOML(data interface{}) ([]byte, error) {
	return toml.Marshal(data)
}
func WriteToTOMLFile(file string, data interface{}) error {
	txt, err := toml.Marshal(data)
	if err != nil {
		return errors.New("Can not marshal datas")
	}
	if file == "-" {
		_, err := os.Stdout.Write(txt)
		return err
	}
	return ioutil.WriteFile(file, txt, 0644)
}
func WriteToTOMLStream(file *os.File, data interface{}) (int, error) {
	txt, err := toml.Marshal(data)
	if err != nil {
		return 0, errors.New("Can not marshal datas")
	}
	return file.Write(txt)
}

/*
 *  Load/Save in-memory database (insterface) to INI file
 *  - file: the full path of the INI file
 *  - data: the generic structure to save
 */
func ReadFromINIStream(s io.Reader, data interface{}) error {
	txt, err := ioutil.ReadAll(s)
	if err != nil {
		return err
	}
	return ReadFromINI(txt, data)
}
func ReadFromINIFile(file string, data interface{}) error {
	if file == "-" {
		return ReadFromINIStream(os.Stdin, data)
	}
	c := ini.New()
	if err := c.LoadExists(file); err != nil {
		return err
	}
	return c.MapTo(&data)
}
func ReadFromINI(txt []byte, data interface{}) error {
	c := ini.New()
	if err := c.LoadStrings(string(txt)); err != nil {
		return err
	}
	return c.MapTo(&data)
}

/*
func WriteToINI(data interface{}) ([]byte, error) {
	return iniv1.Marshal(data)
}
func WriteToINIFile(file string, data interface{}) error {
	txt, err := iniv1.Marshal(data)
	if err != nil { return errors.New("Can not marshal datas") }
	if file=="-" { _,err := os.Stdout.Write(txt) ; return err }
	return ioutil.WriteFile(file, txt, 0644)
}
func WriteToINIStream(file *os.File, data interface{}) (int, error) {
	txt, err := iniv1.Marshal(data)
	if err != nil { return 0,errors.New("Can not marshal datas") }
	return file.Write(txt)
}
*/

/*
 *  Load/Save in-memory database (insterface) to properties file
 *  - file: the full path of the properties file
 *  - data: the generic structure to save
 */

/*
func ReadFromPROPSStream(s io.Reader, data interface{}) error {
	txt, err := ioutil.ReadAll(s)
	if err!=nil { return err }
	return properties.Unmarshal(txt, data)
}
func ReadFromPROPSFile(file string, data interface{}) error {
	var txt []byte
	var err error
	if file == "-" { return ReadFromPROPSStream( os.Stdin, data) }
	txt, err = ioutil.ReadFile(file)
	if err!=nil { return err }
	return properties.Unmarshal(txt, data)
}
func ReadFromPROPS(txt []byte, data interface{}) error {
	return properties.Unmarshal(txt, data)
}
func WriteToPROPS(data interface{}) ([]byte, error) {
	return properties.Marshal(data)
}
func WriteToPROPSFile(file string, data interface{}) error {
	txt, err := properties.Marshal(data)
	if err != nil { return errors.New("Can not marshal datas") }
	if file=="-" { _,err := os.Stdout.Write(txt) ; return err }
	return ioutil.WriteFile(file, txt, 0644)
}
func WriteToPROPSStream(file *os.File, data interface{}) (int, error) {
	txt, err := properties.Marshal(data)
	if err != nil { return 0,errors.New("Can not marshal datas") }
	return file.Write(txt)
}
*/

/*
 *  Load datas from stream into in-memory database
 *  - file: the full path of the JSON file, if not JSON try YAML, then try TOML, then try INI
 *  - data: the generic structure to load
 */
func ReadFromAllStream(s io.Reader, data interface{}) error {
	txt, _ := ioutil.ReadAll(s)
	var err error = nil
	if err = json.Unmarshal(txt, data); err != nil {
		if err = yaml.Unmarshal(txt, data); err != nil {
			if err = toml.Unmarshal(txt, data); err != nil {
				if err = ReadFromINI(txt, data); err != nil {
					return errors.New("Unable to read datas")
				}
			}
		}
	}
	return nil
}

/*
 *  Load datas from file into in-memory database
 *  - file: the full path of the JSON file, if not JSON try YAML, then try TOML, then try INI
 *  - data: the generic structure to load
 */
func ReadFromAllFile(file string, data interface{}) error {
	txt, _ := ioutil.ReadFile(file)
	var err error = nil
	if err = json.Unmarshal(txt, data); err != nil {
		if err = yaml.Unmarshal(txt, data); err != nil {
			if err = toml.Unmarshal(txt, data); err != nil {
				if err = ReadFromINI(txt, data); err != nil {
					return errors.New("Unable to read datas")
				}
			}
		}
	}
	return nil
}

/*
 *  Load datas from []byte into in-memory database
 *  - txt: the data, if not JSON try YAML, then try TOML, then try properties, then try INI
 *  - data: the generic structure to load
 */
func ReadFromAll(txt []byte, data interface{}) error {
	var err error = nil
	if err = json.Unmarshal(txt, data); err != nil {
		if err = yaml.Unmarshal(txt, data); err != nil {
			if err = toml.Unmarshal(txt, data); err != nil {
				if err = ReadFromINI(txt, data); err != nil {
					return errors.New("Unable to read datas")
				}
			}
		}
	}
	return err
}
