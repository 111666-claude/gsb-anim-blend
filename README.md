# anim-blend

动画分层混合：层内归一化、上层覆盖、缺失剪辑回退与按真实时间过渡，只用 Go 标准库。

```
go run ./cmd/anim-blend --sample=override
go run ./cmd/anim-blend --sample=normalize
go run ./cmd/anim-blend --sample=missing
go run ./cmd/anim-blend --sample=transition
go run ./cmd/anim-blend --sample=scan
go test ./...
go vet ./...
```

## 口径（README 为准）

- **层内归一化**：一层的剪辑权重先除以本层权重之和。
- **层权重**：一层的权重不超过 1；从最上层往下算，本层贡献等于
  `层权重 × 剩余比例`，算完把剩余比例减去本层贡献。上层权重到 1 时下层贡献为 0。
- **有效权重**：所有剪辑的有效权重之和恒等于 1。
- **缺失回退**：播放不存在的剪辑时改用 fallback 剪辑，`Missing` 加一，不许静默忽略。
- **过渡**：过渡在 `durMs` 毫秒内线性完成，`Tick(dtMs)` 按真实毫秒推进；
  相同累计时长、不同帧步长必须得到相同结果。
- **代价**：`Sample` 不许扫描未参与混合的剪辑，`Scanned` 不随剪辑总数线性增长。
- **规模**：10 万个剪辑、每秒 10 万次采样，单次摊还 O(层内剪辑数)。

## 输出契约（不改格式）

```
run=.. aim=..
a=.. b=.. sum=..
idle=.. missing=..
a=.. b=..
scanned=..
```

## 目录

```
anim.go              归一化、覆盖、回退与过渡
anim_test.go         单元测试
cmd/anim-blend/      命令行入口
```
