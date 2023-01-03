# JDownloader client in Go

This repository hosts code for [JDownloader](https://jdownloader.org/) client written in Go

### Example usage

#### Add link and start download

```go
package main

import (
	"github.com/rkosegi/jdownloader-go/jdownloader"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	c := jdownloader.NewClient("test@acme.tld", "passw0rd", logger.Sugar())
	err = c.Connect()
	if err != nil {
		panic(err)
	}
	dev, err := c.Device("my-device-name")
	if err != nil {
		panic(err)
	}
	_, err = dev.LinkGrabber().Add([]string{"http://myremoteservice/somefile.zip"},
		jdownloader.AddLinksOptionPackage("Package-Name"),
		jdownloader.AddLinksOptionAutostart(true),
		jdownloader.AddLinksOptionDestinationDir("/mnt/download"),
		)
	if err != nil {
		panic(err)
	}
	_ = c.Disconnect()
}
```

