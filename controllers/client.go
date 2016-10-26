package controllers

import (
	"encoding/json"
	"fmt"
	"moduleab_server/common"
	"moduleab_server/models"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/astaxie/beego"
	"github.com/gorilla/websocket"
)

const (
	ClientWebSocketReplyGot  = "GOT"
	ClientWebSocketReplyDone = "DONE"
	ClientWebSocketReplyBye  = "BYE"
)

const (
	ClientRunStatusAll = iota
	ClientRunStatusRunning
	ClientRunStatusStopped
)

type ClientStatusMsg struct {
	HostId string
	Status int
}

var (
	ClientStatus     map[string]int
	ChanClientStatus chan ClientStatusMsg
)

func init() {
	ClientStatus = make(map[string]int)
	ChanClientStatus = make(chan ClientStatusMsg, 2<<10)
	go clientStatus()
	AddPrivilege("GET", "^/api/v1/client/signal/(.+)/ws$", models.RoleFlagNone)
}

func clientStatus() {
	for {
		select {
		case s := <-ChanClientStatus:
			ClientStatus[s.HostId] = s.Status
		}
	}
}

type ClientController struct {
	beego.Controller
}

func (c *ClientController) Prepare() {
	if c.Ctx.Input.Header("Signature") != "" {
		err := common.AuthWithKey(c.Ctx)
		if err != nil {
			c.Data["json"] = map[string]string{
				"error": err.Error(),
			}
			c.Ctx.Output.SetStatus(http.StatusForbidden)
			c.ServeJSON()
		}
	} else {
		id := c.GetSession("id")
		if id == nil {
			c.Data["json"] = map[string]string{
				"error": "You need login first.",
			}
			c.Ctx.Output.SetStatus(http.StatusUnauthorized)
			c.ServeJSON()
		} else {
			if !CheckPrivileges(id.(string), c.Ctx) {
				c.Data["json"] = map[string]string{
					"error": "No privileges.",
				}
				c.Ctx.Output.SetStatus(http.StatusForbidden)
				c.ServeJSON()
			}
		}
	}
}

// @Title getClientConf
// @Description getClientConf
// @Param	body		body 	models.Hosts	true		"body for host content"
// @Success 200
// @Failure 403 body is empty
// @router /config [get]
func (c *ClientController) GetAliConfig() {
	c.Data["json"] = map[string]string{
		"ali_key":    beego.AppConfig.String("aliapi::apikey"),
		"ali_secret": beego.AppConfig.String("aliapi::secret"),
	}
	c.ServeJSON()
}

// @Title getSignalsWs
// @router /signal/:name/ws [get]
func (c *ClientController) WebSocket() {
	name := c.GetString(":name")
	beego.Debug("[C] Got name:", name)
	if name != "" {
		host := &models.Hosts{
			Name: name,
		}
		hosts, err := models.GetHosts(host, 1, 0)
		if err != nil {
			c.Data["json"] = map[string]string{
				"message": fmt.Sprint("Failed to get with name:", name),
				"error":   err.Error(),
			}
			beego.Warn("[C] Got error:", err)
			c.Ctx.Output.SetStatus(http.StatusInternalServerError)
			c.ServeJSON()
			return
		}
		if len(hosts) == 0 {
			beego.Debug("[C] Got nothing with name:", name)
			c.Ctx.Output.SetStatus(http.StatusNotFound)
			c.ServeJSON()
			return
		}
		HostId := hosts[0].Id

		ws, err := websocket.Upgrade(
			c.Ctx.ResponseWriter, c.Ctx.Request,
			nil, 1024, 1024)
		if err != nil {
			c.Data["json"] = map[string]string{
				"message": "Failed on upgrading to websocket",
				"error":   err.Error(),
			}
			beego.Warn("[C] Got error:", err)
			c.Ctx.Output.SetStatus(http.StatusInternalServerError)
			c.ServeJSON()
			return
		}
		defer ws.Close()

		tick := beego.AppConfig.DefaultInt64("websocket::pingperiod", 5)
		ticker := time.NewTicker(
			time.Duration(tick) * time.Second)
		defer ticker.Stop()

		timeout := beego.AppConfig.DefaultInt64("websocket::timeout", 10)
		ws.SetReadDeadline(time.Now().Add(
			time.Duration(timeout) * time.Second),
		)
		ws.SetWriteDeadline(time.Now().Add(
			time.Duration(timeout) * time.Second),
		)

		var runningStatus = ClientStatusMsg{
			HostId: HostId,
			Status: ClientRunStatusStopped,
		}
		ws.SetPongHandler(func(string) error {
			beego.Debug("Host:", name, "is still alive.")
			ws.SetReadDeadline(time.Now().Add(
				time.Duration(timeout) * time.Second),
			)
			ws.SetWriteDeadline(time.Now().Add(
				time.Duration(timeout) * time.Second),
			)

			runningStatus.Status = ClientRunStatusRunning
			ChanClientStatus <- runningStatus

			return nil
		})

		defer func() {
			runningStatus.Status = ClientRunStatusStopped
			ChanClientStatus <- runningStatus
		}()

		var c chan models.Signal
		c, ok := models.SignalChannels[HostId]
		if !ok {
			c = make(chan models.Signal, 1024)
			models.SignalChannels[HostId] = c
		}

		// Start read routine
		go func() {
			defer ws.Close()
			for {
				_, bConfirm, err := ws.ReadMessage()
				if websocket.IsCloseError(err,
					websocket.CloseGoingAway) {
					beego.Info("Host", name, "is offline.")
					return
				} else if err != nil {
					beego.Warn("Error on reading:", err.Error())
					return
				}
				s := strings.Split(string(bConfirm), " ")
				if s[0] == ClientWebSocketReplyDone {
					models.DeleteSignal(HostId, s[1])
				}
			}
		}()

		for {
			select {
			case s := <-models.SignalChannels[HostId]:
				ws.WriteJSON(s)
			case <-ticker.C:
				beego.Debug("Websocket ping:", name)
				err := ws.WriteMessage(websocket.PingMessage, []byte{})
				if err != nil {
					beego.Warn("Got error on ping", err.Error())
					return
				}
			}
		}
	}
}

