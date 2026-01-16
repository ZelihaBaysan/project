# Port Scanner

A lightweight **network port scanner** written in **Go (Golang)** that checks for open TCP ports on a given host or IP address. It is designed to be simple, dependency-free, and easy to use for basic network discovery and security assessment tasks.

## ⚡ Features

* 🚀 Fast and lightweight — implemented with Go’s standard library
* 🔍 Scans a range of TCP ports for open/closed status
* 📌 No external dependencies required
* 📦 Cross-platform support (Linux, macOS, Windows)
* 💻 Easy to use from the command line

## 🧠 How It Works

This tool attempts to establish a TCP connection to each specified port on a target host. If the connection succeeds, the port is considered *open*. If it fails or times out, the port is considered *closed*.

## 📥 Installation

**Prerequisites:**
✔ Go installed (version 1.18+ recommended)

Clone the repository and build:

```bash
git clone https://github.com/ZelihaBaysan/port-scanner.git
cd port-scanner
go build -o port-scanner
```

> This produces a binary named `port-scanner` that you can run from your terminal.

## ▶️ Usage

```
./port-scanner <target> [startPort] [endPort]
```

### Examples

🔎 Scan the most common ports (1–1024):

```bash
./port-scanner scanme.sh 1 1024
```

📍 Scan a specific port:

```bash
./port-scanner 192.168.1.1 80 80
```

👉 If no port range is provided, a default range (e.g., 1–1024) may be used depending on implementation.

## 🛠️ Command-Line Parameters

| Parameter   | Description                  |
| ----------- | ---------------------------- |
| `target`    | Hostname or IP to scan       |
| `startPort` | First TCP port in scan range |
| `endPort`   | Last TCP port in scan range  |

## ⚠️ Legal and Ethical Notice

Port scanning can be intrusive. Scan only hosts and networks you have permission to test. Unauthorized scanning may be illegal in some jurisdictions.

## 🧾 License

This project is released under the **MIT License** — feel free to use and modify it responsibly.

ner "GitHub - ZelihaBaysan/port-scanner"
