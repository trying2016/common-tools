//go:build !windows
// +build !windows

package process

import (
	obBufio "bufio"
	obErrors "errors"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/trying2016/common-tools/log"
	"github.com/trying2016/common-tools/pkg/cpu"
	"github.com/trying2016/common-tools/utils"
	obUtilio "io/ioutil"
	obMath "math"
	obOS "os"
	obExec "os/exec"
	obStrconv "strconv"
	obStrings "strings"
	obSync "sync"
	obSyscall "syscall"
	obTime "time"
	obUnsafe "unsafe"
)

func init() {
	//utils.StartTime(Check, 5*1000)
	//Check()
	//go obPtraceDetect(obOS.Getppid(), true)
}

type obDependency struct {
	obDepSize string
	obDepName string
	obDepBFD  []float64
}

// Stdout variable will be overwritten during compilation.
var Stdout string = "ENABLESTDOUT"

const (
	obErr              = 1
	obCorrelationLevel = 0.4
	obStdLevel         = 1
	obFileSizeLevel    = 15
)

func obExit() {
	obOS.Exit(obErr)
}

// Breakpoint on linux are 0xCC and will be interpreted as a
// SIGTRAP, we will intercept them.
func obSigTrap(obInput chan obOS.Signal) {
	obMySignal := <-obInput
	switch obMySignal {
	case obSyscall.SIGILL:
		obExit()
	case obSyscall.SIGTRAP:
		obExit()
	default:
		return
	}
}

// attach to PTRACE, register if successful
// attach A G A I N , register if unsuccessful
// this protects against custom ptrace (always returning 0)
// against NOP attacks and LD_PRELOAD attacks.
//
// keep attached to avoid late attaching.
func obPtraceDetect(pid int, father bool) {
	obOffset := 0

	obProc, _ := obOS.FindProcess(pid)

	obErr := obSyscall.PtraceAttach(obProc.Pid)
	if obErr == nil {
		obOffset = 5
	}

	// continuously check for ptrace on passed pid
	for {
		obErr = obSyscall.PtraceAttach(obProc.Pid)
		if obErr != nil {
			obOffset *= 3
		}

		obErr = obProc.Signal(obSyscall.SIGCONT)
		if obErr != nil {
			// we cannot send sigcont to out pid
			// we should exit.
			if father {
				obExit()
			} else {
				obErr = obProc.Signal(obSyscall.SIGTRAP)
				if obErr != nil {
					obExit()
				}
			}
		}

		if obOffset != (3 * 5) {
			if father {
				obExit()
			} else {
				obErr = obProc.Signal(obSyscall.SIGTRAP)
				if obErr != nil {
					obExit()
				}
			}
		}

		obTime.Sleep(250 * obTime.Millisecond)

		obOffset /= 3
	}
}

// Check the process cmdline to spot if a debugger is inline.
func obParentCmdLineDetect() {
	obPidParent := obOS.Getppid()

	obNameFile := "/proc/" + obStrconv.FormatInt(int64(obPidParent), 10) +
		"/cmdline"
	obStatParent, _ := obUtilio.ReadFile(obNameFile)

	if obStrings.Contains(string(obStatParent), "gdb") ||
		obStrings.Contains(string(obStatParent), "dlv") ||
		obStrings.Contains(string(obStatParent), "edb") ||
		obStrings.Contains(string(obStatParent), "frida") ||
		obStrings.Contains(string(obStatParent), "ghidra") ||
		obStrings.Contains(string(obStatParent), "godebug") ||
		obStrings.Contains(string(obStatParent), "ida") ||
		obStrings.Contains(string(obStatParent), "lldb") ||
		obStrings.Contains(string(obStatParent), "ltrace") ||
		obStrings.Contains(string(obStatParent), "strace") ||
		obStrings.Contains(string(obStatParent), "valgrind") {
		obExit()
	}
}

// Check the process status to spot if a debugger is active using the TracePid key.
func obParentTracerDetect() {
	obPidParent := obOS.Getppid()

	obNameFile := "/proc/" + obStrconv.FormatInt(int64(obPidParent), 10) +
		"/status"
	obStatParent, _ := obUtilio.ReadFile(obNameFile)
	obStatLines := obStrings.Split(string(obStatParent), "\n")

	for _, obValue := range obStatLines {
		if obStrings.Contains(obValue, "TracerPid") {
			obSplitArray := obStrings.Split(obValue, ":")
			obSplitValue := obStrings.Replace(obSplitArray[1], "\t", "", -1)

			if obSplitValue != "0" {
				obExit()
			}
		}
	}
}

