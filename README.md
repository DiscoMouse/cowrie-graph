# 📊 Cowrie-Graph

![Status: Work in Progress](https://img.shields.io/badge/status-work%20in%20progress-yellow)

A modern, secure, and performant web-based dashboard for visualizing data from the [Cowrie](https://github.com/cowrie/cowrie) SSH & Telnet honeypot.

This project is a spiritual successor to the original [kippo-graph](https://github.com/ikoniaris/kippo-graph), built from the ground up with a secure Go backend and an interactive Apache ECharts frontend.

---

## ⚠️ Project Status

**This project has a working Minimum Viable Product (MVP).** Core features for visualizing honeypot data are functional. The application is still under active development and should be considered **Alpha** quality. The backend architecture has been refactored for maintainability.

---

## 🚀 Installation & Setup

1.  **Clone the Repository:**
    ```bash
    git clone [https://github.com/DiscoMouse/cowrie-graph.git](https://github.com/DiscoMouse/cowrie-graph.git)
    cd cowrie-graph
    ```

2.  **GeoIP Database (Manual Step):** This project uses the free MaxMind GeoLite2 City database for IP geolocation. Due to licensing, you must download this yourself.
    * Sign up for a free [MaxMind GeoLite2 account](https://www.maxmind.com/en/geolite2/signup).
    * Log in and navigate to "Download Databases".
    * Download the **"GeoLite2 City"** database in the **`GeoIP2 / MMDB`** format.
    * Extract the archive and place the `GeoLite2-City.mmdb` file in the root of this project directory.

3.  **Configuration:**
    * Copy the example configuration file:
        ```bash
        cp config.example.json config.json
        ```
    * Edit `config.json` and enter your MySQL database connection details. Ensure the `geo_db_path` points to your `.mmdb` file.

4.  **Run the Application:**
    ```bash
    go run ./cmd/web
    ```
    The server will start on `http://localhost:8080`. The first time it runs, it will automatically create the `ip_intelligence` table in your database.

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