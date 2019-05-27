package main

import (
	"log"
	"net/http"
	"time"

	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

const port = ":8009"

func main() {
	cfg := config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Second,
		},
	}

	t, closer, err := cfg.New("main-function", config.Logger(jaeger.StdLogger))
	if err != nil {
		log.Fatal(err)
	}
	defer closer.Close()
	opentracing.SetGlobalTracer(t)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		p := opentracing.GlobalTracer().StartSpan("index")
		defer p.Finish()
		log.Print("Log has been printed")
		c := opentracing.GlobalTracer().StartSpan("index-handler", opentracing.ChildOf(p.Context()))
		defer c.Finish()
		w.Write([]byte("Hello World!"))
	})

	log.Printf("Port listened on %s\n", port)
	http.ListenAndServe(port, nethttp.Middleware(opentracing.GlobalTracer(), http.DefaultServeMux))
}
