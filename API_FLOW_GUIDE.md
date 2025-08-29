# EcoLokal API Flow Guide

## 🎯 Overview

EcoLokal adalah sistem Bank Sampah digital yang memfasilitasi penjemputan sampah terjadwal dengan sistem poin reward. API ini mengautomasi proses pengumpulan sampah mulai dari request warga hingga pencatatan petugas.

## 👥 User Roles

### 1. **Warga**

- Penduduk yang ingin menjual sampah
- Dapat membuat request penjemputan
- Mendapatkan poin dari sampah yang dikumpulkan

### 2. **Petugas**

- Pegawai yang bertugas mengambil sampah
- Menimbang sampah aktual dan mencatat hasilnya
- Update status penjemputan

### 3. **Admin**

- Mengelola sistem secara keseluruhan
- Assign petugas ke request penjemputan
- Mengelola jenis sampah dan sistem poin

## 🔄 Complete API Flow

### **Phase 1: Setup & Authentication**

#### 1. **User Registration**

```
POST /api/auth/register
```

**Purpose**: Daftar user baru (warga/petugas)
**Body**:

```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "password123",
  "phone": "081234567890",
  "address": "Jl. Merdeka No.1",
  "role": "warga" // atau "petugas"
}
```

#### 2. **User Login**

```
POST /api/users/login
```

**Purpose**: Login dan dapatkan JWT token
**Body**:

```json
{
  "email": "john@example.com",
  "password": "password123"
}
```

**Response**: JWT token untuk authentication

---

### **Phase 2: Warga Flow (Customer Journey)**

#### 3. **Check Available Waste Types**

```
GET /api/waste-types/active
```

**Purpose**: Lihat jenis sampah yang bisa dijual dan poinnya
**Response**: List waste types dengan point per kg

#### 4. **Create Pickup Request**

```
POST /api/pickups
Header: Authorization: Bearer <token>
```

**Purpose**: Warga membuat request penjemputan sampah
**Body**:

```json
{
  "scheduled_date": "2025-08-30",
  "scheduled_time": "09:00",
  "notes": "Sampah di depan rumah",
  "items": [
    {
      "waste_type_id": 1,
      "estimated_weight": 5.5
    },
    {
      "waste_type_id": 2,
      "estimated_weight": 3.0
    }
  ]
}
```

#### 5. **Track My Pickups**

```
GET /api/pickups/my
Header: Authorization: Bearer <token>
```

**Purpose**: Warga melihat status request penjemputannya

#### 6. **Check Points Balance**

```
GET /api/users/points
Header: Authorization: Bearer <token>
```

**Purpose**: Cek total poin yang dimiliki

#### 7. **View Transaction History**

```
GET /api/transactions/my
Header: Authorization: Bearer <token>
```

**Purpose**: Lihat riwayat poin masuk/keluar

---

### **Phase 3: Admin Flow (Management)**

#### 8. **View Pending Pickups**

```
GET /api/pickups/pending
Header: Authorization: Bearer <token>
```

**Purpose**: Admin melihat semua request yang belum di-assign

#### 9. **Get Available Petugas**

```
GET /api/users/role/petugas
Header: Authorization: Bearer <token>
```

**Purpose**: Lihat list petugas untuk assignment

#### 10. **Assign Pickup to Petugas**

```
PUT /api/pickups/{id}/assign
Header: Authorization: Bearer <token>
```

**Purpose**: Admin assign request ke petugas tertentu
**Body**:

```json
{
  "petugas_id": 3
}
```

#### 11. **Manage Waste Types**

```
POST /api/waste-types          // Create new waste type
PUT /api/waste-types/{id}      // Update existing
DELETE /api/waste-types/{id}   // Delete waste type
```

**Purpose**: Admin kelola jenis sampah dan poin

---

### **Phase 4: Petugas Flow (Field Operations)**

#### 12. **Check Today's Pickups**

```
GET /api/pickups/petugas?date=2025-08-30
Header: Authorization: Bearer <token>
```

**Purpose**: Petugas lihat jadwal penjemputan hari ini

#### 13. **Update Pickup Status**

```
PUT /api/pickups/{id}/status
Header: Authorization: Bearer <token>
```

