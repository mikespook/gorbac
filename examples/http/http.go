// possum & gorbac example
package main

import (
	"encoding/json"
	"fmt"
	"github.com/mikespook/gorbac"
	"github.com/mikespook/possum"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

const addr = "127.0.0.1:12345"

var rbac *myRbac

type myRbac struct {
	*gorbac.Rbac
}

func (rbac *myRbac) Post(w http.ResponseWriter, r *http.Request) (status int, data interface{}) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return http.StatusInternalServerError, err.Error()
	}
	defer r.Body.Close()
	var m gorbac.Map
	if err := json.Unmarshal(body, &m); err != nil {
		return http.StatusInternalServerError, err.Error()
	}
	rbac.Rbac = gorbac.Restore(m)
	return http.StatusOK, rbac.Rbac.Dump()
}

func (rbac *myRbac) Get(w http.ResponseWriter, r *http.Request) (status int, data interface{}) {
	name := r.Form.Get("name")
	if name == "" {
		return http.StatusOK, rbac.Rbac.Dump()
	}
	return http.StatusOK, gorbac.RoleToMap(rbac.Rbac.Get(name))
}

func (rbac *myRbac) Put(w http.ResponseWriter, r *http.Request) (status int, data interface{}) {
	name := r.Form.Get("name")
	permissions := r.Form["permissions"]
	parents := r.Form["parents"]
	rbac.Rbac.Set(name, permissions, parents)
	return http.StatusOK, gorbac.RoleToMap(rbac.Rbac.Get(name))
}

func (rbac *myRbac) Delete(w http.ResponseWriter, r *http.Request) (status int, data interface{}) {
	name := r.Form.Get("name")
	role := rbac.Rbac.Get(name)
	rbac.Rbac.Remove(name)
	return http.StatusOK, gorbac.RoleToMap(role)
}

func (rbac *myRbac) Patch(w http.ResponseWriter, r *http.Request) (status int, data interface{}) {
	name := r.Form.Get("name")
	permissions := r.Form["permissions"]
	parents := r.Form["parents"]
	rbac.Rbac.Add(name, permissions, parents)
	return http.StatusOK, gorbac.RoleToMap(rbac.Rbac.Get(name))
}

func isGranded(w http.ResponseWriter, r *http.Request) (status int, data interface{}) {
	name := r.Form.Get("name")
	permission := r.Form.Get("permission")
	if rbac.Rbac.IsGranted(name, permission, nil) {
		return http.StatusOK, gorbac.RoleToMap(rbac.Rbac.Get(name))
	}
	return http.StatusForbidden, nil
}

func main() {
	rbac = &myRbac{Rbac: gorbac.New()}

	h := possum.NewHandler()
	h.PreHandler = func(r *http.Request) (int, error) {
		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			return http.StatusInternalServerError, err
		}
		if host != "127.0.0.1" {
			return http.StatusForbidden, fmt.Errorf("Localhost only")
		}
		return http.StatusOK, nil
	}

	h.PostHandler = func(r *http.Request, status int, data interface{}) {
		fmt.Printf("[%s] %s %s \"%s\" %d\n", time.Now(), r.RemoteAddr, r.Method, r.URL.String(), status)
	}

	if err := h.AddResource("/rbac", rbac); err != nil {
		fmt.Println(err)
		return
	}

	h.AddRPC("/isgranted", isGranted)
	fmt.Printf("[%s] %s\n", time.Now(), addr)
	if err := http.ListenAndServe(addr, h); err != nil {
		fmt.Println(err)
		return
	}
}
