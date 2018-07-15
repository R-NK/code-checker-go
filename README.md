# code-checker-go

## Application options
|short option|long option|description|
|:--|:--|:--|
|s|status|show files encoding and EOL|
|r|replace|replace EOL with a specified argument. e.g. r=CRLF, r=LF|
|t|target|target dir (default: current dir)|
|e|exts|target file extensions (default: all)|
|o|out|output dir (default: override)|
|c|convert|convert encoding with a specified argument. e.g. c=utf8 <br>utf8, utf-16, shift-jis, etc<br>See [Encoding spec on WHATWG](https://encoding.spec.whatwg.org/#names-and-labels)|

## Example

#### show .cpp and .h files encoding and EOL in target_dir.
```
code-checker-go.exe -s -t target_dir -e cpp -e h
```

#### replace .cpp and .h files EOL to LF in target_dir and convert them to UTF-16(LE), then save to output_dir.

```
code-checker-go.exe -r=LF -t target_dir -e cpp -e h -o output_dir -c=utf-16
```