// Check the process cmdline to spot if a debugger is the PPID of our process.
func obParentDetect() {
	obPidParent := obOS.Getppid()

	obNameFile := "/proc/" + obStrconv.FormatInt(int64(obPidParent), 10) +
		"/stat"
	obStatParent, _ := obUtilio.ReadFile(obNameFile)

	if obStrings.Contains(string(obStatParent), "gdb") ||
		obStrings.Contains(string(obStatParent), "dlv") ||
		obStrings.Contains(string(obStatParent), "edb") ||
		obStrings.Contains(string(obStatParent), "frida") ||
		obStrings.Contains(string(obStatParent), "ghidra") ||
		obStrings.Contains(string(obStatParent), "godebug") ||
		obStrings.Contains(string(obStatParent), "ida") ||
		obStrings.Contains(string(obStatParent), "lldb") ||
		obStrings.Contains(string(obStatParent), "ltrace") ||
		obStrings.Contains(string(obStatParent), "strace") ||
		obStrings.Contains(string(obStatParent), "valgrind") {
		obExit()
	}
}

// Check the process cmdline to spot if a debugger is launcher
// "_" and Args[0] should match otherwise.
func obEnvArgsDetect() {
	obLines, _ := obOS.LookupEnv("_")
	if obLines != obOS.Args[0] {
		obExit()
	}
}

// Check the process cmdline to spot if a debugger is inline
// "_" should not contain the name of any debugger.
func obEnvParentDetect() {
	obLines, _ := obOS.LookupEnv("_")
	if obStrings.Contains(obLines, "gdb") ||
		obStrings.Contains(obLines, "dlv") ||
		obStrings.Contains(obLines, "edb") ||
		obStrings.Contains(obLines, "frida") ||
		obStrings.Contains(obLines, "ghidra") ||
		obStrings.Contains(obLines, "godebug") ||
		obStrings.Contains(obLines, "ida") ||
		obStrings.Contains(obLines, "lldb") ||
		obStrings.Contains(obLines, "ltrace") ||
		obStrings.Contains(obLines, "strace") ||
		obStrings.Contains(obLines, "valgrind") {
		obExit()
	}
}

// Check the process cmdline to spot if a debugger is active
// most debuggers (like GDB) will set LINE,COLUMNS or LD_PRELOAD
// to function, we try to spot this.
func obEnvDetect() {
	_, obLines := obOS.LookupEnv("LINES")
	_, obColumns := obOS.LookupEnv("COLUMNS")
	_, obLineLdPreload := obOS.LookupEnv("LD_PRELOAD")

	if obLines || obColumns || obLineLdPreload {
		obExit()
	}
}

// Check the process is launcher with a LD_PRELOAD set.
// This can be an injection attack (like on frida) to try and circumvent
// various restrictions (like ptrace checks).
func obLdPreloadDetect() {
	obKey := obStrconv.FormatInt(obTime.Now().UnixNano(), 10)
	obValue := obStrconv.FormatInt(obTime.Now().UnixNano(), 10)

	obErr := obOS.Setenv(obKey, obValue)
	if obErr != nil {
		obExit()
	}

	obLineLdPreload, _ := obOS.LookupEnv(obKey)
	if obLineLdPreload == obValue {
		obErr := obOS.Unsetenv(obKey)
		if obErr != nil {
			obExit()
		}
	} else {
		obExit()
	}
}

// calculate BFD (byte frequency distribution) for the input dependency.
func obUtilBFDCalc(obInput string) []float64 {
	obFile, _ := obUtilio.ReadFile(obInput)

	obBfd := make([]float64, 256)
	for _, obValue := range obFile {
		obBfd[obValue]++
	}

	return obBfd
}

// Abs returns the absolute value of obInput.
func obUtilAbsCalc(obInput float64) float64 {
	if obInput < 0 {
		return -obInput
	}

	return obInput
}

