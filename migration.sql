-- migration.sql

-- 1. Tambah tipe hadiah pada tabel rewards
-- Default value di set ke 'CUSTOM' agar data yang sudah ada tidak error
ALTER TABLE rewards ADD reward_type VARCHAR(50) DEFAULT 'CUSTOM' WITH VALUES;

-- 2. Tambah kolom is_locked pada tabel activity_submissions 
-- Digunakan untuk menandai ID GROW yang sedang dalam antrean klaim (agar tidak bisa dipakai double)
ALTER TABLE activity_submissions ADD is_locked BIT DEFAULT 0;

-- Catatan:
-- Untuk status 'HOLD' di activity_submissions.status dan 'PENDING_REDEMPTION' di reward_redemptions.status,
-- karena tipe data kolom status adalah VARCHAR, kita tidak perlu melakukan ALTER TABLE (tidak ada ENUM strict di SQL Server).
-- Cukup ditangani dari sisi logic aplikasi.
