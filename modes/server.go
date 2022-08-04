package modes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/go-redis/redis"
	"github.com/uptrace/bunrouter"
	"github.com/uptrace/bunrouter/extra/reqlog"
)

var SERVER_SETTINGS ServerSettings

type requestHandler struct {
	client *redis.Client
	mutex  *sync.RWMutex
}

func StartServer(settings ServerSettings) error {
	SERVER_SETTINGS = settings

	router := initRouter()

	listenTo := fmt.Sprintf(":%d", SERVER_SETTINGS.IpPort)
	if err := http.ListenAndServe(listenTo, router); err != redis.Nil {
		return fmt.Errorf("server ended")
	}
	return nil
}

func initRouter() *bunrouter.Router {
	router := bunrouter.New(
		bunrouter.Use(reqlog.NewMiddleware()),
	)
	router.WithGroup("/:db/", func(group *bunrouter.Group) {
		group.GET("/", dbHandler)
		group.GET("/get", getHandler)
		group.GET("/set", setHandler)
		group.GET("/delete", deleteHandler)
	})

	return router
}

func dbHandler(w http.ResponseWriter, req bunrouter.Request) error {
	params := req.Params().Map()
	db, _ := strconv.Atoi(params["db"])

	helpString := fmt.Sprintf("Current Database: %d\n\n", db)
	homeHandlerHelp(&helpString)

	w.Write([]byte(helpString))
	return nil
}

func getHandler(w http.ResponseWriter, req bunrouter.Request) error {
	queries := req.URL.Query()

	key := queries.Get("key")
	db, _ := strconv.Atoi(queries.Get("db"))

	handler := initHandler(SERVER_SETTINGS.DbAddress(), db)
	if key != "" {
		value, err := handler.readOne(key)
		if err != nil {
			return fmt.Errorf("redis client cannot get %s", key)
		} else {
			w.Write([]byte(value))
		}
	} else {
		keys, err := handler.readAll()
		if err != nil {
			return fmt.Errorf("redis client cannot get all keys")
		} else {
			keysMarshalled, _ := json.Marshal(keys)
			w.Write(keysMarshalled)
		}
	}
	return nil
}

func setHandler(w http.ResponseWriter, req bunrouter.Request) error {
	queries := req.URL.Query()

	key := queries.Get("key")
	value := queries.Get("value")
	db, _ := strconv.Atoi(queries.Get("db"))

	handler := initHandler(SERVER_SETTINGS.DbAddress(), db)
	if key == "" {
		return fmt.Errorf("no key given")
	}
	if err := handler.writeOne(key, value); err != nil {
		return fmt.Errorf("redis client cannot set %s to %s", key, value)
	}
	return nil
}

func deleteHandler(w http.ResponseWriter, req bunrouter.Request) error {
	queries := req.URL.Query()

	key := queries.Get("key")
	db, _ := strconv.Atoi(queries.Get("db"))

	handler := initHandler(SERVER_SETTINGS.DbAddress(), db)
	if key == "" {
		return fmt.Errorf("no key given")
	}
	if err := handler.deleteOne(key); err != nil {
		return fmt.Errorf("redis client cannot find %s", key)
	}
	return nil
}

func initHandler(addr string, db int) requestHandler {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       db,
	})
	if _, err := redisClient.Ping().Result(); err != nil {
		fmt.Println("redis client cannot be started")
	}
	var handler requestHandler = requestHandler{
		client: redisClient,
		mutex:  &sync.RWMutex{},
	}
	return handler
}

func (h *requestHandler) readOne(key string) (string, error) {
	h.mutex.RLock()
	value, err := h.client.Get(key).Result()
	h.mutex.RUnlock()

	return value, err
}

func (h *requestHandler) readAll() ([]string, error) {
	h.mutex.RLock()
	keys, err := h.client.Keys("*").Result()
	h.mutex.RUnlock()

	return keys, err
}

func (h *requestHandler) writeOne(key string, value string) error {
	minutes := time.Duration(SERVER_SETTINGS.TTLMinutes)

	h.mutex.Lock()
	_, err := h.client.Set(key, value, minutes*time.Minute).Result()
	h.mutex.Unlock()

	return err
}

func (h *requestHandler) deleteOne(key string) error {
	h.mutex.Lock()
	_, err := h.client.Del(key).Result()
	h.mutex.Unlock()

	return err
}

func homeHandlerHelp(helpString *string) {
	*helpString += fmt.Sprintf("%7s | %25s | %s\n", "path", "description", "query")
	*helpString += fmt.Sprintf("---------------------------------------------------------------\n")
	*helpString += fmt.Sprintf("%7s | %25s | %s\n", "/get", "gets all keys", "")
	*helpString += fmt.Sprintf("%7s | %25s | %s\n", "/get", "gets the value by a key", "key string")
	*helpString += fmt.Sprintf("%7s | %25s | %s\n", "/set", "updates a key-value pair", "key string, value string")
	*helpString += fmt.Sprintf("%7s | %25s | %s\n", "/delete", "deletes a key-value pair", "key string")
}
