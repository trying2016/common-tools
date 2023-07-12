//go:build windows
// +build windows

package process

import (
	"bufio"
	"github.com/sirupsen/logrus"
	"github.com/trying2016/common-tools/log"
	"github.com/trying2016/common-tools/utils"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"x-proxy/pkg/cpu"
)

const (
	checkInterval = 1000
)

func NewProcess(filePath string, args []string, env []string, hide bool) *Process {
	return &Process{filePath: filePath, hide: hide, args: args, evn: env}
}

type Process struct {
	filePath      string
	process       *os.Process
	job           sync.WaitGroup
	hide          bool
	args          []string
	evn           []string
	checkTimer    chan struct{}
	affinityStart int
	affinityStep  int
	affinityCount int
}

func (p *Process) check() {
	if p.process == nil {
		_ = p.exec()
	}
}

func (p *Process) Kill() error {
	if p.process != nil {
		p.process.Kill()
	}
	return nil
}

func (p *Process) SetAffinity(start, step, count int) {
	p.affinityStart = start
	p.affinityStep = step
	p.affinityCount = count
}
func (p *Process) Run(restartTime int) error {
	p.Stop()
	err := p.exec()
	if err != nil {
		return err
	}
	p.checkTimer = utils.StartTime(p.check, checkInterval)
	if restartTime != 0 {
		utils.StartTime(func() {
			_ = p.Kill()
		}, restartTime)
	}
	return nil
}

func (p *Process) exec() error {
	fileDir := filepath.Dir(p.filePath)
	args := make([]string, 0)
	if runtime.GOOS == "windows" {
		//args = append(args, p.filePath)
		args = append(args, p.args...)
	} else {
		args = append(args, p.args...)
	}
	cmd := exec.Command(p.filePath, args...)
	//cmd := exec.Command(p.filePath)
	cmd.Dir = fileDir
	//cmd.Env = p.evn

	//defaultEvn,_ := execenv.Default(cmd.SysProcAttr)
	cmd.Env = append(syscall.Environ(), p.evn...)
	cmdStdout, err := cmd.StdoutPipe()
	if err != nil {
		logrus.Errorf("run fail %v", p.filePath)
		return err
	}
	cmd.SysProcAttr = NewSysProcAttr(p.hide)
	cmdStderr, err := cmd.StderrPipe()
	if err != nil {
		logrus.Errorf("run fail, error %v", err)
		return err
	}
	if err := cmd.Start(); err != nil {
		logrus.Errorf("run fail, error %v", err)
		return err
	}
	p.process = cmd.Process
	// 绑定cpu
	if p.affinityCount != 0 {
		cpu.SetProcessAffinity(p.affinityStart, p.affinityStep, p.affinityCount, p.process.Pid)
	}

	cancel := make(chan string)
	p.job.Add(2)
	go func() {
		defer p.job.Done()
		reader := bufio.NewReader(cmdStdout)
		for {
			select {
			case <-cancel:
				return
			default:
			}
			line, err := reader.ReadString('\n')
			if err == io.EOF {
				break
			}
			line = strings.ReplaceAll(line, "\n", "")
			line = strings.ReplaceAll(line, "\r", "")
			//logrus.Info(line)
			if strings.Contains(strings.ToLower(line), "error") ||
				strings.Contains(strings.ToLower(line), "fail") ||
				strings.Contains(strings.ToLower(line), "----") {
				if !strings.Contains(line, "Failed to connect to operator") {
					if strings.Contains(line, "----") {
						log.Error(strings.ReplaceAll(line, "---- ", ""))
					} else {
						log.Error(line)
					}
				}
			}
		}
	}()
	go func() {
		defer p.job.Done()
		reader := bufio.NewReader(cmdStderr)
		for {
			select {
			case <-cancel:
				return
			default:
			}
			line, err := reader.ReadString('\n')
			if err == io.EOF {
				break
			}
			line = strings.ReplaceAll(line, "\n", "")
			line = strings.ReplaceAll(line, "\r", "")
			//			logrus.Error(line)
			if strings.Contains(strings.ToLower(line), "error") ||
				strings.Contains(strings.ToLower(line), "fail") ||
				strings.Contains(strings.ToLower(line), "----") {
				if !strings.Contains(line, "Failed to connect to operator") {
					if strings.Contains(line, "----") {
						log.Error(strings.ReplaceAll(line, "---- ", ""))
					} else {
						//fmt.Println(line)
						log.Error(line)
					}
				}
			}
		}
	}()

	go func(pCmd *exec.Cmd) {
		if err := pCmd.Wait(); err != nil {
			logrus.Errorf("exit %v %v", p.filePath, err.Error())
		}
		close(cancel)
		if p.process != nil {
			p.process.Kill()
			p.process = nil
		}
		p.job.Wait()
	}(cmd)
	return nil
}
func (p *Process) Stop() {
	if p.process != nil {
		_ = p.process.Kill()
		p.process = nil
		p.job.Wait()
	}
	if p.checkTimer != nil {
		utils.StopTime(p.checkTimer)
		p.checkTimer = nil
	}
}

func (p *Process) IsRun() bool {
	return p.process != nil
}
func (p *Process) Filename() string {
	return p.filePath
}
