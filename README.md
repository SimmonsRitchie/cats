# cats
A tool implemented in Go that downloads a random cat image from [The Cat API](https://thecatapi.com/).

## Installation
First install [Go](https://golang.org/doc/install).

Then, to install the binary to your current directory, run the following command:

`GOBIN="$(pwd)" GOPATH="$(mktemp -d)" go get github.com/simmonsritchie/cats`

## Usage

In the directory you installed the binary, run:

`./cats`

Options:

```
  -o string
        Output path for cat image (default "./cat.jpg")
```

## Example cats

![screenshot](./example1.jpg)

![screenshot](./example2.jpg)

## API key

By default, cats fetches data from [The Cat API](https://thecatapi.com/) without an API key. The Cat API allows requesters to make requests without one.

However, heavy users may encounter rate limiting without use of an API key. To provide one to cats, set it as an environment variable:

`export API_KEY=xxxxxxxx`

Or set API_KEY in a .env file in the same directory as the cats binary:

```
// .env file

API_KEY=xxxxxxx
```


