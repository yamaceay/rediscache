package modes

import (
	"bytes"
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
		return fmt.Errorf("server terminated: %s", err)
	}
	return nil
}

func initRouter() *bunrouter.Router {
	router := bunrouter.New(
		bunrouter.Use(reqlog.NewMiddleware()),
	)
	router.WithGroup("/:db/", func(group *bunrouter.Group) {
		group.GET("/help", dbHandler)
		group.GET("/get", getHandler)
		group.GET("/set", setHandler)
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

	params := req.Params().Map()
	db, _ := strconv.Atoi(params["db"])

	handler := newHandler(SERVER_SETTINGS.DbAddress(), db)
	if key == "" {
		if keys, err := handler.readAll(); err != nil {
			return fmt.Errorf("keys cannot be read: %s", err)
		} else {
			keysMarshalled, _ := json.MarshalIndent(keys, "", "  ")
			w.Write(keysMarshalled)
		}
		return nil
	}

	if value, err := handler.readOne(key); err != nil {
		return fmt.Errorf("object with key %s cannot be read: %s", key, err)
	} else {
		var prettyJSON bytes.Buffer
		if err := json.Indent(&prettyJSON, []byte(value), "", "  "); err != nil {
			w.Write([]byte(value))
		} else {
			w.Write(prettyJSON.Bytes())
		}
	}
	return nil
}

func setHandler(w http.ResponseWriter, req bunrouter.Request) error {
	queries := req.URL.Query()
	key := queries.Get("key")
	value := queries.Get("value")

	params := req.Params().Map()
	db, _ := strconv.Atoi(params["db"])

	handler := newHandler(SERVER_SETTINGS.DbAddress(), db)
	if key == "" {
		return fmt.Errorf("no key given")
	}
	if err := handler.writeOne(key, value); err != nil {
		if value == "" {
			return fmt.Errorf("object with key %s couldn't be deleted: %s", key, err)
		} else {
			return fmt.Errorf("object with key %s couldn't be set to %s: %s", key, value, err)
		}
	}
	return nil
}

func newHandler(addr string, db int) *requestHandler {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       db,
	})
	if _, err := redisClient.Ping().Result(); err != nil {
		fmt.Printf("client cannot be started: %s no:%d", addr, db)
	}
	var handler requestHandler = requestHandler{
		client: redisClient,
		mutex:  &sync.RWMutex{},
	}
	return &handler
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

func homeHandlerHelp(helpString *string) {
	*helpString += fmt.Sprintf("%7s | %25s | %s\n", "path", "description", "query")
	*helpString += fmt.Sprintf("---------------------------------------------------------------\n")
	*helpString += fmt.Sprintf("%7s | %25s | %s\n", "/get", "gets all keys", "")
	*helpString += fmt.Sprintf("%7s | %25s | %s\n", "/get", "gets the value by a key", "key string")
	*helpString += fmt.Sprintf("%7s | %25s | %s\n", "/set", "deletes a key-value pair", "key string")
	*helpString += fmt.Sprintf("%7s | %25s | %s\n", "/set", "updates a key-value pair", "key string, value string")
}
