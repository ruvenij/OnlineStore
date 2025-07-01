# **Simple E-commerce Backend (Mini Platform)**

Scenario:
Design and implement a simple e-commerce backend with basic product and order functionality.

### **Requirements:**

- API Endpoints:
POST /products, GET /products
POST /orders, GET /orders/:id
- Track users, products, and orders (use in-memory or SQLite)
- Basic input validation and error handling
Use Go, with proper project structure (cmd/, internal/, pkg/)
- Dockerize the service

### **Bonus Points:**

- Add authentication (JWT-based)
- Support order status transitions (placed → confirmed → shipped)
- Add metrics (Prometheus-compatible /metrics endpoint)
