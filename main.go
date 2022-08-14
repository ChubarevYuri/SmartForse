package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

func main() {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err.Error())
	}
	oldStr, erro := os.LookupEnv("old")
	if !erro {
		fmt.Println("err old env path")
		panic("not env old")
	}
	fmt.Println("[old] = " + oldStr)
	newStr, errn := os.LookupEnv("new")
	if !errn {
		fmt.Println("err new env path")
		panic("not env new")
	}
	fmt.Println("[new] = " + newStr)
	if (len(oldStr) > 0) && (oldStr != newStr) {
		fmt.Println("change start")
		ChangeTxtToPath(pwd+"/files", oldStr, newStr)
	} else {
		fmt.Println("none change params")
	}
}

func ChangeTxtToPath(path, oldStr, newStr string) {
	logfile := time.Now().Format("2006-01-02T15-04-05")
	paths := listDirByReadDir(path)
	fmt.Println(path)
	fmt.Println(fmt.Sprintf("files count = %d", len(paths)))
	for i := range paths {
		fmt.Println(paths[i])
		f, err := os.Open(paths[i])
		if err == nil {
			wr := bytes.Buffer{}
			sc := bufio.NewScanner(f)
			for sc.Scan() {
				wr.WriteString(sc.Text())
			}
			fileTxt := wr.String()
			fmt.Println("File " + paths[i] + " read: " + fileTxt)
			f.Close()
			if strings.Contains(fileTxt, oldStr) {
				if os.Remove(paths[i]) != nil {
					logWrite("File "+paths[i]+" don`t delete\n", logfile)
				}
				_, err := os.Create(paths[i])
				if err != nil {
					logWrite("File "+paths[i]+" don`t create\n", logfile)
				} else {
					f, err := os.OpenFile(paths[i], os.O_WRONLY, 0600)
					if err == nil {
						_, err = f.WriteString(strings.Replace(fileTxt, oldStr, newStr, -1))
						if err != nil {
							logWrite("File "+paths[i]+" don`t write\n", logfile)
							fmt.Println("File " + paths[i] + " don`t save")
						} else {
							logWrite(paths[i]+"\n"+logText(fileTxt, oldStr, newStr), logfile)
							fmt.Println("File " + paths[i] + " save")
						}
					} else {
						logWrite("File "+paths[i]+" don`t write\n", logfile)
						fmt.Println("File " + paths[i] + " don`t save")
					}
					f.Close()
				}
			}
		} else {
			logWrite("File "+paths[i]+" don`t read\n", logfile)
			fmt.Println("File " + paths[i] + " don`t read")
		}
		f.Close()
	}
}

func listDirByReadDir(path string) []string {
	var lst, err = ioutil.ReadDir(path)
	var result []string
	if err != nil {
		return nil
	}
	for _, val := range lst {
		if !val.IsDir() {
			result = append(result, path+"/"+val.Name())
		}
	}
	return result
}

func logWrite(text, fileName string) {
	pwd, err := os.Getwd()
	if err == nil {
		os.MkdirAll(pwd+"/log", 0777)
		path := pwd + "/log/" + fileName + ".log"
		if _, err := os.Stat(path); err != nil {
			os.Create(path)
		}
		log, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0600)
		if err == nil {
			log.WriteString(text)
		}
		log.Close()
	}
}

func logText(fileTxt, oldStr, newStr string) string {
	var mytxt = "[" + oldStr + "->" + newStr + "]"
	result := ""
	for true {
		i := strings.Index(fileTxt, oldStr)
		if i == -1 {
			break
		} else {
			if i > 10 {
				result += fileTxt[i-10 : i]
			} else if i > 0 {
				result += fileTxt[:i]
			}
			result += mytxt
			if i+len(oldStr)+10 < len(fileTxt) {
				result += fileTxt[i+len(oldStr) : i+len(oldStr)+10]
			} else if i+len(oldStr) < len(fileTxt) {
				result += fileTxt[i+len(oldStr):]
			}
			result += "\n"
			fileTxt = strings.Replace(fileTxt, oldStr, "", 1)
		}
	}
	return result
}
