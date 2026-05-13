package backend_skv

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
)

func doRequest(origin, path, hostname string) (int, error) {
	resp, err := http.Get(origin + path)
	if err != nil {
		log.Error().Str("hostname", hostname).Msg(fmt.Sprintf("error=%s", err))
		return -1, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error().Str("hostname", hostname).Msg(fmt.Sprintf("error=%s", err))
		return -1, err
	}

	val, err := strconv.Atoi(strings.TrimSpace(string(body)))
	if err != nil {
		log.Error().Str("hostname", hostname).Msg(fmt.Sprintf("error=%s", err))
		return -1, err
	}

	return val, nil
}

func DoCountSKV(origin, hostname string) (int, error) {
	return doRequest(origin, "/api/incr/counter", hostname)
}

func GetCountSKV(origin, hostname string) (int, error) {
	return doRequest(origin, "/api/get/counter", hostname)
}