**Purpose**: Update status (scheduled → in_progress → completed)
**Body**:

```json
{
  "status": "in_progress" // atau "completed"
}
```

#### 14. **Record Actual Weight**

```
PUT /api/pickups/{id}/items
Header: Authorization: Bearer <token>
```

**Purpose**: Petugas catat berat aktual setelah menimbang
**Body**:

```json
{
  "items": [
    {
      "id": 1,
      "actual_weight": 4.8 // berat sebenarnya
    },
    {
      "id": 2,
      "actual_weight": 2.5
    }
  ]
}
```

#### 15. **Complete Pickup**

```
PUT /api/pickups/{id}/status
```

**Purpose**: Tandai pickup selesai
**Body**:

```json
{
  "status": "completed"
}
```

---

### **Phase 5: Point Management**

#### 16. **Auto Point Calculation**

**Automatic Process**: Sistem otomatis hitung poin berdasarkan:

- Actual weight × point per kg
- Update user_points table
- Create transaction record

#### 17. **Manual Point Transaction** (Admin)

```
POST /api/transactions
Header: Authorization: Bearer <token>
```

**Purpose**: Admin buat transaksi poin manual (penukaran reward)
**Body**:

```json
{
  "user_id": 2,
  "type": "redeem",
  "points": 500,
  "description": "Tukar pulsa Rp 50.000"
}
```

---

## 🎯 Business Value & Use Cases

### **For Warga (Citizens)**

- **Convenience**: Request pickup dari rumah tanpa perlu ke TPS
- **Transparency**: Real-time tracking status penjemputan
- **Rewards**: Sistem poin yang bisa ditukar dengan benefit
- **Eco-friendly**: Mendorong daur ulang sampah

### **For Petugas (Staff)**

- **Route Optimization**: Lihat semua pickup dalam 1 area/hari
- **Digital Recording**: Catat berat digital, tidak perlu kertas
- **Performance Tracking**: History pekerjaan tersimpan
- **Efficient Workflow**: Update status real-time

### **for Admin (Management)**

- **Resource Planning**: Assign petugas berdasarkan lokasi/kapasitas
- **Business Intelligence**: Data analytics untuk optimasi
- **Inventory Management**: Kelola jenis sampah dan harga
- **Financial Control**: Monitor transaksi poin dan reward

### **For Environment**

- **Waste Reduction**: Lebih banyak sampah ter-recycle
- **Community Engagement**: Gamifikasi mendorong partisipasi
- **Data-driven Policy**: Data untuk kebijakan lingkungan
- **Sustainable Economy**: Ekonomi sirkuler berbasis teknologi

## 📊 Success Metrics

### **Operational Metrics**

- Pickup completion rate
- Average response time
- User satisfaction score
- Waste collection volume

### **Business Metrics**

- Active users growth
- Revenue from waste sales
- Cost reduction vs traditional method
- ROI from point reward system

### **Environmental Impact**

- Total waste diverted from landfill
- CO2 reduction from recycling
- Community recycling participation rate
- Waste stream optimization efficiency

## 🚀 Scalability Features

### **Current Capabilities**

- Multi-role user management
- Real-time status tracking
- Automated point calculation
- RESTful API architecture

### **Future Enhancements**

- Mobile app integration
- GPS tracking for petugas
- AI-powered route optimization
- IoT integration for smart bins
- Blockchain for transparency
- Marketplace for point redemption

## 🔒 Security & Compliance

### **Authentication & Authorization**

- JWT-based authentication
- Role-based access control (RBAC)
- API rate limiting
- Input validation & sanitization

### **Data Protection**

- Encrypted password storage (bcrypt)
- HTTPS communication
- Database connection security
- User data privacy compliance

---

## 📱 Integration Possibilities

### **Mobile App**

- Native iOS/Android app
- Push notifications for status updates
- Camera integration for waste photos
- GPS for location services

### **Payment Gateway**

- Digital wallet integration
- Bank transfer for point redemption
- E-commerce marketplace connection
- Cryptocurrency rewards (future)

### **Third-party Services**

- Google Maps for routing
- SMS gateway for notifications
- Email service for communications
- Analytics platforms for insights

---
