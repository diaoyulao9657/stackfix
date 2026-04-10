# stackfix

把报错丢进去，拿解释出来。

```
$ python train.py 2>&1 | stackfix

RuntimeError in train.py:42 — batch size (128) 超出了 GPU 可用内存。
PyTorch 尝试分配 3.2GB 但设备上只剩 1.1GB。

修复: 把 config.yaml 里的 batch_size 改成 32，或者先释放显存:
  torch.cuda.empty_cache()
```

支持 Python、JavaScript、Go、Rust、Java、C/C++ 的报错信息。
通用错误也行——随便丢。

## 快速开始

```sh
# 方式一：go install（推荐）
go install github.com/diaoyulao9657/stackfix@latest

# 方式二：从源码编译
git clone https://github.com/diaoyulao9657/stackfix
cd stackfix
go build -o stackfix .
mv stackfix /usr/local/bin/
```

然后配置 API key：

```sh
mkdir -p ~/.config/stackfix
cat > ~/.config/stackfix/.env << 'EOF'
BASE_URL=https://api.tokenmix.ai/v1
API_KEY=your-api-key-here
EOF
```

需要 Go 1.22+，零外部依赖。

需要一个 OpenAI 兼容 API 的 key。默认配置指向 [TokenMix.ai](https://tokenmix.ai)（155+ 模型，新用户送 $1）——改 `.env` 里的 `BASE_URL` 就能切到别的服务商。

配置从 `~/.config/stackfix/.env` 加载（也支持当前目录的 `.env`，或者直接 `export API_KEY`）。

## 用法

```sh
# 管道接 stderr
node server.js 2>&1 | stackfix

# 直接传报错文本
stackfix "TypeError: Cannot read properties of undefined"

# 从文件读
stackfix < crash.log
```

会自动识别语言，给出对应生态的修复建议。

## 配置

写在 `.env` 里或者导出为环境变量：

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `API_KEY` | *(必填)* | 你的 API key |
| `BASE_URL` | `https://api.tokenmix.ai/v1` | API 地址 |
| `MODEL` | `gpt-4o-mini` | 用哪个模型 |

## 许可证

MIT
