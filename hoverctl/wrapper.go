package hoverctl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	_ "net"
	"net/http"

	l "github.com/op/go-logging"
)

var log = l.MustGetLogger("politoctrl")

/*
func (d *Dataplane) postObject(url string, requestObj interface{}, responseObj interface{}) (err error) {
	b, err := json.Marshal(requestObj)
	if err != nil {
		return
	}
	resp, err := d.client.Post(d.baseUrl+url, "application/json", bytes.NewReader(b))
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		var body []byte
		if body, err = ioutil.ReadAll(resp.Body); err != nil {
			politoctrl.Error.Print(string(body))
		}
		return fmt.Errorf("module server returned %s", resp.Status)
	}
	if responseObj != nil {
		err = json.NewDecoder(resp.Body).Decode(responseObj)
	}
	return
}
*/

func (d *Dataplane) sendObject(method string, url string, requestObj interface{}, responseObj interface{}) (err error) {
	b, er := json.Marshal(requestObj)
	if er != nil {
		log.Warning("error during json marshal.")
		return er
	}

	var resp *http.Response
	var e error
	var req *http.Request

	switch method {
	case "", "POST":
		resp, e = d.client.Post(d.baseUrl+url, "application/json", bytes.NewReader(b))
	case "GET":
		resp, e = d.client.Get(d.baseUrl + url)
	default:
		req, e = http.NewRequest(method, d.baseUrl+url, bytes.NewReader(b))
		if err != nil {
			log.Errorf("%s\n", err)
		}
		resp, e = d.client.Do(req)
	}

	if e != nil {
		log.Warning(e)
		return e
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		var body []byte
		if body, err = ioutil.ReadAll(resp.Body); err != nil {
			log.Error(string(body))
		}
		return fmt.Errorf("module server returned %s", resp.Status)
	}
	if responseObj != nil {
		err = json.NewDecoder(resp.Body).Decode(responseObj)
	}
	return err
}

/*---------LINKS---------------*/

/*
	LinkModule(Dataplane,from,to)
  LinkModule(d,i:veth0,m:12345ab)
*/
func LinkPOST(d *Dataplane, from string, to string) (error, LinkEntry) {
	log.Infof("link POST %s <--> %s\n", from, to)

	request := map[string]interface{}{
		"from": from,
		"to":   to,
	}

	var link LinkEntry
	err := d.sendObject("POST", "/links/", request, &link)
	if err != nil {
		log.Warning(err)
		return err, link
	}

	log.Debugf("link POST %s <--> %s link id: %s", from, to, link.Id)
	return nil, link
}

func LinkGET(d *Dataplane, linkId string) (error, LinkEntry) {
	log.Infof("link GET %s\n", linkId)

	request := map[string]interface{}{}

	var link LinkEntry
	err := d.sendObject("GET", "/links/"+linkId, request, &link)
	if err != nil {
		log.Warning(err)
		return err, link
	}
	log.Debugf("link GET %s OK\n", linkId)
	return nil, link
}

func LinkDELETE(d *Dataplane, linkId string) (error, LinkEntry) {
	log.Infof("link DELETE %s\n", linkId)

	request := map[string]interface{}{}

	var link LinkEntry
	err := d.sendObject("DELETE", "/links/"+linkId, request, &link)
	if err != nil {
		log.Warning(err)
		return err, link
	}
	log.Debugf("link DELETE %s OK\n", linkId)
	return nil, link
}

/*------------MODULES-----------*/

/*
	ModulePOST(d,"bpf","myModulòeName",bpf.Modulename)
*/
func ModulePOST(d *Dataplane, moduleType string, displayName string, code string) (error, Module) {
	log.Infof("module POST %s\n", displayName)

	req := map[string]interface{}{
		"module_type":  moduleType,
		"display_name": displayName,
		"config": map[string]interface{}{
			"code": code,
		},
	}
	var module Module
	err := d.sendObject("POST", "/modules/", req, &module)
	if err != nil {
		log.Warning(err)
		return err, module
	}

	log.Debugf("module POST %s module id : %s\n", displayName, module.Id)

	return nil, module
}

func ModuleGET(d *Dataplane, moduleId string) (error, Module) {
	log.Infof("module GET %s \n", moduleId)

	req := map[string]interface{}{}
	var module Module
	err := d.sendObject("GET", "/modules/"+moduleId, req, &module)
	if err != nil {
		log.Warning(err)
		return err, module
	}

	log.Debugf("module GET %s OK\n", moduleId)

	return nil, module
}

func ModuleDELETE(d *Dataplane, moduleId string) (error, Module) {
	log.Infof("module DELETE %s\n", moduleId)

	req := map[string]interface{}{}
	var module Module
	err := d.sendObject("DELETE", "/modules/"+moduleId, req, &module)
	if err != nil {
		log.Warning(err)
		return err, module
	}

	log.Debugf("module DELETE %s OK\n", moduleId)

	return nil, module
}

