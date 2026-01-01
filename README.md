# SubProxy

A simple subscription reverse proxy service built with Go.

---

## Quick Start

### 1. Build

```bash
git clone https://github.com/awfufu/subproxy.git
cd subproxy/src
go mod tidy
go build ./cmd/subproxy
```

### 2. Edit configure

Create `config.yaml`.

```yaml
listen: ":8042"
routes:
  - name: "sub1" # http://127.0.0.1:8042/sub1
    suburl: "https://example.com/sub1"
  - name: "sub2" # http://127.0.0.1:8042/sub2
    suburl: "https://example.com/sub2"
    proxy: "http://127.0.0.1:7890" # use proxy request
```

### 3. Run

```bash
./subproxy -f config.yaml
```
