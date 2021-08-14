package main

import (
	"encoding/json"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/gofrs/flock"
	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/os/gproc"

	"github.com/gorilla/websocket"
	"github.com/tangr/cicdgo/app/model"
	// "github.com/tangr/cicdgo/app/service"
)

var AgentCICD = agentCICD{}

type agentCICD struct{}

var wsUrl = g.Cfg().GetString("agent.WsUrl")
var syncInterval = g.Cfg().GetInt32("agent.SyncInterval")
var dataPathDir = g.Cfg().GetString("agent.DataPathDir")
var jobFlash = g.Cfg().GetString("agent.JobFlash")
var maxrunningjobs int = g.Cfg().GetInt("agent.MaxRunningJobs")
var runningJobs map[int]*gproc.Process = make(map[int]*gproc.Process)
var envPrefix string = g.Cfg().GetString("agent.EnvPrefix")
var agentInclude string = g.Cfg().GetString("agent.Include")
var wsAgentSend chan model.WsAgentSend = make(chan model.WsAgentSend)

type AgentsMap struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type AgentsList []AgentsMap

var agents AgentsList = make(AgentsList, 0)

func main() {
	AgentCICD.AgentRun()
}

func (s *agentCICD) GetAgentsList(isreload bool) AgentsList {
	if len(agents) != 0 && !isreload {
		return agents
	}

	var newagents AgentsList = make(AgentsList, 0)
	var agentsList AgentsList
	var agentsStr string = g.Cfg().GetString("agents")

	if agentsStr != "" {
		if err := gjson.DecodeTo(agentsStr, &agentsList); err != nil {
			glog.Errorf("%s decode failed. %s", g.Cfg().GetFileName(), err)
		}
		newagents = append(newagents, agentsList...)
	}

	if agentInclude != "" {
		files := s.HanleIncludeConfig(agentInclude)
		for _, f := range files {
			if g.Cfg().Available(f) {
				newconfig := g.Cfg().SetFileName(f)
				agentsStr = newconfig.GetString("agents")
				if agentsStr == "" {
					continue
				}
				if err := gjson.DecodeTo(agentsStr, &agentsList); err != nil {
					glog.Errorf("%s decode failed. %s", f, err)
					continue
				}
				newagents = append(newagents, agentsList...)
			}
		}
	}

	if len(newagents) > 0 {
		agents = newagents
		jobFlashStatus, err := json.Marshal(agents)
		if err != nil {
			glog.Error(err)
		}
		glog.Info(jobFlashStatus)
		jobFlashPath := dataPathDir + jobFlash
		s.WriteFile(jobFlashPath, string(jobFlashStatus))
	}
	return agents
}

func (s *agentCICD) HanleIncludeConfig(pattern string) []string {
	var filenames []string
	files, err := filepath.Glob(pattern)
	if err != nil {
		glog.Error(err)
		panic(err)
	}
	filenames = append(filenames, files...)
	return filenames
}

func (s *agentCICD) SendJson() model.WsAgentSend {

	var agentsList AgentsList

	select {
	case msg := <-wsAgentSend:
		return msg
	default:
		var agentSent = model.WsAgentSend{}
		var agentSentMap = model.WsAgentSendMap{}
		agentsList = s.GetAgentsList(false)
		for _, agent := range agentsList {
			agentSentMap.AgentId = agent.ID
			agentSentMap.AgentName = agent.Name
			agentSent = append(agentSent, agentSentMap)
		}
		return agentSent
	}

}

func (s *agentCICD) GetExecutable(scriptbody string) string {
	if len(scriptbody) < 3 {
		glog.Error("scriptbody is empty")
		return ""
	}
	if scriptbody[0:2] == "#!" {
		return scriptbody[2:strings.Index(scriptbody, "\n")]
	} else {
		return "/usr/bin/env bash"
	}
}

func (s *agentCICD) WriteFile(path string, content string) error {
	glog.Debug("Write file: ", path)
	if err := gfile.PutContents(path, content); err != nil {
		glog.Error(err)
		return err
	}
	return nil
}

func FileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func (s *agentCICD) ReadFile(path string) string {
	glog.Debug("Read file: ", path)
	if !FileExists(path) {
		glog.Debug("file not exist: ", path)
		return ""
	}
	// if _, err := os.Stat(path); os.IsNotExist(err) {
	// 	return ""
	// }
	content := gfile.GetContents(path)
	return content
}

func (s *agentCICD) SetStatus(jobId int, jobStatus string) error {
	jobPathscriptJson := dataPathDir + strconv.Itoa(jobId) + ".json"
	oldJobStatus := s.GetStatus(jobId)
	if jobStatus == oldJobStatus {
		return nil
	}
	var jobMeta = model.JobMeta{}
	jobMeta.ID = jobId
	jobMeta.JobStatus = jobStatus
	jobJson, _ := json.Marshal(&jobMeta)
	fileLock := flock.New(jobPathscriptJson)
	err := fileLock.Lock()
	if err != nil {
		glog.Error(err)
	}
	if err := s.WriteFile(jobPathscriptJson, string(jobJson)); err != nil {
		glog.Error(err)
		return err
	}
	fileLock.Unlock()
	return nil
}

