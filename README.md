caply
==========

Caply is rate limiter - cheap and yet efficient way of securing your API against basic attacks or malicious overuse.
Fixed Time Window rate limiter (also known as fixed window counter) is one of many known algorithms. This algo counts a number of operations per given period and returns true when limiter has enough capacity.


How to use
----------

```go
import "github.com/arkadyb/caply"
```


### Create rate limiter
Use ```ratelimiter.NewFixedTimeWindowRateLimiter``` function to create new rate limiter. Function takes three arguments - maximum allowed number of operations, counter period and reference to ```store.Store```. 
Where ```store.Store``` is data storage used to hold number of operations for given time window.

```go
rl := caply.NewFixedTimeWindowRateLimiter(100, 1*time.Second, myStore) 
```  

### Store
You are free to introduce your own store by implementing `store.Store` interface or use one i created for ```Redis```.
If you decide to create your own, make sure all desired operations are atomic.

To create redis store, use ```NewRedisStore``` function. It takes one argument - ```redis.Pool```.

### Use rate limiter
Rate limiter exposes only one function - ```LimitExceeded(opName string)```. It takes one argument only - operation name.
It then checks its defined capacity for given operation and returns true or false as for whether operation may be run or not.

Operation name may be for example user IP address or license key.

```go
exceeded, err := rl.LimitExceeded("my-key")
if exceeded {
	// no capacity at the moment, try again later
} else {
	// all good, operation can run
}
```  

## License
 
The MIT License (MIT)
