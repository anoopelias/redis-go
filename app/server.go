package main

import (
	"redis-go/app/ev"
	redis "redis-go/app/redis_go"
)

func main() {
	sc := &ev.Syscalls{}
	el := ev.NewSocketEventLoop(sc)
	data := make(map[string]string)

	err := el.Run(func(sr ev.StringReader) string {
		rr := redis.NewRespReader(sr)
		cr := redis.NewCommandReader(rr)

		c, err := cr.Read()
		// TODO: Move to resp protocol
		if err != nil {
			return "-ERR " + err.Error() + "\r\n"
		}
		return c.Execute(&data) + "\r\n"
	})

	if err != nil {
		panic(err)
	}
}
