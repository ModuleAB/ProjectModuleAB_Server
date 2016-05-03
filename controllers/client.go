package controllers

import (
	"encoding/json"
	"fmt"
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

type ClientController struct {
	beego.Controller
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
			this.Ctx.ResponseWriter, this.Ctx.Request,
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
		var c chan models.Signal
		c, ok := models.SignalChannels[HostsId]
		if !ok {
			c = make(chan models.Signal, 1024)
			models.SignalChannels[HostsId] = c
		}

		for {
			select {
			case s := <-models.SignalChannels[HostsId]:
				ws.WriteJSON(s)
				_, bConfirm, _ := ws.ReadMessage()
				if string(bConfirm) == ClientWebSocketReplyDone {
					models.DeleteSignal(HostId, s["id"])
				}
			case <-time.After(5 * time.Second):
				ws.WriteControl(websocket.PingMessage,
					[]byte("Are you alive?"),
					time.Now().Add(10*time.Second))
				mType, _, err := ws.ReadMessage()
				if err != nil || mType != websocket.PongMessage {
					beego.Info("Host", name, "become offline.")
					return
				}
			default:
				mType, _, _ := ws.ReadMessage()
				if mType == websocket.CloseMessage {
					beego.Info("Host", name, "become offline.")
					return
				}
			}
		}
		defer ws.Close()
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
