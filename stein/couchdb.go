package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/beatgammit/semver"
	"github.com/beatgammit/stein"
	"net/http"
)

type CouchDB struct {
	addr, database, design string
	version                semver.Semver
	user, pass             string
}

func NewCouchDB(addr, database, user, pass string) (DB, error) {
	// it's possible the user wants https:// or it's behind
	// a proxy, so let them specify that if they like
	if addr[:4] != "http" {
		addr = "http://" + addr
	}
	if addr[len(addr)-1] != '/' {
		addr += "/"
	}

	resp, err := http.Get(addr)
	if err != nil {
		return nil, err
	}
	var m map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&m)
	if err != nil {
		return nil, err
	} else if _, ok := m["version"]; !ok {
		return nil, fmt.Errorf("Invalid CouchDB server")
	}

	db := CouchDB{
		addr:     addr,
		database: database,
		user:     user,
		pass:     pass,
		design:   "stein",
	}
	db.version, err = semver.Parse(m["version"].(string))
	if err != nil {
		return nil, fmt.Errorf("Error parsing version: %s", err)
	}
	return db, db.init()
}

func (db CouchDB) genByTypeView(project string) (string, map[string]string) {
	return project + "_by_type", map[string]string{
		"map": fmt.Sprintf(`function(doc) {
           if (doc.project == '%s' && typeof doc.Type === 'string') {
               emit(doc.Type, doc._id);
           }
	   }`, project),
		"reduce": `function (keys, values, rereduce) {
           if (rereduce) {
               return values.reduce(function (p, o) {
                   var k;

                   for (k in o) {
                       if (isNaN(o[k])) {
                           continue;
                       }
                       if (isNaN(p[k])) {
                           p[k] = 0;
                       }
                       p[k] += o[k];
                   }
                   return p;
               }, {});
           }

		   return keys.reduce(function (p, key) {
			   var k = key[0];
			   p[k] = k in p ? p[k] + 1 : 1;
			   return p;
		   }, {});
	   }`,
	}
}

// init ensures that the database is configured correctly:
// - creates database
// - creates views
func (db CouchDB) init() error {
	req, err := http.NewRequest("PUT", db.addr+db.database, nil)
	if err != nil {
		return fmt.Errorf("Error sending DB create request")
	}
	if db.user != "" {
		req.SetBasicAuth(db.user, db.pass)
	}
	if resp, err := http.DefaultClient.Do(req); err != nil {
		return err
	} else {
		// TODO: check for errors
		resp.Body.Close()
	}

	designDocUrl := db.addr + db.database + "/_design/" + db.design

	resp, err := http.Get(designDocUrl)
	if err != nil {
		return err
	}
	var m map[string]interface{}
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		// design doc already exists, so update it
		if err = json.NewDecoder(resp.Body).Decode(&m); err != nil {
			return err
		}
	} else {
		m = make(map[string]interface{})
	}
	if _, ok := m["views"]; !ok {
		m["views"] = make(map[string]interface{})
		m["language"] = "javascript"
	}
	views := m["views"].(map[string]interface{})

	if projects, err := db.GetProjects(); err == nil {
		for _, project := range projects {
			viewName, view := db.genByTypeView(project)
			views[viewName] = view
		}
	}

	// by_project maps projects to documents
	// use reduce=false to get all documents, or
	// reduce=true to get counts per project
	views["by_project"] = map[string]interface{}{
		"map": `function(doc) {
		   emit(doc.project, doc);
	   }`,
		"reduce": `function (keys, values, rereduce) {
           if (rereduce) {
               return values.reduce(function (p, o) {
                   var k;

                   for (k in o) {
                       if (isNaN(o[k])) {
                           continue;
                       }
                       if (isNaN(p[k])) {
                           p[k] = 0;
                       }
                       p[k] += o[k];
                   }
                   return p;
               }, {});
           }

		   return keys.reduce(function (p, key) {
			   var k = key[0];
			   p[k] = k in p ? p[k] + 1 : 1;
			   return p;
		   }, {});
	   }`,
	}

	b, err := json.Marshal(m)
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer(b)
	if rev, ok := m["_rev"].(string); ok {
		req, err = http.NewRequest("PUT", designDocUrl+"?rev="+rev, buf)
	} else {
		req, err = http.NewRequest("PUT", designDocUrl, buf)
	}
	if db.user != "" {
		req.SetBasicAuth(db.user, db.pass)
	}
	if err != nil {
		return err
	}
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return err
	} else if resp.StatusCode >= 400 {
		fmt.Println(designDocUrl)
		return fmt.Errorf("Error updating view: %d", resp.StatusCode)
	}
	return nil
}

