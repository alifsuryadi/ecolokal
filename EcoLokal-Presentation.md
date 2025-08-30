# EcoLokal - Bank Sampah Digital

## Presentasi Project

---

## 1. Cover

# EcoLokal

## Bank Sampah Digital dengan Sistem Penjemputan Terjadwal

**Developer:** Alif Suryadi  
**Course:** Sanbercode  
**Tech Stack:** Go, Gin, PostgreSQL, JWT, Docker

---

## 2. Problem Statement

### 🚯 Permasalahan Pengelolaan Sampah

- **Rendahnya partisipasi masyarakat** dalam program bank sampah tradisional
- **Keterbatasan waktu** warga untuk datang ke TPS (Tempat Pengumpulan Sampah)
- **Pencatatan manual** yang tidak efisien dan rawan error
- **Kurangnya transparansi** dalam sistem poin dan reward
- **Kesulitan tracking** progress dan riwayat transaksi

---

## 3. Solution Overview

### 💡 Solusi: EcoLokal Digital Platform

**EcoLokal** adalah REST API untuk mengelola ekosistem Bank Sampah digital dengan fitur:

✅ **Penjemputan Terjadwal** - Request pickup dari rumah  
✅ **Multi-Role System** - Warga, Petugas, Admin  
✅ **Sistem Poin Otomatis** - Perhitungan reward real-time  
✅ **Tracking Real-time** - Monitor status penjemputan  
✅ **Riwayat Digital** - Semua transaksi tercatat

---

## 4. Target Users

### 👥 User Personas

#### 🏠 **Warga (Citizens)**

- Penduduk RT/RW yang ingin menjual sampah
- Ingin kemudahan dan fleksibilitas waktu
- Tertarik dengan sistem reward

#### 🚛 **Petugas (Field Staff)**

- Pegawai pengumpul sampah
- Butuh tool digital untuk efisiensi kerja
- Perlu tracking route dan schedule

#### 👨‍💼 **Admin (Management)**

- Pengelola bank sampah
- Butuh kontrol sistem dan analytics
- Mengelola resource dan inventory

---

## 5. Key Features

### 🌟 Fitur Utama EcoLokal

#### **Core Features:**

- **User Management** dengan 3 role berbeda
- **Katalog Jenis Sampah** dengan sistem poin per kg
- **Penjadwalan Penjemputan** yang flexible
- **Sistem Poin Reward** yang transparan
- **Riwayat Transaksi** lengkap

#### **Technical Features:**

- **JWT Authentication** & Role-based Access
- **RESTful API** architecture
- **Auto-generated Documentation** (Swagger)
- **Database Migration** & Setup automation

---

## 6. System Architecture

### 🏗️ Arsitektur Aplikasi

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   Mobile App    │    │   Web Client     │    │   Admin Panel   │
│   (Future)      │    │   (Future)       │    │   (Future)      │
└─────────┬───────┘    └─────────┬────────┘    └─────────┬───────┘
          │                      │                       │
          └──────────────────────┼───────────────────────┘
                                 │
                    ┌────────────▼─────────────┐
                    │      EcoLokal API        │
                    │   (Gin + JWT + CORS)     │
                    └────────────┬─────────────┘
                                 │
                    ┌────────────▼─────────────┐
                    │    PostgreSQL DB         │
                    │  (Users, Pickups, etc)   │
                    └──────────────────────────┘
