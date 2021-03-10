# Knox Webhdfs Library


## usage:

```
go get github.com/arnobroekhof/knoxwebhdfs
```


## List

```
conf := &Conf{
        Username: "some-username",
		Password:      "some-password",
		Addr: "localhost",
		Port: "9443",
		AuthType:      AuthTypeBasic,
		SSLSkipVerify: false,
	}

	c, err := NewClient(conf)
	if err != nil {
		t.Error(err)
	}

	list, err = c.List("/tmp/hello-world")
	if err != nil {
		t.Error(err)
	}
```

## Get file

```
    conf := &Conf{
        Username: "some-username",
		Password:      "some-password",
		Addr: "localhost",
		Port: "9443",
		AuthType:      AuthTypeBasic,
		SSLSkipVerify: false,
	}

	c, err := NewClient(conf)
	if err != nil {
		panic(err)
	}

	reader, err = c.Get("/tmp/hello-world.csv")
	if err != nil {
		panic(err)
	}
	defer reader.Close()
```

## Put file

```
    conf := &Conf{
        Username: "some-username",
		Password:      "some-password",
		Addr: "localhost",
		Port: "9443",
		AuthType:      AuthTypeBasic,
		SSLSkipVerify: false,
	}

	c, err := NewClient(conf)
	if err != nil {
		panic(err)
	}

    f, err := os.Open("/some/file.txt")
    if err != nil {
        panic(err)
    }
    defer f.Close()

	reader, err = c.Put("/tmp/hello-world.csv", f)
	if err != nil {
		panic(err)
	}
	defer reader.Close()
```