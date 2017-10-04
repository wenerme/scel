# scel
Sougou scel dict - 搜狗 scel 词库工具

## Tips

```bash
# Build js
yarn build
```

## Dev

```bash
# Generate pb
protoc --go_out=plugins=grpc,import_path=telattr:$HOME/go/src/ *.proto
```