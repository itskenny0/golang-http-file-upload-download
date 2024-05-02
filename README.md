# golang-http-file-upload-download

 A simple example of an HTTP upload and download in Go 

 Start with

 ```
 go run main.go
 ```

## Fork

I forked this repo from @zupzup to make it fit my use case.

Specifically, I:
- changed the output format to JSON
- changed the default port to :8979
- changed the print to correctly point out that it is running on 0.0.0.0
- I removed the file download capability
- Temp files will have a "httpup-" prefix to their name
- The program returns the location of the file in the filesystem instead of a download link
