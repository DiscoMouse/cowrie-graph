# 📊 Cowrie-Graph

![Status: Work in Progress](https://img.shields.io/badge/status-work%20in%20progress-yellow)

A modern, secure, and performant web-based dashboard for visualizing data from the [Cowrie](https://github.com/cowrie/cowrie) SSH & Telnet honeypot.

This project is a spiritual successor to the original [kippo-graph](https://github.com/ikoniaris/kippo-graph), built from the ground up with a secure Go backend and an interactive Apache ECharts frontend.

---

## ⚠️ Project Status

**This project has a working Minimum Viable Product (MVP).** Core features for visualizing honeypot data are functional. The application is still under active development and should be considered **Alpha** quality. The backend architecture has been refactored for maintainability.

---

## ✅ Current Features

* **Multi-Page Dashboard:** A clean, responsive UI with separate, detailed pages for different data views.
* **Light/Dark Mode:** A theme toggle with the user's preference saved in their browser.
* **Time-Series Charts:** Interactive, zoomable charts for daily and monthly attack trends.
* **Top 10 Statistics:** A dedicated page with bar charts for:
    * Most common passwords
    * Most common usernames
    * Top attacking IP addresses
    * Most common SSH client versions
* **Secure API Backend:** All data is served via a secure, API-driven backend written in Go using the Gin framework.
* **Organized Codebase:** The backend logic is organized into a clean, package-based architecture.

---

## 🗺️ Roadmap (Planned Features)

* **Data Enrichment:** A service to add more context to attacker IPs.
    * [ ] Local, cached GeoIP and ASN lookups.
    * [ ] TOR exit node detection.
* **Geospatial Analysis:**
    * [ ] World map visualization of attacker origins.
* **Security & Operations:**
    * [ ] User authentication to protect the dashboard.
    * [ ] Full containerization with Docker for easy deployment.
* **Experimental:**
    * [ ] A Terminal User Interface (TUI) as an alternative frontend.

---

## 🛠️ Technology Stack

* **Backend:** Go (Gin)
* **Database:** MySQL
* **Frontend:** Apache ECharts, HTML5, CSS3

---

## 🙏 Inspiration & Credit

This project would not exist without the pioneering work of the original [kippo-graph](https://github.com/ikoniaris/kippo-graph) team. Their tool was the standard for visualizing Kippo and early Cowrie data for many years, and this project aims to build upon their legacy using modern technologies and security practices.