func (s *agentCICD) GetStatus(jobId int) string {
	jobPathscriptJson := dataPathDir + strconv.Itoa(jobId) + ".json"
	fileLock := flock.New(jobPathscriptJson)
	err := fileLock.RLock()
	if err != nil {
		glog.Error(err)
	}
	jobJson := s.ReadFile(jobPathscriptJson)
	fileLock.Unlock()
	if jobJson == "" {
		glog.Debugf("fileName %s with content empty!", jobPathscriptJson)
		return ""
	}
	var jobMeta = model.JobMeta{}
	if err := json.Unmarshal([]byte(jobJson), &jobMeta); err != nil {
		glog.Error(err)
		glog.Debugf("fileName %s with content: %s !", jobPathscriptJson, jobJson)
		return ""
	}
	return jobMeta.JobStatus
}

func (s *agentCICD) KillJob(jobId int) {
	glog.Error("Begin in KillJob")
	if runningProcess, ok := runningJobs[jobId]; ok {
		if runningProcess.ProcessState.Exited() {
			exitCode := runningProcess.ProcessState.ExitCode()
			delete(runningJobs, jobId)
			glog.Debugf("Exit newjob: %d exitcode: %d", jobId, exitCode)
		} else {
			if err := runningProcess.Kill(); err != nil {
				glog.Error(err)
			}
			if runningProcess.ProcessState.Exited() {
				exitCode := runningProcess.ProcessState.ExitCode()
				delete(runningJobs, jobId)
				glog.Debugf("Kill job: %d exitcode: %d", jobId, exitCode)
			}
		}
	}
}

func (s *agentCICD) RunCommand(jobId int, runCommand string, scriptEnvs []string) {
	// defer delete(runningJobs, jobId)
	glog.Debugf("recvScriptEnvs: %+v", scriptEnvs)
	glog.Debugf("recvScriptEnvs: %#v", scriptEnvs)
	newprocess := gproc.NewProcessCmd(runCommand, scriptEnvs)
	newpid, err := newprocess.Start()
	if err != nil {
		glog.Error(newpid, err)
	}
	glog.Debugf("Run newjob: %d pid: %d", jobId, newpid)
	if err := s.SetStatus(jobId, "running"); err != nil {
		glog.Error(newpid, err)
	}
	runningJobs[jobId] = newprocess
	if err = newprocess.Wait(); err != nil {
		glog.Errorf("Command finished with error: %v", err)
	}
	glog.Debugf("Finished Run newjob: %d pid: %d", jobId, newpid)
	// glog.Debug(newprocess.ProcessState.Exited())
	// glog.Debug(newprocess.ProcessState.ExitCode())

	newprocess.ProcessState.Exited()
	newprocess.ProcessState.ExitCode()
	if newprocess.ProcessState.Exited() {
		exitCode := newprocess.ProcessState.ExitCode()
		glog.Debugf("Exit newjob: %d pid: %d exitcode: %d", jobId, newpid, exitCode)
		if exitCode == 0 {
			if err := s.SetStatus(jobId, "success"); err != nil {
				glog.Error(newpid, err)
			}
		} else {
			if err := s.SetStatus(jobId, "failed"); err != nil {
				glog.Error(newpid, err)
			}
		}
		delete(runningJobs, jobId)
	}

}

func (s *agentCICD) HandleJob(jobv *model.WsServerSendMap) *model.WsAgentSendMap {
	var sendMap = &model.WsAgentSendMap{}
	jobId := jobv.JobID
	jobStatus := jobv.JobStatus
	sendMap.AgentId = jobv.AgentId
	sendMap.AgentName = jobv.AgentName
	sendMap.JobId = jobId
	if jobStatus == "success" || jobStatus == "failed" {
		sendMap.JobStatus = jobStatus
		jobPath := dataPathDir + strconv.Itoa(jobId)
		jobPathOutput := jobPath + ".output"
		output := s.ReadFile(jobPathOutput)
		sendMap.JobOutput = output
		return sendMap
	}
	if jobStatus == "running" {
		localJobStatus := s.GetStatus(jobId)
		if localJobStatus == "running" {
			if _, ok := runningJobs[jobId]; !ok {
				sendMap.JobStatus = "pending"
				return sendMap
			}
		}
		sendMap.JobStatus = localJobStatus
		jobPath := dataPathDir + strconv.Itoa(jobId)
		jobPathOutput := jobPath + ".output"
		output := s.ReadFile(jobPathOutput)
		sendMap.JobOutput = output
		return sendMap
	}
	if jobStatus == "aborted" {
		s.KillJob(jobId)
		sendMap.JobStatus = s.GetStatus(jobId)
		return sendMap
	}
	oldJobStatus := s.GetStatus(jobId)
	if oldJobStatus == "" {
		if err := s.SetStatus(jobId, "pending"); err != nil {
			glog.Error(jobId, err)
		}
	}
	jobPath := dataPathDir + strconv.Itoa(jobId)
	jobPathOutput := jobPath + ".output"
	if _, ok := runningJobs[jobId]; !ok {
		scriptBody := jobv.Body + "\n"
		scriptBody = strings.Replace(scriptBody, "\r\n", "\n", -1)
		jobPathscriptBody := jobPath + ".scriptbody"
		s.WriteFile(jobPathscriptBody, scriptBody)
		scriptArgs := jobv.Args + "\n"
		scriptArgs = strings.Replace(scriptArgs, "\r\n", "\n", -1)
		jobPathscriptArgs := jobPath + ".scriptargs"
		s.WriteFile(jobPathscriptArgs, scriptArgs)
		var scriptEnvs []string
		envAgentName := strings.Split(jobv.AgentName, ":")[0]
		scriptEnvs = append(scriptEnvs, envPrefix+"AGENTNAME"+"="+envAgentName)
		for k, v := range jobv.Envs {
			scriptEnvs = append(scriptEnvs, envPrefix+k+"="+v)
		}
		execommand := s.GetExecutable(scriptBody)
		runcommand := execommand + " " + jobPathscriptBody + " " + jobPathscriptArgs + " >>" + jobPathOutput + " 2>&1"
		glog.Debugf("Run jobId: %d with Command: %s and scriptEnvs: %s", jobId, runcommand, scriptEnvs)
		go s.RunCommand(jobId, runcommand, scriptEnvs)
	}
	sendMap.JobStatus = s.GetStatus(jobId)
	output := s.ReadFile(jobPathOutput)
	sendMap.JobOutput = output
	return sendMap
}

