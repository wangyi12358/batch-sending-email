package main

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"gopkg.in/gomail.v2"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type Addr struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type TplConfig struct {
	Path     string
	Filename string
	Name     string
}

type Config struct {
	Email    string `yaml:"email"`
	Password string `yaml:"password"`
}

func questionsConfig() (*Config, error) {
	var config Config

	emailQuestion := &survey.Input{
		Message: "请输入你的Email",
	}
	passwordQuestion := &survey.Input{
		Message: "请输入你的密码",
	}
	questions := []*survey.Question{
		{
			Name:   "email",
			Prompt: emailQuestion,
		},
		{
			Name:   "password",
			Prompt: passwordQuestion,
		},
	}
	err := survey.Ask(questions, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func questionTpl(options []TplConfig) *TplConfig {
	var names []string
	for _, o := range options {
		names = append(names, o.Name)
	}
	var name string
	question := &survey.Select{
		Message: "请选择你要发送的模版",
		Options: names,
	}
	_ = survey.AskOne(question, &name)
	for _, o := range options {
		if name == o.Name {
			return &o
		}
	}
	return nil
}

func setEmailConfig(c *Config) {
	out, _ := yaml.Marshal(c)
	wd := getDw()
	fileName := fmt.Sprintf("%s/config.yaml", wd)
	_ = ioutil.WriteFile(fileName, out, 0666)
}

func getEmailConfig() (*Config, error) {
	wd := getDw()
	var config Config
	content, err := os.ReadFile(fmt.Sprintf("%s/config.yaml", wd))
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(content, &config)
	if err != nil {
		return nil, err
	}
	return &config, err
}

func getAddrConfig() *Addr {
	wd := getDw()
	var config Addr
	content, err := os.ReadFile(fmt.Sprintf("%s/addr.yaml", wd))
	if err != nil {
		panic("获取Addr失败")
	}
	err = yaml.Unmarshal(content, &config)
	if err != nil {
		panic("解析Addr失败")
	}
	return &config
}

// 只获取.html文件
func filterHtmlFile(files []string) (filterFiles []string) {
	for _, file := range files {
		if strings.Contains(file, ".html") {
			filterFiles = append(filterFiles, file)
		}
	}
	return filterFiles
}

func getTplOptions() []TplConfig {
	dw := getDw()
	dir := filepath.Join(dw, "tpl", "*")
	files, _ := filepath.Glob(dir)
	files = filterHtmlFile(files)
	var options []TplConfig
	for _, file := range files {
		arr := strings.Split(file, "/")
		filename := arr[len(arr)-1]
		name := strings.Split(filename, ".")[0]
		options = append(options, TplConfig{
			Path:     file,
			Filename: filename,
			Name:     name,
		})
	}
	return options
}

func getTpl(path string) string {
	content, _ := os.ReadFile(path)
	return string(content)
}

func getEmails() []string {
	wd := getDw()
	content, _ := os.ReadFile(fmt.Sprintf("%s/email.txt", wd))
	return strings.Split(string(content), "\n")
}

func getDw() string {
	wd, _ := os.Getwd()
	if isRun() {
		return wd
	}
	return fmt.Sprintf("%s/Desktop/automatic_send_email", wd)
}

func isRun() bool {

	// 获取调用者的文件路径和行号
	_, callerFile, _, _ := runtime.Caller(0)

	// 比较调用者文件和可执行文件的路径
	if filepath.Base(callerFile) == "main.go" {
		return true
	}
	return false
}

func main() {
	c, err := getEmailConfig()
	if err != nil {
		c, err = questionsConfig()
		setEmailConfig(c)
	}
	tplOptions := getTplOptions()
	tpl := questionTpl(tplOptions)
	emails := getEmails()
	addr := getAddrConfig()
	for _, to := range emails {
		m := gomail.NewMessage()
		m.SetHeader("From", c.Email)
		m.SetHeader("To", to)
		m.SetHeader("Subject", tpl.Name)
		body := getTpl(tpl.Path)
		m.SetBody("text/html", body)
		d := gomail.NewDialer(addr.Host, addr.Port, c.Email, c.Password)

		if err := d.DialAndSend(m); err != nil {
			fmt.Printf("%s 邮件发送失败\n", to)
		}
	}
}