// @Title getSignals
// @router /signal/:name [get]
func (c *ClientController) GetSignals() {
	name := c.GetString(":name")
	defer c.ServeJSON()
	beego.Debug("[C] Got name:", name)
	if name != "" {
		host := &models.Hosts{
			Name: name,
		}
		hosts, err := models.GetHosts(host, 0, 0)
		if err != nil {
			c.Data["json"] = map[string]string{
				"message": fmt.Sprint("Failed to get with name:", name),
				"error":   err.Error(),
			}
			beego.Warn("[C] Got error:", err)
			c.Ctx.Output.SetStatus(http.StatusInternalServerError)
			return
		}
		if len(hosts) == 0 {
			beego.Debug("[C] Got nothing with name:", name)
			c.Ctx.Output.SetStatus(http.StatusNotFound)
			return
		}
		c.Data["json"] = models.GetSignals(hosts[0].Id)
		c.Ctx.Output.SetStatus(http.StatusOK)
	}
}

// @Title getSignal
// @router /signal/:name/:id [get]
func (c *ClientController) GetSignal() {
	name := c.GetString(":name")
	id := c.GetString(":id")
	defer c.ServeJSON()
	beego.Debug("[C] Got name:", name)
	if name != "" {
		host := &models.Hosts{
			Name: name,
		}
		hosts, err := models.GetHosts(host, 0, 0)
		if err != nil {
			c.Data["json"] = map[string]string{
				"message": fmt.Sprint("Failed to get with name:", name),
				"error":   err.Error(),
			}
			beego.Warn("[C] Got error:", err)
			c.Ctx.Output.SetStatus(http.StatusInternalServerError)
			return
		}
		if len(hosts) == 0 {
			beego.Debug("[C] Got nothing with name:", name)
			c.Ctx.Output.SetStatus(http.StatusNotFound)
			return
		}
		signal, err := models.GetSignal(hosts[0].Id, id)
		if err != nil {
			c.Data["json"] = map[string]string{
				"message": fmt.Sprint("Got nothing with id:", id),
			}
			c.Ctx.Output.SetStatus(http.StatusNotFound)
			return
		}

		c.Data["json"] = signal
		c.Ctx.Output.SetStatus(http.StatusOK)
	}
}

