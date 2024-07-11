# IMRes

> *This project is in an early experimental state, but is currently functioning and usable.*

A lightweight, dep-free library aiming for quickly measuring image size dimensions without reading the entire bytes or relying on heavy image handling libraries.

IMRes' goal is to provide a light and efficient way for extracting image dimensions from various well-known (and most) image formats by reading only the required bytes.

✅ Multi-format support<br>
✅ Dependency-free<br>
✅ Really faster

## Supported Formats

(TODO: Detailed format)

- JPEG
- PNG
- GIF
- WEBP: VP8, VP8L, and VP8X
- AVIF

## Benchmarks

```sh
$ go test -bench=. -benchmem
```

See [https://github.com/ssut/imres/wiki/Benchmarks](https://github.com/ssut/imres/wiki/Benchmarks).



## Contributing

Any and all contributions are always welcome.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.


