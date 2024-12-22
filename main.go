package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nicolas-vp/groupcache"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var myCache *groupcache.Group

const cacheName = "test"
const keyName = "123"

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Параметры запуска: порт_http нода_кеша ноды_кеша")
		fmt.Println("Пример для первой ноды: 3001 localhost:4001 http://localhost:4001,http://localhost:4002")
		fmt.Println("Пример для второй ноды: 3002 localhost:4002 http://localhost:4002,http://localhost:4001")
		return
	}

	startCache(os.Args[3], os.Args[2])

	router := mux.NewRouter()
	router.HandleFunc("/write", writeHandler).Methods("GET")
	router.HandleFunc("/read", readHandler).Methods("GET")

	port := os.Args[1]
	log.Println("server listen on port: " + port)

	myCache = createCache(cacheName).WithUpdateOtherNodes()

	err := http.ListenAndServe(":"+port, router)

	if err != nil {
		log.Fatal(err)
	}
}

func startCache(nodes string, host string) {
	p := strings.Split(nodes, ",")

	// вот тут обратите внимание, что p[0] это первая нода в списке, именно она и инциализируется и становится Self
	pool := groupcache.NewHTTPPoolOpts(p[0], &groupcache.HTTPPoolOptions{Replicas: 5000})
	pool.Set(p...)

	server := http.Server{
		Addr:    host,
		Handler: pool,
	}

	go func() {
		log.Println("Кеш инициализирован на хосте: " + host)
		if err := server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()
}

func createCache(cacheName string) *groupcache.Group {
	fillFunction := groupcache.GetterFunc(
		func(ctx context.Context, key string, dest groupcache.Sink) error {
			// не будем использовать рефилл
			return nil
		})
	return groupcache.NewGroup(cacheName, 1000, fillFunction) //.WithUpdateOtherNodes()
}

func readHandler(writer http.ResponseWriter, _ *http.Request) {
	var s []byte
	err := myCache.Get(context.Background(), keyName, groupcache.AllocatingByteSliceSink(&s))
	if err != nil {
		fmt.Println(err)
	}
	resultString := string(s)
	writeJsonToHTTP(writer, resultString)
	fmt.Printf("GET %v\n", resultString)
}

func writeHandler(writer http.ResponseWriter, _ *http.Request) {
	// будем записывать в кеш рандомное значение в виде строки
	resultString := strconv.Itoa(rand.Int())

	err := myCache.Set(context.Background(), keyName, []byte(resultString), time.Now().Add(time.Second*1000), true)
	if err != nil {
		fmt.Println(err)
	}
	writeJsonToHTTP(writer, resultString)
	fmt.Printf("WRITE %v\n", resultString)
}

func writeJsonToHTTP(writer http.ResponseWriter, resultString string) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)

	_, err := writer.Write([]byte(fmt.Sprintf("{\"value\": \"%v\"}", resultString)))
	if err != nil {
		fmt.Println(err)
	}
}