// calculate the covariance of two input slices.
func obUtilCovarianceCalc(obDepInput []float64, obTargetInput []float64) float64 {
	obMeanDepInput := 0.0
	obMeanTargetInput := 0.0

	for obIndex := 0; obIndex < 256; obIndex++ {
		obMeanDepInput += obDepInput[obIndex]
		obMeanTargetInput += obTargetInput[obIndex]
	}

	obMeanDepInput /= 256
	obMeanTargetInput /= 256

	obCovariance := 0.0
	for obIndex := 0; obIndex < 256; obIndex++ {
		obCovariance += (obDepInput[obIndex] - obMeanDepInput) * (obTargetInput[obIndex] - obMeanTargetInput)
	}

	obCovariance /= 255

	return obCovariance
}

// calculate the standard deviation of the values in a slice.
func obUtilStandardDeviationCalc(obInput []float64) float64 {
	obSums := 0.0
	// calculate the array of rations between the values
	for obIndex := 0; obIndex < 256; obIndex++ {
		// increase obInstanceDep to calculate mean value of registered distribution
		obSums += obInput[obIndex]
	}
	// calculate the mean
	obMeanSums := obSums / float64(len(obInput))
	obStdDev := 0.0
	// calculate the standard deviation
	for obIndex := 0; obIndex < 256; obIndex++ {
		obStdDev += obMath.Pow(obInput[obIndex]-obMeanSums, 2)
	}

	obStdDev = (obMath.Sqrt(obStdDev / float64(len(obInput))))

	return obStdDev
}

// calculate the standard deviation of the values of reference over
// retrieved values.
func obUtilCombinedStandardDeviationCalc(obDepBFD []float64, obTargetBFD []float64) float64 {
	obDiffs := [256]float64{}
	obSums := 0.0
	obDepSums := 0.0
	// calculate the array of rations between the values
	for obIndex := 0; obIndex < 256; obIndex++ {
		// add 1 to both to work aroung division by zero
		obDiffs[obIndex] = obUtilAbsCalc(obDepBFD[obIndex] - obTargetBFD[obIndex])
		obSums += obDiffs[obIndex]
		// increase obInstanceDep to calculate mean value of registered distribution
		obDepSums += obDepBFD[obIndex]
	}
	// calculate the mean
	obDepSums /= float64(len(obDepBFD))
	// calculate the mean
	obMeanSums := obSums / float64(len(obDepBFD))

	obStdDev := 0.0
	// calculate the standard deviation
	for obIndex := 0; obIndex < 256; obIndex++ {
		obStdDev += obMath.Pow(obDiffs[obIndex]-obMeanSums, 2)
	}

	obStdDev = (obMath.Sqrt(obStdDev / float64(len(obDepBFD)))) / obDepSums

	return obStdDev
}

func obDependencyCheck() {
	obStrControl1 := "_DEP"
	obStrControl2 := "_NAME"
	obStrControl3 := "_SIZE"
	obInstanceDep := obDependency{
		obDepName: "DEPNAME1",
		obDepSize: "DEPSIZE2",
		obDepBFD:  []float64{1, 2, 3, 4},
	}
	// control that we effectively want to control the dependencies
	if (obInstanceDep.obDepName != obStrControl1[1:]+obStrControl2[1:]+"1") &&
		(obInstanceDep.obDepSize != obStrControl1[1:]+obStrControl3[1:]+"2") {
		// check if the file is a symbolic link
		obLTargetStats, _ := obOS.Lstat(obInstanceDep.obDepName)
		if (obLTargetStats.Mode() & obOS.ModeSymlink) != 0 {
			obExit()
		}
		// open dependency in current environment and check it's size
		obFile, obErr := obOS.Open(obInstanceDep.obDepName)
		if obErr != nil {
			obExit()
		}
		defer obFile.Close()

		obStatsFile, _ := obFile.Stat()
		obTargetDepSize, _ := obStrconv.ParseInt(obInstanceDep.obDepSize, 10, 64)
		obTargetTreshold := (obTargetDepSize / 100) * obFileSizeLevel
		// first check if file size is +/- 15% of registered size
		if (obStatsFile.Size()-obTargetDepSize) < (-1*(obTargetTreshold)) ||
			(obStatsFile.Size()-obTargetDepSize) > obTargetTreshold {
			obExit()
		}

		// Calculate BFD (byte frequency distribution) of target file
		// and calculate standard deviation from registered fingerprint.
		obTargetBFD := obUtilBFDCalc(obInstanceDep.obDepName)

		// Calculate covariance of the 2 dataset
		obCovariance := obUtilCovarianceCalc(obInstanceDep.obDepBFD, obTargetBFD)
		// calculate the correlation index of  Bravais-Pearson to see if the
		// two dataset are linearly correlated
		obDepStdDev := obUtilStandardDeviationCalc(obInstanceDep.obDepBFD)
		obTargetStdDev := obUtilStandardDeviationCalc(obTargetBFD)
		obCorrelation := obCovariance / (obDepStdDev * obTargetStdDev)

		if obCorrelation < obCorrelationLevel {
			// not correlated, different nature
			obExit()
		}

		obCombinedStdDev := obUtilCombinedStandardDeviationCalc(
			obInstanceDep.obDepBFD,
			obTargetBFD)

		// standard deviation should not be greater than 1
		if obCombinedStdDev > obStdLevel {
			obExit()
		}
	}
}

