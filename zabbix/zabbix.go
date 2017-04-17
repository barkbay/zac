package zabbix

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"

	"github.com/AlekSi/zabbix"
)

type Zabbix struct {
	r8AlertServer  string // http://alert.my.company.com
	zabbixServer   string // https://zabbix.my.company.com
	zabbixUser     string
	zabbixPassword string
	rateURL        string // URL where the rate service can be reached
	api            *zabbix.API
	debug          bool
}

func NewZabbix(zabbixServer string, zabbixUser string, zabbixPassword string, rateURL string) *Zabbix {

	debug := false
	if debugVar := os.Getenv("DEBUG"); debugVar != "" {
		debug = true
	}

	return &Zabbix{
		zabbixServer:   zabbixServer,
		zabbixUser:     zabbixUser,
		zabbixPassword: zabbixPassword,
		api:            zabbix.NewAPI(zabbixServer + "/api_jsonrpc.php"),
		debug:          debug,
		rateURL:        rateURL,
	}
}

func (z *Zabbix) NewOrUpdateMonitoring(namespace string) {

	scenario := "Warning rate in " + namespace

	// Login
	z.api.Login(z.zabbixUser, z.zabbixPassword)
	// Get the Zabbix server
	zbxSrv := z.getZabbixServer()

	// Compute URL
	u, parseErr := url.Parse(z.rateURL)
	if parseErr != nil {
		panic(parseErr.Error())
	}
	relative, _ := url.Parse(namespace)
	sanitizedRateURL := u.ResolveReference(relative).String()

	if zbxSrv != nil {
		// Check if Web Scenario already exist
		httptestid, exists := z.getWebScenario(scenario)

		// Compute Web Senario steps
		checkName := "Retrieve rate for " + namespace
		step1 := map[string]string{"name": checkName, "url": sanitizedRateURL, "status_codes": "404,503", "no": "1"}
		steps := []map[string]string{step1}
		parameters := map[string]interface{}{"name": scenario, "hostid": zbxSrv.HostId, "steps": steps}
		if z.debug {
			b, _ := json.MarshalIndent(parameters, "", "    ")
			fmt.Printf("[->ZBX] %+v\n", string(b))
		}

		if !exists {
			createRes, createErr := z.api.Call("httptest.create", zabbix.Params{
				"name":   scenario,
				"hostid": zbxSrv.HostId,
				"delay":  10,
				"steps":  steps,
				"output": "shorten"})
			if createErr != nil {
				panic(createErr.Error())
			}
			if z.debug {
				createB, _ := json.MarshalIndent(createRes, "", "    ")
				fmt.Printf("[C<-ZBX] %+v\n", string(createB))
			}
		} else {
			updateRes, updateErr := z.api.Call("httptest.update", zabbix.Params{
				"name":       scenario,
				"httptestid": httptestid,
				"hostid":     zbxSrv.HostId,
				"delay":      10,
				"steps":      steps,
				"output":     "shorten"})
			if updateErr != nil {
				panic(updateErr.Error())
			}
			if z.debug {
				updateB, _ := json.MarshalIndent(updateRes, "", "    ")
				fmt.Printf("[U<-ZBX] %+v\n", string(updateB))
			}
		}
		// Ok now create trigger
		expression := "{Zabbix server:web.test.fail[" + scenario + "].last()}<>0"
		description := "Taux d'erreur trop haut pour " + namespace
		triggerid, exists := z.getTrigger(description)
		if !exists {
			z.createTrigger(namespace, expression, description)
		} else {
			z.updateTrigger(triggerid, namespace, expression, description)
		}
	}
}

func (z *Zabbix) getZabbixServer() *zabbix.Host {
	filter := map[string]string{"host": "Zabbix server"}
	res, err := z.api.HostsGet(zabbix.Params{"output": "extend", "filter": filter})
	if err != nil {
		panic(err.Error())
	}
	if len(res) > 0 {
		return &res[0]
	}
	return nil
}

func (z *Zabbix) getWebScenario(ws string) (string, bool) {
	filter := map[string]string{"name": ws}
	res, err := z.api.Call("httptest.get", zabbix.Params{"output": "extend", "filter": filter})
	if err != nil {
		panic(err.Error())
	}
	f := res.Result.([]interface{})
	if len(f) < 1 {
		return "", false
	}
	d := f[0].(map[string]interface{})
	if val, ok := d["httptestid"]; ok {
		return val.(string), true
	}
	return "", false
}

func (z *Zabbix) getTrigger(description string) (string, bool) {
	filter := map[string]string{"description": description}
	res, err := z.api.Call("trigger.get", zabbix.Params{"output": "extend", "filter": filter})
	if err != nil {
		panic(err.Error())
	}
	if z.debug {
		updateB, _ := json.MarshalIndent(res, "", "    ")
		fmt.Printf("[GT<-ZBX] %+v\n", string(updateB))
	}
	f := res.Result.([]interface{})
	if len(f) < 1 {
		return "", false
	}
	d := f[0].(map[string]interface{})
	if val, ok := d["triggerid"]; ok {
		return val.(string), true
	}
	return "", false
}

func (z *Zabbix) createTrigger(namespace string, expression string, description string) {
	foo, err := z.api.Call("trigger.create",
		zabbix.Params{
			"expression":  expression,
			"description": description,
			"priority":    4})
	if err != nil {
		panic(err.Error())
	}
	if z.debug {
		updateB, _ := json.MarshalIndent(foo, "", "    ")
		fmt.Printf("[CT<-ZBX] %+v\n", string(updateB))
	}
}

func (z *Zabbix) updateTrigger(triggerid string, namespace string, expression string, description string) {
	foo, err := z.api.Call("trigger.update",
		zabbix.Params{
			"triggerid":   triggerid,
			"expression":  expression,
			"description": description,
			"priority":    4})
	if err != nil {
		panic(err.Error())
	}
	if z.debug {
		updateB, _ := json.MarshalIndent(foo, "", "    ")
		fmt.Printf("[UT<-ZBX] %+v\n", string(updateB))
	}
}
