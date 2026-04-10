# stackfix

Pipe an error, get an explanation. That's it.

```
$ python train.py 2>&1 | stackfix

RuntimeError in train.py:42 — batch size (128) exceeds available GPU memory.
PyTorch tried to allocate 3.2GB but only 1.1GB is free on your device.

Fix: lower batch_size to 32 in config.yaml, or free memory first:
  torch.cuda.empty_cache()
```

Works with Python, JavaScript, Go, Rust, Java, and C/C++ stack traces.
Also handles generic error messages — just throw anything at it.

## Get started

```sh
# option A: go install (recommended)
go install github.com/diaoyulao9657/stackfix@latest

# option B: build from source
git clone https://github.com/diaoyulao9657/stackfix
cd stackfix
go build -o stackfix .
mv stackfix /usr/local/bin/   # or anywhere in your PATH
```

Then configure your API key:

```sh
mkdir -p ~/.config/stackfix
cat > ~/.config/stackfix/.env << 'EOF'
BASE_URL=https://api.tokenmix.ai/v1
API_KEY=your-api-key-here
EOF
```

Requires Go 1.22+. Zero external dependencies.

You need an API key from any OpenAI-compatible provider. Default config
points to [TokenMix.ai](https://tokenmix.ai) (155+ models, $1 free credit) —
change `BASE_URL` in `.env` to use a different provider.

Config is loaded from `~/.config/stackfix/.env` (or `.env` in the current
directory, or just export `API_KEY` in your shell).

## Usage

```sh
# pipe stderr from any command
node server.js 2>&1 | stackfix

# pass error text directly
stackfix "TypeError: Cannot read properties of undefined (reading 'map')"

# from a file
stackfix < crash.log
```

It auto-detects the language from the error patterns, so the explanation
and fix suggestions match the right ecosystem.

## Config

Everything goes in `.env` (or export as env vars):

| Variable | Default | What it does |
|----------|---------|-------------|
| `API_KEY` | *(required)* | Your API key |
| `BASE_URL` | `https://api.tokenmix.ai/v1` | API endpoint |
| `MODEL` | `gpt-4o-mini` | Model to use for analysis |

## How it works

1. Reads error text from stdin or command-line args
2. Detects the language (Python, JS, Go, Rust, Java, C/C++)
3. Sends it to an LLM with a focused debugging prompt
4. Streams the explanation back to your terminal

No indexing, no config files, no setup beyond an API key.

## License

MIT
