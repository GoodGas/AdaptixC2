package connector

import (
	"AdaptixServer/core/utils/logs"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	adaptix "github.com/Adaptix-Framework/axc2"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AgentConfig struct {
	ListenerName string `json:"listener_name"`
	ListenerType string `json:"listener_type"`
	AgentName    string `json:"agent"`
	Config       string `json:"config"`
}

func (tc *TsConnector) TcAgentGenerate(ctx *gin.Context) {
	var (
		agentConfig     AgentConfig
		err             error
		listenerProfile []byte
		listenerWM      string
		fileContent     []byte
		fileName        string
	)

	err = ctx.ShouldBindJSON(&agentConfig)
	if err != nil {
		_ = ctx.Error(errors.New("invalid agent config"))
		return
	}

	listenerWM, listenerProfile, err = tc.teamserver.TsListenerGetProfile(agentConfig.ListenerName, agentConfig.ListenerType)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": err.Error(), "ok": false})
		return
	}

	fileContent, fileName, err = tc.teamserver.TsAgentGenerate(agentConfig.AgentName, agentConfig.Config, listenerWM, listenerProfile)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": err.Error(), "ok": false})
		return
	}

	encodedContent := base64.StdEncoding.EncodeToString([]byte(fileName)) + ":" + base64.StdEncoding.EncodeToString(fileContent)

	ctx.JSON(http.StatusOK, gin.H{"message": encodedContent, "ok": true})
}

type CommandData struct {
	AgentName string `json:"name"`
	AgentId   string `json:"id"`
	UI        bool   `json:"ui"`
	CmdLine   string `json:"cmdline"`
	Data      string `json:"data"`
	HookId    string `json:"ax_hook_id"`
}

func (tc *TsConnector) TcAgentCommandExecute(ctx *gin.Context) {
	var (
		username    string
		commandData CommandData
		args        map[string]any
		ok          bool
		err         error
	)

	err = ctx.ShouldBindJSON(&commandData)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": "invalid JSON data", "ok": false})
		return
	}

	value, exists := ctx.Get("username")
	if !exists {
		ctx.JSON(http.StatusOK, gin.H{"message": "Server error: username not found in context", "ok": false})
		return
	}

	username, ok = value.(string)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"message": "Server error: invalid username type in context", "ok": false})
		return
	}

	err = json.Unmarshal([]byte(commandData.Data), &args)
	if err != nil {
		logs.Debug("", "Error parsing commands JSON: %s\n", err.Error())
	}

	err = tc.teamserver.TsAgentCommand(commandData.AgentName, commandData.AgentId, username, commandData.HookId, commandData.CmdLine, commandData.UI, args)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": err.Error(), "ok": false})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "", "ok": true})
}

type CommandData2 struct {
	ObjectId string `json:"object_id"`
}

func (tc *TsConnector) TcAgentCommandFile(ctx *gin.Context) {
	var (
		username     string
		commandData  CommandData
		commandData2 CommandData2
		args         map[string]any
		ok           bool
		err          error
	)

	err = ctx.ShouldBindJSON(&commandData2)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": "invalid JSON data", "ok": false})
		return
	}

	value, exists := ctx.Get("username")
	if !exists {
		ctx.JSON(http.StatusOK, gin.H{"message": "Server error: username not found in context", "ok": false})
		return
	}

	username, ok = value.(string)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"message": "Server error: invalid username type in context", "ok": false})
		return
	}

	content, err := tc.teamserver.TsUploadGetFileContent(commandData2.ObjectId)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": err.Error(), "ok": false})
		return
	}

	err = json.Unmarshal(content, &commandData)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": err.Error(), "ok": false})
	}

	err = json.Unmarshal([]byte(commandData.Data), &args)
	if err != nil {
		logs.Debug("", "Error parsing commands JSON: %s\n", err.Error())
	}

	err = tc.teamserver.TsAgentCommand(commandData.AgentName, commandData.AgentId, username, commandData.HookId, commandData.CmdLine, commandData.UI, args)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": err.Error(), "ok": false})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "", "ok": true})
}

type AgentRemove struct {
	AgentIdArray []string `json:"agent_id_array"`
}