func (db CouchDB) GetProjects() ([]string, error) {
	resp, err := http.Get(db.addr + db.database + "/_design/" + db.design + "/_view/by_project")
	if err != nil {
		return nil, err
	}
	var m map[string]interface{}
	if err = json.NewDecoder(resp.Body).Decode(&m); err != nil {
		return nil, err
	}

	projects := []string{}
	rows, ok := m["rows"].([]interface{})
	if !ok || len(rows) == 0 {
		// no projects yet
		return projects, nil
	}
	match := rows[0].(map[string]interface{})
	counts := match["value"].(map[string]interface{})
	for project := range counts {
		projects = append(projects, project)
	}
	return projects, nil
}

func (db CouchDB) GetTests(project string) ([]string, error) {
	resp, err := http.Get(db.addr + db.database + "/_design/" + db.design + "/_view/by_project?reduce=false")
	if err != nil {
		return nil, err
	}
	var m map[string]interface{}
	if err = json.NewDecoder(resp.Body).Decode(&m); err != nil {
		return nil, err
	}

	var tests []string
	rows := m["rows"].([]interface{})
	for _, row := range rows {
		// TODO: make this safer
		val := row.(map[string]interface{})["value"].(map[string]interface{})
		id := val["_id"].(string)
		tests = append(tests, id)
	}
	return tests, nil
}

func (db CouchDB) GetTestTypes(project string) ([]string, error) {
	urlPath := fmt.Sprintf("%s%s/_design/%s/_view/%s_by_type?group=true", db.addr, db.database, db.design, project)
	resp, err := http.Get(urlPath)
	if err != nil {
		return nil, err
	}
	var m map[string]interface{}
	if err = json.NewDecoder(resp.Body).Decode(&m); err != nil {
		return nil, err
	}

	var testTypes []string
	for _, row := range m["rows"].([]interface{}) {
		// TODO: make this safer
		key := row.(map[string]interface{})["key"].(string)
		testTypes = append(testTypes, key)
	}
	return testTypes, nil
}

func (db CouchDB) GetTestsByType(project, typ string) ([]string, error) {
	urlPath := fmt.Sprintf("%s%s/_design/%s/_view/%s_by_type?key=\"%s\"&reduce=false", db.addr, db.database, db.design, project, typ)
	resp, err := http.Get(urlPath)
	if err != nil {
		return nil, err
	}
	var m map[string]interface{}
	if err = json.NewDecoder(resp.Body).Decode(&m); err != nil {
		return nil, err
	}

	var tests []string
	for _, row := range m["rows"].([]interface{}) {
		// TODO: make this safer
		id := row.(map[string]interface{})["id"].(string)
		tests = append(tests, id)
	}
	return tests, nil
}

func (db CouchDB) GetTest(project, test string) (*stein.Suite, error) {
	resp, err := http.Get(db.addr + db.database + "/" + test)
	if err != nil {
		return nil, err
	}

	var s stein.Suite
	return &s, json.NewDecoder(resp.Body).Decode(&s)
}

func (db CouchDB) Save(project, test string, s *stein.Suite) error {
	b, err := json.Marshal(s)
	if err != nil {
		return err
	}
	var m map[string]interface{}
	err = json.Unmarshal(b, &m)
	if err != nil {
		return err
	}
	m["project"] = project
	b, err = json.Marshal(m)
	if err != nil {
		return err
	}

	testAddr := db.addr + db.database + "/" + test

	var rev string
	resp, err := http.Get(testAddr)
	if err != nil {
		return err
	} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var cur map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&cur)
		if err != nil {
			return err
		}
		rev = cur["_rev"].(string)
	}

	var req *http.Request
	buf := bytes.NewBuffer(b)
	if rev == "" {
		// create a new one
		req, err = http.NewRequest("PUT", testAddr, buf)
	} else {
		// update existing
		req, err = http.NewRequest("PUT", testAddr+"?rev="+rev, buf)
	}
	if db.user != "" {
		req.SetBasicAuth(db.user, db.pass)
	}
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return err
	} else if resp.StatusCode >= 400 {
		return fmt.Errorf("Error updating document: %d", resp.StatusCode)
	}
	resp.Body.Close()
	return nil
}
