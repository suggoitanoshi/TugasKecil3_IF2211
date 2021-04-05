# Tanoshi
Implementasi algoritma A* dalam menentukan jalan.

## Isi
[Deskripsi](#deskripsi-umum)
[Teknologi](#teknologi-yang-digunakan)
[Menjalankan](#persiapanmenjalankan-program)

## Deskripsi Umum
Proyek ini adalah proyek yang menggunakan algoritma A* pada graf untuk mencari
jalan terpendek antar dua simpul. Proyek ini ditulis dengan bahasa golang.
Modul graf proyek ini dapat digunakan secara independen, namun implementasi
interaksi dengan pengguna dibuat terkait dengan teknologi web, menggunakan
wasm.
Proyek ini dibuat sebagai bagian dari penugasan dalam mata kuliah IF2211
Strategi Algoritma tahun 2021.

## Teknologi yang Digunakan
- [Overpass API](https://wiki.openstreetmap.org/wiki/Overpass_API)
- [Golang](https://golang.org)
  - [syscalls/js](https://golang.org/pkg/syscalls/js)
  - [Overpass-Go](https://github.com/serjvanilla/overpass)
  - [mercator](https://github.com/davvo/mercator)
- [WebAssembly](https://webassembly.org)
- [Leaflet](https://leafletjs.com)
- [Sigma js](https://sigmajs.org)

## Persiapan/Menjalankan Program
Kompilasi modul graf dapat dilakukan dimana saja, namun untuk menjalankan UI
web disarankan menggunakan OS Linux atau menggunakan msys2 di windows.

1. [Install go](https://golang.org/dl)
2. Install dependency: `go mod download`
3. Jalankan program:
  1. Dengan Linux atau msys2:
    1. Atur `WASM_EXEC`: `export WASM_EXEC="$(go env GOROOT)/misc/wasm/wasm_exec.js"`
    2. `make run`
  2. Dengan Windows dan PowerShell:
    1. buat folder out: `mkdir out`
    2. copy file `wasm_exec.js` dari GOROOT ke out: `cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" out/`
    3. compile `main.go`: `$env:GOOS="js"; $env:GOARCH="wasm"; go build -o out/main.wasm main/main.wasm`
    4. jalankan server, serve out: `go run server.go -dir=out`