package controllers

import (
	"encoding/json"
	"fmt"
	"moduleab_server/common"
	"moduleab_server/models"
	"net/http"
	"time"

	"github.com/astaxie/beego"
	"github.com/gorilla/websocket"
)

const (
	ClientWebSocketReplyGot  = "GOT"
	ClientWebSocketReplyDone = "DONE"
	ClientWebSocketReplyBye  = "BYE"
)

func init() {
	AddPrivilege("GET", "^/api/v1/client/signal/(.+)/ws$", models.RoleFlagNone)
}

type ClientController struct {
	beego.Controller
}

func (h *ClientController) Prepare() {
	if h.Ctx.Input.Header("Signature") != "" {
		err := common.AuthWithKey(h.Ctx)
		if err != nil {
			h.Data["json"] = map[string]string{
				"error": err.Error(),
			}
			h.Ctx.Output.SetStatus(http.StatusForbidden)
			h.ServeJSON()
		}
	} else {
		id := h.GetSession("id")
		if id == nil {
			h.Data["json"] = map[string]string{
				"error": "You need login first.",
			}
			h.Ctx.Output.SetStatus(http.StatusUnauthorized)
			h.ServeJSON()
		} else {
			if !CheckPrivileges(id.(string), h.Ctx) {
				h.Data["json"] = map[string]string{
					"error": "No privileges.",
				}
				h.Ctx.Output.SetStatus(http.StatusForbidden)
				h.ServeJSON()
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

		tick, err := beego.AppConfig.Int64("websocket::pingperiod")
		if err != nil {
			tick = 5
		}
		ticker := time.NewTicker(
			time.Duration(tick) * time.Second)
		defer ticker.Stop()

		timeout, err := beego.AppConfig.Int64("websocket::timeout")
		if err != nil {
			timeout = 10
		}
		ws.SetReadDeadline(time.Now().Add(
			time.Duration(timeout) * time.Second),
		)
		ws.SetPongHandler(func(string) error {
			ws.SetReadDeadline(time.Now().Add(
				time.Duration(timeout) * time.Second),
			)
			return nil
		})
		var c chan models.Signal
		c, ok := models.SignalChannels[HostId]
		if !ok {
			c = make(chan models.Signal, 1024)
			models.SignalChannels[HostId] = c
		}

		for {
			select {
			case s := <-models.SignalChannels[HostId]:
				ws.WriteJSON(s)
				_, bConfirm, err := ws.ReadMessage()
				if websocket.IsCloseError(err,
					websocket.CloseGoingAway) {
					beego.Info("Host", name, "is offline.")
					break
				} else if err != nil {
					beego.Warn("Error on reading:", err.Error())
					break
				}
				if string(bConfirm) == ClientWebSocketReplyDone {
					models.DeleteSignal(HostId, s["id"].(string))
				}
			case <-ticker.C:
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
			c.ServeJSON()
			return
		}
		if len(hosts) == 0 {
			beego.Debug("[C] Got nothing with name:", name)
			c.Ctx.Output.SetStatus(http.StatusNotFound)
			c.ServeJSON()
			return
		}
		c.Data["json"] = models.GetSignals(hosts[0].Id)
		c.Ctx.Output.SetStatus(http.StatusOK)
		c.ServeJSON()
		return
	}
}

// @Title getSignal
// @router /signal/:name/:id [get]
func (c *ClientController) GetSignal() {
	name := c.GetString(":name")
	id := c.GetString(":id")
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
			c.ServeJSON()
			return
		}
		if len(hosts) == 0 {
			beego.Debug("[C] Got nothing with name:", name)
			c.Ctx.Output.SetStatus(http.StatusNotFound)
			c.ServeJSON()
			return
		}
		signal, err := models.GetSignal(hosts[0].Id, id)
		if err != nil {
			c.Data["json"] = map[string]string{
				"message": fmt.Sprintf("Got nothing with id:", id),
			}
			c.Ctx.Output.SetStatus(http.StatusNotFound)
			c.ServeJSON()
			return
		}

		c.Data["json"] = signal
		c.Ctx.Output.SetStatus(http.StatusOK)
		c.ServeJSON()
		return
	}
}

// @Title createSignal
// @router /signal/:name [post]
func (c *ClientController) PostSignal() {
	name := c.GetString(":name")
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
			c.ServeJSON()
			return
		}
		if len(hosts) == 0 {
			beego.Debug("[C] Got nothing with name:", name)
			c.Ctx.Output.SetStatus(http.StatusNotFound)
			c.ServeJSON()
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
			c.ServeJSON()
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
			c.ServeJSON()
			return
		}
		c.Data["json"] = map[string]string{
			"id": id,
		}
		c.Ctx.Output.SetStatus(http.StatusCreated)
		c.ServeJSON()
		return
	}
}

// @Title deleteSignal
// @router /signal/:name/:id [delete]
func (c *ClientController) DeleteSignal() {
	name := c.GetString(":name")
	id := c.GetString(":id")
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
			c.ServeJSON()
			return
		}
		if len(hosts) == 0 {
			beego.Debug("[C] Got nothing with name:", name)
			c.Ctx.Output.SetStatus(http.StatusNotFound)
			c.ServeJSON()
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
			c.ServeJSON()
			return
		}

		c.Data["json"] = map[string]string{
			"message": fmt.Sprintf("Got nothing with id:", id),
		}
		c.Ctx.Output.SetStatus(http.StatusNotFound)
		c.ServeJSON()
		return
	}
}

// @Title notifySignal
// @router /signal/:name/:id [post]
func (c *ClientController) NotifySignal() {
	name := c.GetString(":name")
	id := c.GetString(":id")
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
			c.ServeJSON()
			return
		}
		if len(hosts) == 0 {
			beego.Debug("[C] Got nothing with name:", name)
			c.Ctx.Output.SetStatus(http.StatusNotFound)
			c.ServeJSON()
			return
		}
		if id == "" {
			c.Ctx.Output.SetStatus(http.StatusBadRequest)
			c.ServeJSON()
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
			c.ServeJSON()
			return
		}
		c.Ctx.Output.SetStatus(http.StatusOK)
		c.ServeJSON()
		return
	}

}
