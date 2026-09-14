# GROW Point — API Reference

Server berjalan di port `8090` (lokal) via IIS Reverse Proxy di `https://grow.astra-visteon.com`.

## Auth Endpoints (tidak butuh JWT)

| Method | Path | Body | Deskripsi |
|--------|------|------|-----------|
| POST | `/api/auth/login` | `{"npk":"...","password":"..."}` | Login, terima access_token |
| POST | `/api/auth/refresh` | — (cookie otomatis) | Rotasi refresh token |
| POST | `/api/auth/logout` | — | Hapus refresh token |

## Endpoint Karyawan (butuh JWT Bearer Token)

| Method | Path | Deskripsi |
|--------|------|-----------|
| GET | `/api/categories` | Daftar kategori aktif |
| GET | `/api/activities?category_id=X` | Daftar aktivitas (filter per kategori) |
| POST | `/api/submissions` | Buat pengajuan baru (multipart/form-data) |
| GET | `/api/submissions/me` | Riwayat pengajuan + saldo poin |
| GET | `/api/submissions/:id` | Detail pengajuan (hanya milik sendiri) |
| PUT | `/api/submissions/:id/resubmit` | Kirim ulang pengajuan REJECTED |
| GET | `/api/rewards` | Daftar reward tersedia |
| POST | `/api/rewards/redeem/:id` | Tukar reward (cek saldo + stok, race-safe) |
| GET | `/api/rewards/history` | Riwayat penukaran reward |
| GET | `/api/leaderboard` | Peringkat global (`?category_id=&department=`) |
| GET | `/api/leaderboard/my-balance` | Saldo & peringkat user saat ini |

## Endpoint Admin (butuh JWT dengan role=admin)

| Method | Path | Body | Deskripsi |
|--------|------|------|-----------|
| GET | `/api/admin/submissions/pending` | — | Antrian pengajuan |
| GET | `/api/admin/submissions/:id` | — | Detail pengajuan mana pun |
| POST | `/api/admin/submissions/:id/approve` | `{"points_override":N}` (opsional) | Approve + set poin |
| POST | `/api/admin/submissions/:id/reject` | `{"admin_notes":"..."}` (wajib) | Reject dengan alasan |
| POST | `/api/admin/categories` | `{"name":"..."}` | Buat kategori |
| PUT | `/api/admin/categories/:id` | `{"name":"..."}` | Update kategori |
| DELETE | `/api/admin/categories/:id` | — | Soft-delete kategori |
| POST | `/api/admin/activities` | `{"category_id":N,"name":"...","default_points":N,"is_custom_input":bool}` | Buat aktivitas |
| PUT | `/api/admin/activities/:id` | *(sama)* | Update aktivitas |
| DELETE | `/api/admin/activities/:id` | — | Soft-delete aktivitas |
| POST | `/api/admin/rewards` | `{"title":"...","points_required":N,"stock":N}` | Buat reward |
| PUT | `/api/admin/rewards/:id` | *(sama)* | Update reward |
| DELETE | `/api/admin/rewards/:id` | — | Soft-delete reward |
| GET | `/api/admin/rewards/redemptions` | — | Semua riwayat penukaran |
| GET | `/api/admin/users/stats` | — | Statistik poin semua user |

## Format Submission Upload

`Content-Type: multipart/form-data`

| Field | Tipe | Keterangan |
|-------|------|-----------|
| `activity_id` | integer | ID aktivitas dari master |
| `activity_date` | string | Format `YYYY-MM-DD` |
| `custom_reference` | string | Nomor SS (Innovation) / keterangan (Sport) |
| `evidence` | file[] | Maks 5 file JPG/PNG. Otomatis dioptimasi ke kualitas 80% sebelum disimpan. Jika `cwebp` tersedia di PATH server, output berupa WebP asli; jika tidak, output berupa JPEG teroptimasi. |

## Catatan Keamanan

- Access Token disimpan di memory (Pinia), **tidak** di localStorage.
- Refresh Token via `httpOnly` Cookie (`/api/auth` path).
- Login: maks 5 kali gagal per NPK per 15 menit → dikunci (`429`).
- Rate limit per IP: maks 10 request/menit di endpoint login.
- Redeem reward menggunakan SQL Server `WITH (UPDLOCK, ROWLOCK)` untuk mencegah *race condition* stok.
