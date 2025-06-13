// Taken from https://michaelheap.com/golang-using-memcached/

package main

import (
	"fmt"
	"github.com/bradfitz/gomemcache/memcache"
)

func main() {
	// Connect to our memcache instance
	mc := memcache.New("127.0.0.1:11211")

	// Set some values
	mc.Set(&memcache.Item{Key: "key_one", Value: []byte("ourclass")})
	mc.Set(&memcache.Item{Key: "key_two", Value: []byte("network services")})

	// Get a single value
	val, err := mc.Get("key_one")

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("-- %s\n", val.Value)

	// Get multiple values
	it, err := mc.GetMulti([]string{"key_one", "key_two"})

	if err != nil {
		fmt.Println(err)
		return
	}

	// It's important to note here that `range` iterates in a *random* order
	for k, v := range it {
		fmt.Printf("## %s => %s\n", k, v.Value)
	}
}
/mem