// @Title createSignal
// @router /signal/:name [post]
func (c *ClientController) PostSignal() {
	name := c.GetString(":name")
	defer c.ServeJSON()
	beego.Debug("[C] Got name:", name)
	if name != "" {
		host := &models.Hosts{
			Name: name,
		}
		hosts, err := models.GetHosts(host, 0, 0)
		if err != nil {
			c.Data["json"] = map[string]string{
				"message": fmt.Sprint("Failed to get with name:", name),
				"error":   err.Error(),
			}
			beego.Warn("[C] Got error:", err)
			c.Ctx.Output.SetStatus(http.StatusInternalServerError)
			return
		}
		if len(hosts) == 0 {
			beego.Debug("[C] Got nothing with name:", name)
			c.Ctx.Output.SetStatus(http.StatusNotFound)
			return
		}
		var signal models.Signal
		err = json.Unmarshal(c.Ctx.Input.RequestBody, &signal)
		if err != nil {
			beego.Warn("[C] Got error:", err)
			c.Data["json"] = map[string]string{
				"message": "Bad request",
				"error":   err.Error(),
			}
			c.Ctx.Output.SetStatus(http.StatusBadRequest)
			return
		}
		beego.Debug("[C] Got data:", signal)
		id, err := models.AddSignal(hosts[0].Id, signal)
		if err != nil {
			beego.Warn("[C] Got error:", err)
			c.Data["json"] = map[string]string{
				"message": "Failed to add new signal",
				"error":   err.Error(),
			}
			c.Ctx.Output.SetStatus(http.StatusBadRequest)
			return
		}
		c.Data["json"] = map[string]string{
			"id": id,
		}
		c.Ctx.Output.SetStatus(http.StatusCreated)
	}
}

// @Title deleteSignal
// @router /signal/:name/:id [delete]
func (c *ClientController) DeleteSignal() {
	name := c.GetString(":name")
	id := c.GetString(":id")
	defer c.ServeJSON()
	beego.Debug("[C] Got name:", name)
	if name != "" {
		host := &models.Hosts{
			Name: name,
		}
		hosts, err := models.GetHosts(host, 0, 0)
		if err != nil {
			c.Data["json"] = map[string]string{
				"message": fmt.Sprint("Failed to get with name:", name),
				"error":   err.Error(),
			}
			beego.Warn("[C] Got error:", err)
			c.Ctx.Output.SetStatus(http.StatusInternalServerError)
			return
		}
		if len(hosts) == 0 {
			beego.Debug("[C] Got nothing with name:", name)
			c.Ctx.Output.SetStatus(http.StatusNotFound)
			return
		}
		err = models.DeleteSignal(hosts[0].Id, id)
		if err != nil {
			beego.Warn("[C] Got error:", err)
			c.Data["json"] = map[string]string{
				"message": "Delete failed",
				"error":   err.Error(),
			}
			if err == models.ErrorSignalNotFound {
				c.Ctx.Output.SetStatus(http.StatusNotFound)
			} else {
				c.Ctx.Output.SetStatus(http.StatusInternalServerError)
			}
			return
		}

		c.Data["json"] = map[string]string{
			"message": fmt.Sprint("Got nothing with id:", id),
		}
		c.Ctx.Output.SetStatus(http.StatusNotFound)
	}
}

// @Title notifySignal
// @router /signal/:name/:id [post]
func (c *ClientController) NotifySignal() {
	name := c.GetString(":name")
	id := c.GetString(":id")
	defer c.ServeJSON()
	if name != "" {
		host := &models.Hosts{
			Name: name,
		}
		hosts, err := models.GetHosts(host, 0, 0)
		if err != nil {
			c.Data["json"] = map[string]string{
				"message": fmt.Sprint("Failed to get with name:", name),
				"error":   err.Error(),
			}
			beego.Warn("[C] Got error:", err)
			c.Ctx.Output.SetStatus(http.StatusInternalServerError)
			return
		}
		if len(hosts) == 0 {
			beego.Debug("[C] Got nothing with name:", name)
			c.Ctx.Output.SetStatus(http.StatusNotFound)
			return
		}
		if id == "" {
			c.Ctx.Output.SetStatus(http.StatusBadRequest)
			return
		}

		err = models.NotifySignal(hosts[0].Id, id)
		if err != nil {
			c.Data["json"] = map[string]string{
				"message": fmt.Sprint("Failed to notify to:", name),
				"error":   err.Error(),
			}
			beego.Warn("[C] Got error:", err)
			c.Ctx.Output.SetStatus(http.StatusInternalServerError)
			return
		}
		c.Ctx.Output.SetStatus(http.StatusOK)
		return
	}
}

// @Title getClientStatus
// @router /config/status [get]
func (c *ClientController) GetStatus() {
	defer c.ServeJSON()
	var lock = new(sync.Mutex)
	lock.Lock()
	defer lock.Unlock()

	c.Data["json"] = ClientStatus
	c.Ctx.Output.SetStatus(http.StatusOK)
}
