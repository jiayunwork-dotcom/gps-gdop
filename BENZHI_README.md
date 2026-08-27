# gps-gdop：Go 卫星几何精度因子核算服务，入口 go run . -http :8080 与 POST /api/dop

给定接收机概略 ECEF 坐标与一组卫星 ECEF 坐标，对每颗可见星算单位视线并组成 n×4 方向余弦矩阵 H（第四列为钟差 1），由同一法矩阵 N=HᵀH 求逆后开方得到 GDOP/PDOP/TDOP，再把位置协方差旋到当地东北天得到 HDOP/VDOP。可见星少于 4 颗、与接收机重合、N 奇异、截止角非法必须拒绝。交叉规则：PDOP²+TDOP²=GDOP²，HDOP²+VDOP²=PDOP²；提高仰角截止使参与星变少时 GDOP 不得下降；增加不共线星 GDOP 不得变差。ECEF↔大地与 ENU 旋转共用同一套 WGS84 椭球常数。

## 构建

```
go run . -http :8080
go build ./...
go test ./...
```

## API

- `POST /api/dop` — 接收机与卫星 ECEF → 各 DOP、所用星数、条件数
- `POST /api/sky` — 同上 → 方位/仰角
- `GET /api/examples` — 内置 four-good / four-poor
- `GET /api/health` — 若页面同源存活则打开控制台

Module targets Go 1.21.
