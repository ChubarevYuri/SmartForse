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
	fmt.Print("Directory path: ")
	var path string
	fmt.Scanln(&path)
	fmt.Print("Old text: ")
	var oldStr string
	fmt.Scanln(&oldStr)
	fmt.Print("New text: ")
	var newStr string
	fmt.Scanln(&newStr)
	changeTxtToPath(path, oldStr, newStr)
	fmt.Println("Finish.")
	fmt.Scanln()
}

func changeTxtToPath(path, oldStr, newStr string) {
	logfile := time.Now().Format("2006-01-02T15-04-05")
	paths := listDirByReadDir(path)
	for i := range paths {
		f, err := os.Open(paths[i])
		if err == nil {
			wr := bytes.Buffer{}
			sc := bufio.NewScanner(f)
			for sc.Scan() {
				wr.WriteString(sc.Text())
			}
			fileTxt := wr.String()
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
						} else {
							logWrite(paths[i]+"\n"+logText(fileTxt, oldStr, newStr), logfile)
						}
					} else {
						logWrite("File "+paths[i]+" don`t write\n", logfile)
					}
				}
			}
		} else {
			logWrite("File "+paths[i]+" don`t read\n", logfile)
		}
		f.Close()
	}
}

func listDirByReadDir(path string) []string {
	lst, err := ioutil.ReadDir(path)
	var result []string
	if err != nil {
		return nil
	}
	for _, val := range lst {
		if !val.IsDir() {
			result = append(result, path+"\\"+val.Name())
		}
	}
	return result
}

func logWrite(text, fileName string) {
	pwd, err := os.Getwd()
	if err == nil {
		os.MkdirAll(pwd+"\\log", 0777)
		path := pwd + "\\log\\" + fileName + ".log"
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
	changetxt := "[" + oldStr + "->" + newStr + "]"
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
			result += changetxt
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
