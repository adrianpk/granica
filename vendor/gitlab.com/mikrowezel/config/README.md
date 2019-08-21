# Config

A [twelve-factor](https://12factor.net/config) config manager.

# Overview (draft)

**Example**

Configured these environment variables
```shell
export MYSVC_SERVER_HOST="localhost"
export MYSVC_SERVER_PORT="8080"
export MYSVC_ENV_DEV="true"
export MYSVC_CONST_PI="3.14159"
```

A code like this

```go
package mypackage

import(
	"gitlab.com/mikrowezel/config"
)

func main() {
	printConfig()
}

func printConfig() {
	// Setup a config handler.
	cfg := config.Load("mysvc")

	// Print raw values
	values := cfg.Get()
	fmt.Println("Host: ", values["server.host"])
	fmt.Println("Port: ", values["server.port"])

	fmt.Println("----")

	// Get a specific value
	host, _ := cfg.Val("server.host")
	fmt.Println("Host: ", host)

	fmt.Println("----")

	// A default value if var does not exist
	remote := cfg.Val("server.remote", "github.com")
	fmt.Println("Remote: ", remote)

	fmt.Println("----")

	// As an int
	port := cfg.ValAsInt("server.port", 8080)
	fmt.Println("Port: ", port)

	fmt.Println("----")

	// As a float
	pi := cfg.ValAsFloat("const.pi", 3.14)
	fmt.Println("Pi: ", pi)

	fmt.Println("----")

	// With defaults
	e := cfg.ValAsFloat("const.e", 2.71)
	fmt.Println("Euler number: ", e)

	fmt.Println("----")

	// As boolean
	dev := cfg.ValAsBool("env.dev", false)
	fmt.Println("Is Dev: ", dev)
}
```

Will print the following result
```shell
Host: localhost
Port: 8080
----
Host: localhost
----
Remote: github.com
----
Port: 8080
----
Pi: 3.14159
----
Euler number: 2.71
----
Is Dev: true
```