```

#### **Clean Architecture Pattern:**

- **Domain Layer** - Business entities & rules
- **Repository Layer** - Data access abstraction
- **Usecase Layer** - Business logic implementation
- **Handler Layer** - HTTP request/response

---

## 7. Database Design

### 🗄️ Database Schema

#### **Core Tables:**

- **users** - User authentication & profile
- **waste_types** - Catalog jenis sampah & poin
- **pickup_requests** - Request penjemputan
- **pickup_items** - Detail item per pickup
- **user_points** - Akumulasi poin user
- **transactions** - History poin masuk/keluar

#### **Key Relationships:**

- Users (1:Many) Pickup Requests
- Pickup Requests (1:Many) Pickup Items
- Waste Types (1:Many) Pickup Items
- Users (1:1) User Points

---

## 8. API Endpoints

### 🔌 REST API Endpoints

#### **Authentication:**

```
POST /api/auth/register  - Register new user
POST /api/auth/login     - User login
```

#### **User Management:**

```
GET  /api/users/profile  - Get user profile
PUT  /api/users/profile  - Update profile
GET  /api/users/points   - Check points balance
```

#### **Pickup Management:**

```
POST /api/pickups        - Create pickup request
GET  /api/pickups        - Get user pickups
PUT  /api/pickups/:id/status - Update status
PUT  /api/pickups/:id/items  - Update actual weight
```

#### **Admin Functions:**

```
GET  /api/admin/pickups       - View all pickups
PUT  /api/admin/pickups/:id/assign - Assign to petugas
POST /api/admin/waste-types   - Manage waste catalog
```

---

## 9. User Journey Flow

### 🔄 Complete User Flow

#### **Phase 1: Warga Journey**

1. **Register** → Login → Get JWT Token
2. **Browse** waste types & check points
3. **Create** pickup request dengan estimasi berat
4. **Track** status penjemputan real-time
5. **Receive** poin otomatis setelah pickup selesai

#### **Phase 2: Admin Management**

1. **Review** pending pickup requests
2. **Assign** petugas ke request tertentu
3. **Monitor** overall system performance
4. **Manage** waste types & point system

#### **Phase 3: Petugas Operations**

1. **Check** jadwal pickup hari ini
2. **Update** status during pickup process
3. **Record** actual weight setelah ditimbang
4. **Complete** pickup & trigger point calculation

---

## 10. Tech Stack Deep Dive

### 🛠️ Technology Stack

#### **Backend:**

- **Go 1.21+** - Fast, concurrent, type-safe
- **Gin Web Framework** - Lightweight HTTP router
- **PostgreSQL 12+** - Robust relational database
- **JWT** - Secure stateless authentication
- **bcrypt** - Password hashing

#### **DevOps & Tools:**

- **Docker & Docker Compose** - Containerization
- **Swagger/OpenAPI** - Auto API documentation
- **Database Migrations** - Version-controlled schema
- **Makefile** - Development automation

#### **Project Structure:**

```
cmd/api/          - Application entry point
internal/domain/  - Business entities
internal/usecase/ - Business logic
internal/repository/ - Data layer
internal/delivery/http/ - HTTP handlers
pkg/utils/       - Shared utilities
```

---

## 11. Security & Best Practices

### 🔒 Security Implementation

#### **Authentication & Authorization:**

- **JWT-based authentication** dengan expiry
- **Role-based access control** (RBAC)
- **bcrypt password hashing** (cost factor 10)
- **Input validation** dengan struct tags

#### **API Security:**

- **CORS handling** untuk web clients
- **Request validation** di semua endpoints
- **SQL injection prevention** dengan prepared statements
- **Error handling** yang tidak expose internal details

#### **Development Best Practices:**

- **Clean Architecture** pattern
- **Dependency injection** untuk testability
- **Environment-based configuration**
- **Database connection pooling**

---

## 12. Performance & Scalability

### ⚡ Performance Considerations

#### **Database Optimization:**

- **Indexed columns** untuk query performance
- **Database triggers** untuk auto-update timestamps
- **Connection pooling** untuk concurrent access
- **Prepared statements** untuk query caching

#### **API Performance:**

- **Lightweight Gin framework** (minimal overhead)
- **JSON serialization** yang efisien
- **Stateless JWT** (no session storage)
- **Modular architecture** untuk horizontal scaling

#### **Scalability Features:**

- **Docker containerization** untuk easy deployment
- **Environment-based config** untuk multi-stage
- **Database migration** untuk version control
- **RESTful design** untuk microservices readiness

---

## 13. Business Impact

### 📈 Value Proposition

#### **For Citizens:**

- **Time Saving:** No need to visit TPS manually
- **Convenience:** Schedule pickup from home
- **Transparency:** Real-time tracking & point balance
- **Incentive:** Clear reward system

#### **For Environment:**

- **Increased Participation:** Easier access = more users
- **Better Tracking:** Data for environmental analytics
- **Waste Reduction:** More efficient collection
- **Community Engagement:** Gamification encourages recycling

#### **For Management:**

- **Operational Efficiency:** Digital workflow & automation
- **Cost Reduction:** Optimized route & resource allocation
- **Data-Driven Decisions:** Analytics for better planning
- **Scalability:** System can grow with community needs

---

## 14. Future Enhancements

### 🚀 Roadmap & Next Steps

#### **Phase 2 - Mobile Integration:**

- **Mobile Apps** (iOS/Android) untuk better UX
- **Push Notifications** untuk status updates
- **GPS Tracking** untuk real-time petugas location
- **Photo Upload** untuk waste verification

#### **Phase 3 - Advanced Features:**

- **AI Route Optimization** untuk pickup efficiency
- **IoT Smart Bins** dengan sensor berat
- **Blockchain** untuk transparency & trust
- **Marketplace** untuk point redemption

#### **Phase 4 - Ecosystem Expansion:**

- **Multi-location** support (multi-RT/RW)
- **B2B Integration** dengan waste management companies
- **Government Dashboard** untuk policy insights
- **Carbon Credit** integration

---

## 15. Demo & Documentation

### 🎯 Live Demo Access

#### **API Documentation:**

- **Swagger UI:** http://localhost:8080/swagger/index.html
- **Complete API specs** dengan example requests/responses
- **Try it out** feature untuk testing endpoints

#### **Quick Start:**

```bash
# Clone & Setup
git clone https://github.com/alifsuryadi/ecolokal.git
cd ecolokal

# Start dengan Docker
docker-compose up -d

# Manual Setup
make build && make run
```

#### **Testing Credentials:**

- **Admin:** admin@ecolokal.com / admin123
- **Test User:** Create via /api/auth/register

#### **Repository:** github.com/alifsuryadi/ecolokal

---

## 16. Q&A

### ❓ Questions & Discussion

#### **Technical Questions Welcome:**

- Architecture decisions & trade-offs
- Scalability & performance considerations
- Security implementation details
- Database design choices

#### **Business Questions:**

- Market potential & user adoption
- Revenue model possibilities
- Competition analysis
- Implementation challenges

#### **Future Collaboration:**

- Open source contributions
- Feature requests & suggestions
- Partnership opportunities

---

## 17. Thank You

### 🙏 Terima Kasih

## **EcoLokal**

### Revolutionizing Waste Management Through Technology

**Contact Information:**

- **GitHub:** github.com/alifsuryadi
- **Email:** alifsuryadi037@gmail.com
- **Project Repo:** github.com/alifsuryadi/ecolokal

### _"Technology for a Cleaner, Greener Community"_

**Questions?**
