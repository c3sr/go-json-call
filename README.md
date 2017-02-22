# go-json-call

Library that allows arbitrary Go functions to be called with JSON-formatted arguments and return values. For example:

let `args` be `{"0": 6, "1": 3}`

```
func div(x,y int) int {
    return x / y
}

results := CallWithJSON(div, args)
```

results is then `{"0": 2}`

The JSON is indexed by the number of the argument and result.