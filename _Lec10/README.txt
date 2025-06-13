
$ go run lookuphost.go go.dev

compare to

$ dig go.dev A go.dev AAAA +short

$ go run lookupport.go tcp telnet
$ go run lookupport.go tcp http
