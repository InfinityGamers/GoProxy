# 🚀 GoProxy (Archived)

A lightweight proxy server written in Go, originally designed to forward and inspect Minecraft server traffic.

> ⚠️ **Status:** Archived — this project is no longer maintained and is kept for portfolio purposes only.  
> Changes in the Minecraft protocol (notably encryption requirements) have made this implementation obsolete.

---

## 🧠 Overview

GoProxy was built as an experimental project to explore:

- Low-level network programming in Go  
- UDP proxying and packet forwarding  
- Real-time client-server communication  
- Reverse engineering and understanding the Minecraft protocol (pre-encryption)

At the time of development, the proxy successfully intercepted and forwarded traffic between a Minecraft client and server.

---

## ✨ Features (at the time)

- UDP proxy for Minecraft traffic  
- Real-time packet forwarding  
- Lightweight and minimal dependencies  
- Clean and readable Go codebase  

---

## ⚠️ Limitations

This project no longer works with modern Minecraft servers due to:

- Introduction of encrypted communication  
- Changes in handshake and authentication mechanisms  

Supporting current versions would require:

- Implementing encryption handling  
- Updating protocol compatibility  
- Supporting modern authentication flows  

---

## 📦 Installation (Historical)

```bash
git clone https://github.com/InfinityGamers/GoProxy.git
cd GoProxy
go build -o goproxy
./goproxy
```

---

## 🎮 Example Use Cases

This project was used as a foundation for experimenting with gameplay automation and client-side behavior modification through network-level control.

---

## 🎥 Demonstrations

### ⚔️ Combat Automation (Killaura)
[![Combat Automation Demo](https://img.youtube.com/vi/w4J-HzICmWc/0.jpg)](https://www.youtube.com/watch?v=w4J-HzICmWc)

---

### ⚡ Movement & Combat Enhancements
[![Movement & Combat Demo](https://img.youtube.com/vi/h43zbAWX7zQ/0.jpg)](https://www.youtube.com/watch?v=h43zbAWX7zQ)

---

### 🤖 Automated PvP Bot
[![PvP Bot Demo](https://img.youtube.com/vi/2NqTOBOu0-Y/0.jpg)](https://www.youtube.com/watch?v=2NqTOBOu0-Y)

---

## 🧪 Learning Outcomes

This project demonstrates:

- Building concurrent network services in Go  
- Working with raw UDP sockets and streams  
- Inspecting and manipulating protocol-level data  
- Designing simple but effective proxy architectures  

---

## 💡 Note for Reviewers

This project is included as part of my portfolio to showcase early experience with:

- Systems programming  
- Networking and proxies  
- Real-time data handling  

While it is no longer compatible with current Minecraft versions, it reflects strong foundational knowledge in backend and network engineering.
