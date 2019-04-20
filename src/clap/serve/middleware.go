package serve

import (
	"clap/staging/TBLogger"
	"net/http"
	"time"
)

func LogApiAccess(httpHandler http.HandlerFunc)http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		TBLogger.TbLogger.Info("ApiAccess:",
			"host:",r.Host,
			"AccessApi:",r.URL.Path,
			"Time:",time.Now().Format("2006-01-02 15:04:05"),
		)
		httpHandler(w,r)
	}
}

type DownFileAccessMid struct {
	handlerFunc http.Handler
}

func (dfm DownFileAccessMid)ServeHTTP(w http.ResponseWriter, r *http.Request){
	TBLogger.TbLogger.Info("ApiAccess",
		"host:",r.Host,
		"AccessApi:",r.URL.Path,
		"Method",r.Method,
	)
	dfm.handlerFunc.ServeHTTP(w,r)
}

func LogDownFileAccess(handlerFunc http.Handler) http.Handler{
	var dfm DownFileAccessMid
	dfm.handlerFunc = handlerFunc
	return dfm
}