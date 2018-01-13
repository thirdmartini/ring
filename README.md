# go-ring

Ring API and Example Code for ring.com doorbots. It is by no means complete

See: cmd/example/

This API is based on research by https://github.com/AppleTechy/Ring

## Things We Can Do
* Get Account Information
* Get Doorbot History
* Download event videos
* Listen for Doorbell Events

## Usage:

Build Example:

```$ go build github.com/thirdmartini/ring/cmd/example```

Run Example:
./example --username=<your ring.com e-mail>  --password=<your ring.com password> --save-recordings=<number of recordings to save from history> 