/*it returns map[module-id]module provided by hover
eg. map["m:12345ab"] = module {}
*/
func ModuleListGET(d *Dataplane) (error, map[string]Module) {
	log.Info("getting modules list")
	modules := map[string]Module{}

	resp, e := d.client.Get(d.baseUrl + "/modules/")
	if e != nil {
		return e, modules
	}
	defer resp.Body.Close()

	var data []map[string]interface{}

	err := json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		log.Errorf("%s\n", err)
		return err, modules
	}

	for i := range data {
		m := Module{}
		item := data[i]

		id, _ := item["id"].(string)
		m.Id = id
		module_type, _ := item["module_type"].(string)
		m.ModuleType = module_type
		display_name, _ := item["display_name"].(string)
		m.DisplayName = display_name
		permissions, _ := item["permissions"].(string)
		m.Perm = permissions
		m.Config, _ = item["config"].(map[string]interface{})
		modules[id] = m
		log.Debugf("module-id:%15s   DisplayName: %s \n", id, display_name)
	}

	//log.Noticef("%+v\n", modules)
	log.Debug("getting modules list OK\n")
	return nil, modules
}

/*------------EXTERNAL-INTERFACES-----------*/

/*it returns map[iface-name]iface provided by hover
eg. map[veth1] = iface {Name:veth1, Id:42}
*/
func ExternalInterfacesListGET(d *Dataplane) (error, map[string]ExternalInterface) {
	log.Info("getting external_interfaces list")
	external_interfaces := map[string]ExternalInterface{}

	resp, e := d.client.Get(d.baseUrl + "/external_interfaces/")
	if e != nil {
		return e, external_interfaces
	}
	defer resp.Body.Close()

	var data []map[string]interface{}

	err := json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		log.Errorf("%s\n", err)
		return err, external_interfaces
	}

	for i := range data {
		item := data[i]
		name, _ := item["name"].(string)
		id, _ := item["id"].(string)
		ext_iface := ExternalInterface{id, name}
		external_interfaces[name] = ext_iface
	}

	//log.Noticef("%+v\n", external_interfaces)
	log.Debug("getting modules list OK\n")
	return nil, external_interfaces
}

/*-----------TABLES-------------*/

func TableEntryPUT(d *Dataplane, moduleId string, tableId string, entryId string, entryValue string) (error, ModuleTableEntry) {
	log.Infof("table entry PUT /modules/"+moduleId+"/tables/"+tableId+"/entries/"+entryId+" {%s,%s}\n", entryId, entryValue)

	req := map[string]interface{}{
		"key":   entryId,
		"value": entryValue,
	}
	var moduleTableEntry ModuleTableEntry
	err := d.sendObject("PUT", "/modules/"+moduleId+"/tables/"+tableId+"/entries/"+entryId, req, &moduleTableEntry)
	if err != nil {
		log.Warning(err)
		return err, moduleTableEntry
	}

	log.Debugf("table entry PUT /modules/"+moduleId+"/tables/"+tableId+"/entries/"+entryId+" {%s,%s} OK\n", moduleTableEntry.Key, moduleTableEntry.Value)

	return nil, moduleTableEntry
}

func TableEntryGET(d *Dataplane, moduleId string, tableId string, entryId string) (error, ModuleTableEntry) {
	log.Infof("table entry GET /modules/" + moduleId + "/tables/" + tableId + "/entries/" + entryId + "\n")

	req := map[string]interface{}{}
	var moduleTableEntry ModuleTableEntry
	err := d.sendObject("GET", "/modules/"+moduleId+"/tables/"+tableId+"/entries/"+entryId, req, &moduleTableEntry)
	if err != nil {
		log.Warning(err)
		return err, moduleTableEntry
	}

	log.Debugf("table entry GET /modules/"+moduleId+"/tables/"+tableId+"/entries/"+entryId+" {%s,%s} OK\n", moduleTableEntry.Key, moduleTableEntry.Value)

	return nil, moduleTableEntry
}

/*Not Working? Depending on Hover delete entryId on arrays? */
func TableEntryDELETE(d *Dataplane, moduleId string, tableId string, entryId string) (error, ModuleTableEntry) {
	log.Infof("table entry DELETE /modules/" + moduleId + "/tables/" + tableId + "/entries/" + entryId + "\n")

	req := map[string]interface{}{}
	var moduleTableEntry ModuleTableEntry
	err := d.sendObject("DELETE", "/modules/"+moduleId+"/tables/"+tableId+"/entries/"+entryId, req, &moduleTableEntry)
	if err != nil {
		log.Warning(err)
		return err, moduleTableEntry
	}

	log.Debugf("table entry DELETE /modules/" + moduleId + "/tables/" + tableId + "/entries/" + entryId + " OK\n")

	return nil, moduleTableEntry
}

/*
func TableEntryDELETE(d *Dataplane, moduleId string) (error, ModuleEntry) {
	fmt.Printf("deleting module %s\n", moduleId)

	req := map[string]interface{}{}
	var module ModuleEntry
	err := d.sendObject("DELETE", "/modules/"+moduleId, req, &module)
	if err != nil {
		return err, module
	}
	return nil, module
}
*/