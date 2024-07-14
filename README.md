# IMRes

> *This project is in an early experimental state, but is currently functioning and usable.*

A lightweight, dep-free library aiming for quickly measuring image size dimensions without reading the entire bytes or relying on heavy image handling libraries.

IMRes' goal is to provide a light and efficient way for extracting image dimensions from various well-known (and most) image formats by reading only the required bytes.

✅ Multi-format support<br>
✅ Dependency-free<br>
✅ Really faster

## How to use

```sh
$ go get github.com/ssut/imres
```

## Supported Formats

- JPEG
- PNG
- GIF
- WEBP: VP8, VP8L, and VP8X
- AVIF
- HEIF (very experimental): I must note that HEIF is not thoroughly tested due to the nature of HEIF support in the real world, some challenges come from obtaining sample images (probably due to the license requirements for the codec.)
- TIF (TIFF)
- BMP

## Benchmarks

IMRes is faster than the native Go image library

The Go library `golang.org/x/image` does not support some image formats such as AVIF, TIFF. IMRes not only supports these formats but is also way faster for the formats it supports. Still I'm willing to make the performance even better than now though.



```sh
$ go test ./tests/... -bench=. -benchmem
```

See [https://github.com/ssut/imres/wiki/Benchmarks](https://github.com/ssut/imres/wiki/Benchmarks).



## Contributing

Any and all contributions are always welcome.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.


