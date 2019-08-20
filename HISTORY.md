### History for developement

- 2018/08/20
    - modify the web access part by Choong Hee Park
    - add FileServer(static) in signal/[http.go](internal/signal/http.go)
    - modify demo.html to include [demo.js](static/demo.js)
```
$ sfu-server -h
Usage of sfu-server:
  -dir string
        base directory (default "static")
  -port int
        http server port (default 8080)
```

- 2019/08/19
    - read [PLI (Picture Loss Indication)](https://webrtcglossary.com/pli/) over RTCP
    - add Makefile, HISTORY.md, go.{mod.sum}
    - copy example from pion/webrtc/examples/[sfu-minimal](https://github.com/pion/webrtc/tree/master/examples/sfu-minimal)