func (s *agentCICD) HandleRecvJson(recvJson *model.WsServerSend) {
	var sendJson model.WsAgentSend
	recvData := *recvJson
	for _, jobv := range recvData {
		if len(runningJobs) > maxrunningjobs {
			continue
		}
		if jobv.ErrMsg != "" {
			glog.Errorf("jobId: %d errmsg: %s", jobv.JobID, jobv.ErrMsg)
			continue
		}
		if jobv.JobID == 0 || jobv.JobStatus == "" {
			continue
		}
		glog.Debugf("recvjson: %+v", jobv)
		glog.Debugf("recvjson: %#v", jobv)
		var sendMap = s.HandleJob(&jobv)
		sendJson = append(sendJson, *sendMap)
	}

	if len(sendJson) > 0 {
		wsAgentSend <- sendJson
	}
}

func (s *agentCICD) AgentRun() {
	if err := gfile.Mkdir(dataPathDir); err != nil {
		glog.Error(err)
		panic(err)
		// os.Exit(1)
	}

	interrupt := make(chan os.Signal, 1)
	// signal.Notify(interrupt, os.Interrupt, syscall.SIGUSR1)
	signal.Notify(interrupt, os.Interrupt)
	// signal.Notify(interrupt, syscall.SIGTERM)
	reload := make(chan os.Signal, 1)
	signal.Notify(reload, syscall.SIGUSR1)

	var recvJson = new(model.WsServerSend)

	client := ghttp.NewWebSocketClient()
	client.HandshakeTimeout = time.Second    // 设置超时时间
	client.Proxy = http.ProxyFromEnvironment // 设置代理

	for i := 0; i < 10; i++ {
		// time.Sleep(time.Second)
		select {
		case <-interrupt:
			glog.Info("interrupt2")
			os.Exit(1)
		case <-time.After(time.Second):
		}

		// conn, _, err := client.Dial("ws://127.0.0.1:8070/wsv1/wsci", nil)
		conn, _, err := client.Dial(wsUrl, nil)
		if err != nil {
			// panic(err)
			glog.Error("dial:", err)
			continue
		}
		defer conn.Close()

		done := make(chan struct{})

		go func() {
			defer close(done)
			for {
				err := conn.ReadJSON(&recvJson)
				if err != nil {
					glog.Error("read:", err)
					return
				}
				// glog.Infof("recv+v: %+v", recvJson)

				newjobs, _ := json.Marshal(recvJson)
				glog.Infof("recvjson: %s", string(newjobs))
				s.HandleRecvJson(recvJson)

			}
		}()

		ticker := time.NewTicker(time.Duration(1000000000 * syncInterval))
		defer ticker.Stop()

	L:
		for {
			// T:
			select {
			case <-done:
				break L
			case <-ticker.C:
				// glog.Info("*********************************")
				sendJson := s.SendJson()
				err := conn.WriteJSON(sendJson)
				if err != nil {
					glog.Error("write:", err)
					return
				}
				// glog.Infof("send+v: %+v", sendJson)
				// glog.Infof("send#v: %#v", sendJson)
				newjobs, _ := json.Marshal(sendJson)
				glog.Infof("sendjson: %s", string(newjobs))
				// glog.Info("###################################")
			case <-interrupt:
				glog.Info("interrupt1")
				err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				if err != nil {
					glog.Warningf("write close:", err)
					return
				}
				select {
				case <-done:
				case <-time.After(time.Second):
				}
				return
			case <-reload:
				glog.Info("reload")
				s.GetAgentsList(true)

			}
		}
	}
}
