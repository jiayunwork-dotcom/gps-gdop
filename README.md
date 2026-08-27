# gps-gdop

gps-gdop 是一个卫星几何精度因子（DOP）核算服务。输入为接收机概略位置的 ECEF 坐标（米）与一组卫星的 ECEF 坐标（米），后端对每颗可见星计算单位视线 e = (r_sat − r_rx)/|r_sat − r_rx|，组成 n×4 方向余弦矩阵 H（第四列为钟差常数 1），由同一法矩阵 N = HᵀH 求逆并取迹开方得到 GDOP = √(tr(N⁻¹))，PDOP 取位置三对角、TDOP 取钟差对角，HDOP/VDOP 由 ECEF 协方差位置块旋转到接收机当地东北天（ENU）后取水平/垂直方差开方；同时返回使用卫星数、被截止剔除数与 N 的条件数。边界：可见卫星少于 4 颗、某颗卫星与接收机重合、N 奇异（例如全部视线共面导致无法估钟）、仰角截止 ≥90° 或为负值，一律返回 error。仰角截止可选（默认 5°），低于截止的卫星不进 H；同一星历提高截止使参与卫星变少时 GDOP 不得下降，增加一颗不共线卫星 GDOP 不得变差。所有坐标转换共用同一套 WGS84 椭球常数：a = 6378137.0 m、f = 1/298.257223563（故 b = a(1−f) ≈ 6356752.314245 m）、e² = f(2−f)，ECEF↔大地坐标换算与 ENU 旋转矩阵都取自此常数，视线、法矩阵与 ENU 协方差旋转没有第二套互不一致的实现。

## 启动

```bash
go run . -http :8080
```

打开 http://localhost:8080 进入交互页面，或直接用 curl 调用 API。页面可加载 `example/four-good.json`（四星分散，GDOP 小）与 `example/four-poor.json`（近共面差几何，GDOP 大），点计算后列出各 DOP，SVG 天顶图用 `/api/sky` 返回的方位/仰角画点。

## API

### POST /api/dop

请求：接收机 ECEF + 卫星列表 + 可选仰角截止。

```json
{
  "receiver_ecef": { "x": -2178657.083, "y": 4388876.234, "z": 4069505.748 },
  "satellites": [
    { "id": "G01", "x": -2148953.7, "y": 4388876.2, "z": 30351114.8 },
    { "id": "G02", "x": -3398115.4, "y": 24924566.8, "z": 3416720.4 },
    { "id": "G03", "x": -10441238.1, "y": -9806172.5, "z": 3416720.4 },
    { "id": "G04", "x": -1133260.1, "y": -30125314.3, "z": 3416720.4 }
  ],
  "elevation_mask_deg": 5
}
```

响应含 `gdop`/`pdop`/`tdop`/`hdop`/`vdop`、`satellites_used`、`satellites_total`、`satellites_rejected`、`condition_number` 与 `elevation_mask_deg`。

### POST /api/sky

请求体同上，响应返回每颗星的 `azimuth_deg` 与 `elevation_deg`（ECEF 视线经同一 ENU 旋转后 el = atan2(u, √(e²+n²))）及 `used` 标记，供天顶图绘制。

### GET /api/examples

返回内置算例 `four-good` 与 `four-poor` 的原始 JSON，供页面一键加载。

## 算例

- `example/four-good.json`：四星分散（天顶 + 三颗 19.47° 仰角、120° 间隔），GDOP ≈ 2.35。
- `example/four-poor.json`：四星挤在低仰角同一方向，近共面差几何，GDOP 远大于 four-good（约 111）。

curl 验证成功路径：

```bash
curl -s -X POST http://localhost:8080/api/dop \
  -H 'Content-Type: application/json' \
  --data-binary @example/four-good.json
```

失败路径（少于 4 颗）返回 HTTP 400 与 `{"error":"..."}`：

```bash
curl -s -X POST http://localhost:8080/api/dop \
  -H 'Content-Type: application/json' \
  -d '{"receiver_ecef":{"x":0,"y":0,"z":0},"satellites":[{"id":"G01","x":20000000,"y":0,"z":0},{"id":"G02","x":0,"y":20000000,"z":0},{"id":"G03","x":0,"y":0,"z":20000000}],"elevation_mask_deg":5}'
```

## 测试

```bash
go test ./...
```

测试覆盖失败可见路径（少于 4 星、卫星与接收机重合、N 奇异、截止角非法）与交叉规则（PDOP²+TDOP²=GDOP²、HDOP²+VDOP²=PDOP²、单位视线模长为 1、增加不共线星 GDOP 不劣化、提高截止 GDOP 不下降、接收机米级平移 DOP 变化远小于星座几何差异）。

## 实现

领域内核在 `internal/`：`wgs84`（WGS84 椭球常数、ECEF↔大地坐标、ENU 旋转与协方差旋转）、`los`（单位视线、H 矩阵、仰角/方位、截止过滤、天顶视图）、`dop`（法矩阵、4×4 求逆、各 DOP 提取、条件数、求解编排与交叉规则对比分析）。`web` 包提供 `/api/*` 路由与静态页托管；`main.go` 只负责接线。不做轨道积分，不做 Kalman 滤波，不做轨迹打卡。
