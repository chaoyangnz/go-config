[![Travis Build Status](https://travis-ci.org/mexisme/go-config.svg)](https://travis-ci.org/mexisme/go-config)

# go-config

This is a relatively simple wrapper around configuring projects that I've been using as a pattern for a while.

I've pulled it out of one of the original projects, generalised it better, and gave it a better API.

## Why?

The problem I was facing was two-fold:
- Configuring logging in a way that I can easily manipulate it from config (args, env-vars, etc)
- Adding new config items from anywhere in my code; i.e. have it tied more-closely to the piece of code that needs it, rather than in `main.go`
and then having to pass it through.

Tools like Flags and Viper support me really well in configuration, but I was still having to write the same pattern repetitively, and that
seemed a good candidate for abstraction.

Also, whilst they are technically separated, configuring logging seemed to have some overlaps in need:  they both need to be configured early
in the run-time of a tool, and there was interdependence in a few places, e.g. enabling debug logging, sending log-messages (and errors) somewhere.

# How to use it?

## By Example

The [godoc doc's](https://godoc.org/github.com/mexisme/go-config) are the definitive guide, but a few simple exmples are:

Importing it:
```go
import "github.com/mexisme/go-config"
```

Initialising:
```go
func init() {
	config.Init(config.Config{
		File:       ".subprocess",
		EnvPrefix:  "subprocess",
		Name:       version.Application(),
		Release:    version.Release(),
		FromConfig: true,
	})
}
```

Adding new config items (in another file):
```go
func init() {
	config.AddConfigItems([]string{"http.cache.url"})
}

[...]

func (s *Service) FromConfig() *Service {
	config.ApplyWith("http.cache.url", func(val interface{}) {
		s.SetURL(val.(string))
	})

	return s
}
```
