/*
 * @time 2021/10/1 15:06
 * @version 1.00
 * @author huangsiyi
 *
 * 1.接收客户端 request，并将 request 中带的 header 写入 response header
 * 2.读取当前系统的环境变量中的 VERSION 配置，并写入 response header
 * 3.Server 端记录访问日志包括客户端 IP，HTTP 返回码，输出到 server 端的标准输出
 * 4.当访问 localhost/healthz 时，应返回200
 *
 */

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

const (
	StatusOK = 200
)

func main() {
	fmt.Printf("----- http server start -----")
	os.Setenv("VERSION", "1.00") // 设置当前系统环境VERSION为1.00

	mux := http.NewServeMux()
	mux.HandleFunc("/", defaultHandler)        // 当访问localhost时
	mux.HandleFunc("/healthz", healthzHandler) // 当访问 localhost/healthz 时

	err := http.ListenAndServe(":800", mux)
	if err != nil {
		log.Fatal(err)
	}
}

func defaultHandler(w http.ResponseWriter, req *http.Request) {

	// 1. 接收客户端 request，并将 request 中带的 header 写入 response header
	fmt.Fprintf(w, "Before write request header\n")
	for name, header := range req.Header {
		fmt.Fprintf(w, "Header[%v] : %v\n", name, header)
		w.Header().Add(name, strings.Join(header, ","))
	}
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "After write request header\n")
	for name, header := range req.Header {
		fmt.Fprintf(w, "Header[%v] : %v\n", name, header)
	}
	fmt.Fprintf(w, "\n")

	// 2. 读取当前系统的环境变量中的 VERSION 配置，并写入 response header
	fmt.Fprintf(w, "Version : %v\n", os.Getenv("VERSION"))
	w.Header().Add("VERSION", os.Getenv("VERSION"))
	fmt.Fprintf(w, "\n")

	// 3. Server 端记录访问日志包括客户端 IP，HTTP 返回码，输出到 server 端的标准输出
	w.WriteHeader(StatusOK)
	fmt.Fprintf(w, "IP : %v\n", req.Host)
	fmt.Fprintf(w, "Status : %v\n", StatusOK)
}

// 4.当访问 localhost/healthz 时，应返回200
func healthzHandler(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(StatusOK)
}