func (tc *TsConnector) TcAgentConsoleRemove(ctx *gin.Context) {
	var (
		agentRemove AgentRemove
		err         error
	)

	err = ctx.ShouldBindJSON(&agentRemove)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": "invalid JSON data", "ok": false})
		return
	}

	var errorsSlice []string
	for _, agentId := range agentRemove.AgentIdArray {
		err = tc.teamserver.TsAgentConsoleRemove(agentId)
		if err != nil {
			errorsSlice = append(errorsSlice, err.Error())
		}
	}

	if len(errorsSlice) > 0 {
		message := ""
		for i, errorMessage := range errorsSlice {
			message += fmt.Sprintf("%d. %s\n", i+1, errorMessage)
		}

		ctx.JSON(http.StatusOK, gin.H{"message": message, "ok": false})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "", "ok": true})
}

func (tc *TsConnector) TcAgentRemove(ctx *gin.Context) {
	var (
		agentRemove AgentRemove
		err         error
	)

	err = ctx.ShouldBindJSON(&agentRemove)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": "invalid JSON data", "ok": false})
		return
	}

	var errorsSlice []string
	for _, agentId := range agentRemove.AgentIdArray {
		err = tc.teamserver.TsAgentRemove(agentId)
		if err != nil {
			errorsSlice = append(errorsSlice, err.Error())
		}
	}

	if len(errorsSlice) > 0 {
		message := ""
		for i, errorMessage := range errorsSlice {
			message += fmt.Sprintf("%d. %s\n", i+1, errorMessage)
		}

		ctx.JSON(http.StatusOK, gin.H{"message": message, "ok": false})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "", "ok": true})
}

/// Setters

type AgentTag struct {
	AgentIdArray []string `json:"agent_id_array"`
	Tag          string   `json:"tag"`
}

func (tc *TsConnector) TcAgentSetTag(ctx *gin.Context) {
	var (
		agentTag AgentTag
		err      error
	)

	err = ctx.ShouldBindJSON(&agentTag)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": "invalid JSON data", "ok": false})
		return
	}

	var errorsSlice []string
	for _, agentId := range agentTag.AgentIdArray {
		err = tc.teamserver.TsAgentSetTag(agentId, agentTag.Tag)
		if err != nil {
			errorsSlice = append(errorsSlice, err.Error())
		}
	}

	if len(errorsSlice) > 0 {
		message := ""
		for i, errorMessage := range errorsSlice {
			message += fmt.Sprintf("%d. %s\n", i+1, errorMessage)
		}

		ctx.JSON(http.StatusOK, gin.H{"message": message, "ok": false})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "", "ok": true})
}

type AgentMark struct {
	AgentIdArray []string `json:"agent_id_array"`
	Mark         string   `json:"mark"`
}

func (tc *TsConnector) TcAgentSetMark(ctx *gin.Context) {
	var (
		agentMark AgentMark
		err       error
	)

	err = ctx.ShouldBindJSON(&agentMark)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": "invalid JSON data", "ok": false})
		return
	}

	var errorsSlice []string
	for _, agentId := range agentMark.AgentIdArray {
		err = tc.teamserver.TsAgentSetMark(agentId, agentMark.Mark)
		if err != nil {
			errorsSlice = append(errorsSlice, err.Error())
		}
	}

	if len(errorsSlice) > 0 {
		message := ""
		for i, errorMessage := range errorsSlice {
			message += fmt.Sprintf("%d. %s\n", i+1, errorMessage)
		}

		ctx.JSON(http.StatusOK, gin.H{"message": message, "ok": false})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "", "ok": true})
}

type AgentColor struct {
	AgentIdArray []string `json:"agent_id_array"`
	Background   string   `json:"bc"`
	Foreground   string   `json:"fc"`
	Reset        bool     `json:"reset"`
}

func (tc *TsConnector) TcAgentSetColor(ctx *gin.Context) {
	var (
		agentColor AgentColor
		err        error
	)

	err = ctx.ShouldBindJSON(&agentColor)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": "invalid JSON data", "ok": false})
		return
	}

	var errorsSlice []string
	for _, agentId := range agentColor.AgentIdArray {
		err = tc.teamserver.TsAgentSetColor(agentId, agentColor.Background, agentColor.Foreground, agentColor.Reset)
		if err != nil {
			errorsSlice = append(errorsSlice, err.Error())
		}
	}

	if len(errorsSlice) > 0 {
		message := ""
		for i, errorMessage := range errorsSlice {
			message += fmt.Sprintf("%d. %s\n", i+1, errorMessage)
		}

		ctx.JSON(http.StatusOK, gin.H{"message": message, "ok": false})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "", "ok": true})
}

