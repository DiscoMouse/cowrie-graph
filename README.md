# 📊 Cowrie-Graph

![Status: Work in Progress](https://img.shields.io/badge/status-work%20in%20progress-yellow)

A modern, secure, and performant web-based dashboard for visualizing data from the [Cowrie](https://github.com/cowrie/cowrie) SSH & Telnet honeypot.

This project is a spiritual successor to the original `kippo-graph`, built from the ground up with a secure Go backend and an interactive Apache ECharts frontend.

---

## ⚠️ Project Status

**This project is in the very early stages of development and is NOT ready for any form of use.** The repository currently contains planning documents and initial scaffolding only. Please check back later for a functional release.

---

## ✨ Planned Features

- [ ] Interactive, time-series graphs for attack trends (per hour, per day).
- [ ] "Top 10" lists for most common usernames, passwords, attacking IPs, and countries.
- [ ] World map visualization of attacker origins.
- [ ] SSH client version fingerprinting and analysis.
- [ ] Local, cached GeoIP and ASN lookups using the free GeoLite2 database.
- [ ] TOR exit node detection.
- [ ] A secure, API-driven architecture.
- [ ] (Experimental) A Terminal User Interface (TUI) as a potential alternative frontend.

---

## 🛠️ Technology Stack

* **Backend:** Go
* **Database:** MySQL
* **Frontend:** Apache ECharts, HTML5, CSS3
* **Deployment:** Docker (planned)

---

## 🙏 Inspiration & Credit

This project would not exist without the pioneering work of the original `kippo-graph` team. Their tool was the standard for visualizing Kippo and early Cowrie data for many years, and this project aims to build upon their legacy using modern technologies and security practices.