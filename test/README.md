# `/test`

Additional external test apps and test data. Feel free to structure the `/test` directory anyway you want. For bigger projects it makes sense to have a data subdirectory. For example, you can have `/test/data` or `/test/testdata` if you need Go to ignore what's in that directory. Note that Go will also ignore directories or files that begin with "." or "_", so you have more flexibility in terms of how you name your test data directory.

## http test usage
```
func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "pong")
}

//Test_Ping test ping handler
func Test_Ping(t *testing.T) {
	req, err := http.NewRequest("GET", path.Join(common.API, "ping"), nil)
	assert.NoError(t, err)
	pingRes, err := test.InvokeRequest(req,handler, path.Join(common.API, "ping"))
	assert.NoError(t, err)
	assert.True(t, pingRes.Code == 200)
}


```