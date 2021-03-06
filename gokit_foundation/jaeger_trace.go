package gokit_foundation

import (
	"github.com/opentracing/opentracing-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"io"
	"os"
)

func InitTracer(svcName string) (io.Closer, error) {
	setJaegerEnvConf_ForTest()
	cfg, err := jaegercfg.FromEnv()
	if err != nil {
		return nil, err
	}

	cfg.ServiceName = svcName
	tracer, closer, err := cfg.NewTracer()
	if err != nil {
		return nil, err
	}
	opentracing.SetGlobalTracer(tracer)
	return closer, nil
}

// 实际环境中通过其他方式配置env
func setJaegerEnvConf_ForTest() {
	// 完整的env配置项参考： https://github.com/jaegertracing/jaeger-client-go#environment-variables
	//os.Setenv("JAEGER_SERVICE_NAME", "TestInitTracer")
	//os.Setenv("JAEGER_AGENT_HOST", "5.181.135.29")
	//os.Setenv("JAEGER_AGENT_PORT", "6831")
	os.Setenv("JAEGER_ENDPOINT", "http://localhost:14268/api/traces") // 若指定，则agent的host和port被忽略，配置的是collector的HTTP地址，如 http://jaeger-collector:14268/api/traces
	//os.Setenv("JAEGER_USER", "")     // collector认证
	//os.Setenv("JAEGER_PASSWORD", "") // collector认证
	//os.Setenv("JAEGER_REPORTER_LOG_SPANS", "")      // true or false
	//os.Setenv("JAEGER_REPORTER_MAX_QUEUE_SIZE", "") // The reporter's maximum queue size (default 100).
	os.Setenv("JAEGER_SAMPLER_TYPE", "const") // The sampler type: remote, const, probabilistic, ratelimiting (default remote).
	os.Setenv("JAEGER_SAMPLER_PARAM", "1")    // number
	//os.Setenv("JAEGER_SAMPLING_ENDPOINT", "") // when using sampler type remote (default http://127.0.0.1:5778/sampling).
	//os.Setenv("JAEGER_TAGS", "")              // 逗号分隔的k=v格式，会被添加到此服务的所有上报的spans上，例如：svc=user,level=important
	//os.Setenv("JAEGER_DISABLED", "")          // bool, 如果true，全局就会使用一个空的tracer `opentracing.NoopTracer` (default false).
	/*
		默认情况下，client会通过UDP协议发送span数据到 localhost:6831，我们只需要配置 JAEGER_AGENT_HOST和JAEGER_AGENT_PORT来发往指定的agent；
		但也可以指定JAEGER_ENDPOINT，让client直接把span数据发往collector
	*/
}