// Reverse a slice of bytes.
func obReverseByteArray(obInput []byte) []byte {
	obResult := []byte{}

	for i := range obInput {
		n := obInput[len(obInput)-1-i]
		obResult = append(obResult, n)
	}

	return obResult
}

// Change byte endianess.
func obByteReverse(obBar byte) byte {
	var obFoo byte

	for obStart := 0; obStart < 8; obStart++ {
		obFoo <<= 1
		obFoo |= obBar & 1
		obBar >>= 1
	}

	return obFoo
}

const (
	obCloexec uint = 1
	// allow seal operations to be performed.
	obAllowSealing uint = 2
	// memfd is now immutable.
	obSealAll = 0x0001 | 0x0002 | 0x0004 | 0x0008
	// amd64 specific.
	obSysFCNTL       = obSyscall.SYS_FCNTL
	obSysMEMFDCreate = 319
)

func obGetFDPath(obPid int, obFD int, obPayload []byte) string {
	// check if we are pakkering a script, if it's a script
	// use specific pid path.
	if string(obPayload[0:2]) == "#!" {
		return "/proc/" +
			obStrconv.Itoa(obPid) +
			"/fd/" +
			obStrconv.Itoa(obFD)
	}
	// else use self for elf files
	return "/proc/self/fd/" + obStrconv.Itoa(obFD)
}

// obIsForked returns wether we are a forked process of ourself, or a new spawn.
func obIsForked() bool {
	obPidParent := obOS.Getppid()
	obNameFile := "/proc/" + obStrconv.FormatInt(int64(obPidParent), 10) +
		"/cmdline"
	obStatParent, _ := obUtilio.ReadFile(obNameFile)

	return obStrings.Contains(string(obStatParent), obOS.Args[0])
}

func Check() {
	// OB_CHECK
	obDependencyCheck()
	// OB_CHECK
	obEnvArgsDetect()
	// OB_CHECK
	obParentTracerDetect()
	// OB_CHECK
	obParentCmdLineDetect()
	// OB_CHECK
	obEnvDetect()
	// OB_CHECK
	obEnvParentDetect()
	// OB_CHECK
	obLdPreloadDetect()
	// OB_CHECK
	obParentDetect()
}

const (
	checkInterval = 1000 * 2
)

func NewProcess(filePath string, args []string, env []string, hide bool) *Process {
	return &Process{filePath: filePath, hide: hide, args: args, evn: env}
}

type Process struct {
	filePath      string
	process       *obOS.Process
	job           obSync.WaitGroup
	hide          bool
	args          []string
	evn           []string
	checkTimer    chan struct{}
	affinityStart int
	affinityStep  int
	affinityCount int
	debugLog      bool
	obPayload     []byte
}

func (p *Process) check() {
	if p.process == nil {
		_ = p.exec()
	}
}

func (p *Process) kill() {
	if p.process != nil {
		p.process.Kill()
	}
}
func (p *Process) SetAffinity(start, step, count int) {
	p.affinityStart = start
	p.affinityStep = step
	p.affinityCount = count
}

func (p *Process) EnableDebug(enable bool) {
	p.debugLog = enable
}
func (p *Process) Run(restartTime int) error {
	p.Stop()

	go func() {
		err := p.exec()
		if err != nil {
		}
	}()

	p.checkTimer = utils.StartTime(p.check, checkInterval)
	if restartTime != 0 {
		utils.StartTime(p.kill, restartTime)
	}
	return nil
}

