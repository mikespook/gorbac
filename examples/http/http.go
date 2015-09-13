// possum & gorbac example
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/mikespook/gorbac"
	"github.com/mikespook/possum"
	"github.com/mikespook/possum/router"
	"github.com/mikespook/possum/view"
)

const addr = "127.0.0.1:12345"

var rbac = gorbac.New()

func postHandler(ctx *possum.Context) error {
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		return err
	}
	defer ctx.Request.Body.Close()
	var m gorbac.Map
	if err := json.Unmarshal(body, &m); err != nil {
		return err
	}
	rbac = gorbac.Restore(m)
	ctx.Response.Data = rbac.Dump()
	return nil
}

func getHandler(ctx *possum.Context) error {
	name := ctx.Request.Form.Get("name")
	if name == "" {
		ctx.Response.Data = rbac.Dump()
	} else {
		ctx.Response.Data = gorbac.RoleToMap(rbac.Get(name))
	}
	return nil
}

func putHandler(ctx *possum.Context) error {
	name := ctx.Request.Form.Get("name")
	permissions := ctx.Request.Form["permissions"]
	parents := ctx.Request.Form["parents"]
	rbac.Set(name, permissions, parents)
	ctx.Response.Data = gorbac.RoleToMap(rbac.Get(name))
	return nil
}

func deleteHandler(ctx *possum.Context) error {
	name := ctx.Request.Form.Get("name")
	role := rbac.Get(name)
	rbac.Remove(name)
	ctx.Response.Data = gorbac.RoleToMap(role)
	return nil
}

func patchHandler(ctx *possum.Context) error {
	name := ctx.Request.Form.Get("name")
	permissions := ctx.Request.Form["permissions"]
	parents := ctx.Request.Form["parents"]
	rbac.Add(name, permissions, parents)
	ctx.Response.Data = gorbac.RoleToMap(rbac.Get(name))
	return nil
}

func rbacHandler(ctx *possum.Context) error {
	switch ctx.Request.Method {
	case "PATCH":
		return patchHandler(ctx)
	case "GET":
		return getHandler(ctx)
	case "POST":
		return postHandler(ctx)
	case "DELETE":
		return deleteHandler(ctx)
	case "PUT":
		return putHandler(ctx)
	}
	return nil
}

func isGrantedHandler(ctx *possum.Context) error {
	ctx.Response.Header().Set("Test", "Hello world")
	name := ctx.Request.Form.Get("name")
	permission := ctx.Request.Form.Get("permission")
	if rbac.IsGranted(name, permission, nil) {
		ctx.Response.Status = http.StatusOK
		ctx.Response.Data = true
		return nil
	}
	ctx.Response.Status = http.StatusForbidden
	ctx.Response.Data = false
	return nil
}

func main() {
	mux := possum.NewServerMux()

	mux.PreRequest = func(ctx *possum.Context) error {
		host, _, err := net.SplitHostPort(ctx.Request.RemoteAddr)
		if err != nil {
			return err
		}
		if host != "127.0.0.1" {
			return possum.NewError(http.StatusForbidden, "Localhost only")
		}
		return nil
	}

	mux.PostResponse = func(ctx *possum.Context) error {
		fmt.Printf("[%d] %s:%s \"%s\"", ctx.Response.Status,
			ctx.Request.RemoteAddr, ctx.Request.Method,
			ctx.Request.URL.String())
		return nil
	}
	mux.HandleFunc(router.Simple("/rbac"), rbacHandler, view.Json(view.CharSetUTF8))
	mux.HandleFunc(router.Simple("/isgranted"), isGrantedHandler, view.Json(view.CharSetUTF8))
	fmt.Printf("[%s] %s\n", time.Now(), addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		fmt.Println(err)
		return
	}
}
