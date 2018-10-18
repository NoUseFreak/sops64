package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"gopkg.in/yaml.v2"
)

func main() {

	args := os.Args[1 : len(os.Args)-1]
	file := os.Args[len(os.Args)-1]

	var action string
	sopsArgs := args[:0]
	for _, v := range args {
		switch v {
		case "-d":
			action = "--decrypt"
		case "--decrypt":
			action = v
		case "-e":
			action = "--encrypt"
		case "--encrypt":
			action = v
		default:
			sopsArgs = append(sopsArgs, v)
		}
	}

	switch action {
	case "--decrypt":
		fmt.Print(decrypt(file, sopsArgs))
	case "--encrypt":
		fmt.Print(encrypt(file, sopsArgs))
	default:
		log.Fatal("Unknown action")
	}
}

func encrypt(file string, sopsArgs []string) string {
	yamlFile, err := ioutil.ReadFile(file)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	c := make(map[interface{}]interface{})
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	switch data := c["data"].(type) {
	case map[interface{}]interface{}:
		c["data"] = b64enc(data)
	}

	content := getYaml(c)
	tmpfile, err := ioutil.TempFile("", "sops64")
	if err != nil {
		log.Fatal(err)
	}

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		log.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		log.Fatal(err)
	}

	os.Rename(tmpfile.Name(), tmpfile.Name()+".yml")
	defer os.Remove(tmpfile.Name() + ".yml")

	cmd := fmt.Sprintf("sops --encrypt %s %s", strings.Join(sopsArgs, " "), tmpfile.Name()+".yml")
	out, err := exec.Command("bash", "-c", cmd).CombinedOutput()
	if err != nil {
		fmt.Print(string(out))
		os.Exit(1)
	}

	return string(out)
}

func decrypt(file string, sopsArgs []string) string {
	cmd := fmt.Sprintf("sops --decrypt %s %s", strings.Join(sopsArgs, " "), file)
	out, err := exec.Command("bash", "-c", cmd).CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}

	c := make(map[interface{}]interface{})
	err = yaml.Unmarshal([]byte(out), c)
	if err != nil {
		fmt.Print(string(out))
		os.Exit(1)
	}

	switch data := c["data"].(type) {
	case map[interface{}]interface{}:
		c["data"] = b64dec(data)
	}

	return getYaml(c)
}

func getYaml(c map[interface{}]interface{}) string {
	d, err := yaml.Marshal(&c)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return string(d)

}

func b64dec(c map[interface{}]interface{}) map[interface{}]interface{} {
	for k, v := range c {
		switch v := v.(type) {
		case string:
			bla, _ := base64.StdEncoding.DecodeString(v)
			c[k] = string(bla)
		case map[interface{}]interface{}:
			c[k] = b64dec(v)
		}
	}

	return c
}

func b64enc(c map[interface{}]interface{}) map[interface{}]interface{} {
	for k, v := range c {
		switch v := v.(type) {
		case string:
			c[k] = base64.StdEncoding.EncodeToString([]byte(v))
		case map[interface{}]interface{}:
			c[k] = b64enc(v)
		default:
			fmt.Printf("I don't know about type %T!\n", v)
		}
	}

	return c
}