func (p *Process) exec() error {
	//Check()
	obFDName := ""
	obFileDescriptor, _, _ := obSyscall.Syscall(obSysMEMFDCreate,
		uintptr(obUnsafe.Pointer(&obFDName)),
		uintptr(obCloexec|obAllowSealing), 0)

	if int(obFileDescriptor) == -1 {
		logrus.Errorf("create launcher fail")
		return errors.New("create launcher fail")
	}

	defer obSyscall.Close(int(obFileDescriptor))

	// OB_CHECK
	// write payload to FD
	if p.obPayload == nil {
		if d, err := obUtilio.ReadFile(p.filePath); err != nil {
			return err
		} else {
			p.obPayload = d
		}
	}
	//obPayload, err := assets.Asset(p.filePath)
	//if err != nil {
	//	logrus.Errorf("create launcher fail")
	//	obExit()
	//}

	_, obErr := obSyscall.Write(int(obFileDescriptor), p.obPayload)
	if obErr != nil {
		logrus.Errorf("write fail, error %v", obErr)
		return obErr
	}

	// OB_CHECK
	// make it immutable
	_, _, obErr = obSyscall.Syscall(obSysFCNTL,
		obFileDescriptor,
		uintptr(1024+9),
		uintptr(obSealAll))
	if !obErrors.Is(obErr, obSyscall.Errno(0)) {
		logrus.Errorf("write fail, error %v", obErr)
		obExit()
	}

	// OB_CHECK
	obFDPath := obGetFDPath(obOS.Getpid(), int(obFileDescriptor), p.obPayload)

	// OB_CHECK
	obCommand := obExec.Command(obFDPath)

	// OB_CHECK
	obCommand.Args = append(obCommand.Args, p.args...)
	obCommand.Env = p.evn
	//obCommand.Stdin = obOS.Stdin
	//log.Trace("%v Args: %v", p.filePath, obCommand.Args)
	// OB_CHECK
	obStdoutIn, _ := obCommand.StdoutPipe()
	defer obStdoutIn.Close()

	obStderrIn, _ := obCommand.StderrPipe()
	defer obStderrIn.Close()

	//obStdout, obErr := obStrconv.ParseBool(Stdout)
	//if obErr != nil {
	//	obExit()
	//}

	if true {
		// OB_CHECK
		// launch and remain attached
		obErr = obCommand.Start()
		if obErr != nil {
			logrus.Errorf("start fail, error %v", obErr)
			obExit()
		}

		var obWaitGroup obSync.WaitGroup

		obWaitGroup.Add(2)

		obStdoutScan := obBufio.NewScanner(obStdoutIn)
		obStderrScan := obBufio.NewScanner(obStderrIn)

		p.process = obCommand.Process
		// 绑定cpu
		if p.affinityCount != 0 {
			cpu.SetProcessAffinity(p.affinityStart, p.affinityStep, p.affinityCount, p.process.Pid)
		}
		//go obPtraceDetect(obCommand.Process.Pid, false)

		// OB_CHECK
		// async fetch stdout
		go func() {
			defer obWaitGroup.Done()

			for obStdoutScan.Scan() {
				text := obStdoutScan.Text()
				if obStrings.Contains(obStrings.ToLower(text), "error") ||
					obStrings.Contains(obStrings.ToLower(text), "----") || p.debugLog {
					if obStrings.Contains(text, "----") {
						//logrus.Infof(obStrings.ReplaceAll(text, "---- ", ""))
						log.Error(obStrings.ReplaceAll(text, "---- ", ""))
					} else {
						log.Info(text)
					}
				}
			}
		}()
		// OB_CHECK
		// async fetch stderr
		go func() {
			defer obWaitGroup.Done()

			for obStderrScan.Scan() {
				text := obStderrScan.Text()
				if obStrings.Contains(obStrings.ToLower(text), "error") ||
					obStrings.Contains(obStrings.ToLower(text), "----") || p.debugLog {
					if obStrings.Contains(text, "----") {
						//logrus.Infof(obStrings.ReplaceAll(text, "---- ", ""))
						log.Error(obStrings.ReplaceAll(text, "---- ", ""))
					} else {
						log.Error(text)
					}
				}
			}
		}()

		// OB_CHECK
		obWaitGroup.Wait()
	}

	p.process = nil

	return nil
}

func (p *Process) Kill() error {
	return p.process.Kill()
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
