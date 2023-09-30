package httpclient

import (
	"net/http"
	"net/http/httputil"
	"time"

	"go.uber.org/zap"
)

type httpLogger struct {
	log *zap.Logger
}

func (l *httpLogger) LogRequest(req *http.Request) {
	body, err := httputil.DumpRequest(req, true)
	if err != nil {
		l.log.Error("[HTTP] Failed to dump request", zap.Error(err))
	} else {
		l.log.Info(
			"[HTTP] External api request",
			zap.String("path", req.URL.String()),
			zap.String("method", req.Method),
			zap.String("request", string(body)),
		)
	}
}

func (l *httpLogger) LogResponse(req *http.Request, res *http.Response, err error, duration time.Duration) {
	if err != nil {
		l.log.Error("[HTTP] Error log", zap.Error(err))
	}
	if res == nil {
		return
	}
	body, err := httputil.DumpResponse(res, true)
	if err != nil {
		l.log.Error("[HTTP] Failed to dump response", zap.Error(err))
	} else {
		l.log.Info(
			"[HTTP] External api response",
			zap.String("path", req.URL.String()),
			zap.String("method", req.Method),
			zap.Int("status", res.StatusCode),
			zap.Duration("duration", duration),
			zap.Error(err),
			zap.String("response", string(body)),
		)
	}
}