type AgentImpersonate struct {
	AgentId     string `json:"agent_id"`
	Impersonate string `json:"impersonate"`
	Elevated    bool   `json:"elevated"`
}

func (tc *TsConnector) TcAgentSetImpersonate(ctx *gin.Context) {
	var (
		agentImpersonate AgentImpersonate
		err              error
	)

	err = ctx.ShouldBindJSON(&agentImpersonate)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": "invalid JSON data", "ok": false})
		return
	}

	_ = tc.teamserver.TsAgentSetImpersonate(agentImpersonate.AgentId, agentImpersonate.Impersonate, agentImpersonate.Elevated)

	ctx.JSON(http.StatusOK, gin.H{"message": "", "ok": true})
}

/// Tasks

type AgentTaskDelete struct {
	AgentId string   `json:"agent_id"`
	TasksId []string `json:"tasks_array"`
}

func (tc *TsConnector) TcAgentTaskCancel(ctx *gin.Context) {
	var (
		agentTasks AgentTaskDelete
		err        error
	)

	err = ctx.ShouldBindJSON(&agentTasks)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": "invalid JSON data", "ok": false})
		return
	}

	var errorsSlice []string
	for _, taskId := range agentTasks.TasksId {
		err = tc.teamserver.TsTaskCancel(agentTasks.AgentId, taskId)
		if err != nil {
			errorsSlice = append(errorsSlice, err.Error())
		}
	}

	if len(errorsSlice) > 0 {
		message := ""
		for i, errorMessage := range errorsSlice {
			message += fmt.Sprintf("%d. %s\n", i+1, errorMessage)
		}

		ctx.JSON(http.StatusOK, gin.H{"message": message, "ok": false})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "", "ok": true})
}

func (tc *TsConnector) TcAgentTaskDelete(ctx *gin.Context) {
	var (
		agentTasks AgentTaskDelete
		err        error
	)

	err = ctx.ShouldBindJSON(&agentTasks)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": "invalid JSON data", "ok": false})
		return
	}

	var errorsSlice []string
	for _, taskId := range agentTasks.TasksId {
		err = tc.teamserver.TsTaskDelete(agentTasks.AgentId, taskId)
		if err != nil {
			errorsSlice = append(errorsSlice, err.Error())
		}
	}

	if len(errorsSlice) > 0 {
		message := ""
		for i, errorMessage := range errorsSlice {
			message += fmt.Sprintf("%d. %s\n", i+1, errorMessage)
		}

		ctx.JSON(http.StatusOK, gin.H{"message": message, "ok": false})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "", "ok": true})
}

type AgentTaskHook struct {
	AgentId     string `json:"a_id"`
	TaskId      string `json:"a_task_id"`
	HookId      string `json:"a_hook_id"`
	JobIndex    int    `json:"a_job_index"`
	MessageType int    `json:"a_msg_type"`
	Message     string `json:"a_message"`
	Text        string `json:"a_text"`
	Completed   bool   `json:"a_completed"`
}

func (tc *TsConnector) TcAgentTaskHook(ctx *gin.Context) {
	var (
		username  string
		tasksHook AgentTaskHook
		err       error
		ok        bool
	)

	err = ctx.ShouldBindJSON(&tasksHook)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": "invalid JSON data", "ok": false})
		return
	}

	value, exists := ctx.Get("username")
	if !exists {
		ctx.JSON(http.StatusOK, gin.H{"message": "Server error: username not found in context", "ok": false})
		return
	}

	username, ok = value.(string)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"message": "Server error: invalid username type in context", "ok": false})
		return
	}

	hookData := adaptix.TaskData{
		AgentId:     tasksHook.AgentId,
		TaskId:      tasksHook.TaskId,
		HookId:      tasksHook.HookId,
		Client:      username,
		MessageType: tasksHook.MessageType,
		Message:     tasksHook.Message,
		ClearText:   tasksHook.Text,
		Completed:   tasksHook.Completed,
	}

	err = tc.teamserver.TsTaskPostHook(hookData, tasksHook.JobIndex)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": err.Error(), "ok": false})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "", "ok": true})
}
