package main

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type Addr struct {
	Host string `yaml:"host"`
	Port int64  `yaml:"port"`
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

var (
	host = "smtp.qq.com"
	port = 25
)

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

func getTplOptions() []string {
	dw := getDw()
	dir := filepath.Join(dw, "tpl", "*")
	fmt.Printf("dir: %s\n", dir)
	files, _ := filepath.Glob(dir)
	return files
}

func getTpl(filename string) string {
	wd := getDw()
	content, _ := os.ReadFile(fmt.Sprintf("%s/tpl/%s", wd, filename))
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

	// 获取可执行文件的路径
	execFile, err := os.Executable()
	if err != nil {
		panic(err)
	}

	// 比较调用者文件和可执行文件的路径
	if filepath.Base(callerFile) == "main.go" {
		fmt.Println("程序是通过 go run 命令启动的")
		return true
	} else if filepath.Base(callerFile) == filepath.Base(execFile) {
		fmt.Println("程序是通过编译后的可执行文件启动的")
	} else {
		fmt.Println("无法确定程序是如何启动的")
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
	fmt.Printf("tpls: %s\n", tplOptions)
	//emails := getEmails()
	//for _, to := range emails {
	//	m := gomail.NewMessage()
	//	m.SetHeader("From", c.Email)
	//	m.SetHeader("To", to)
	//	m.SetHeader("Subject", "Subject")
	//	body := getTpl("welcome.html")
	//	m.SetBody("text/html", body)
	//	d := gomail.NewDialer("smtp.qq.com", 25, c.Email, c.Password)
	//
	//	if err := d.DialAndSend(m); err != nil {
	//		fmt.Printf("%s 邮件发送失败\n", to)
	//	}
	//}
}
