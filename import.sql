-- Description: SQL script to import initial consumer data into the database.
INSERT INTO consumers (
	id, fullname, username, email, phone, address, birth_date, status
) VALUES
-- 1
(gen_random_uuid(), 'John Doe', 'johndoe', 'john.doe@example.com', '6281234567890', 'Jl. Merdeka No. 123, Jakarta', '1990-05-10', 'active'),
-- 2
(gen_random_uuid(), 'Jane Smith', 'janesmith', 'jane.smith@example.com', '6289876543210', 'Jl. Sudirman No. 45, Bandung', '1988-11-23', 'inactive'),
-- 3
(gen_random_uuid(), 'Ahmad Yusuf', 'ahmadyusuf', 'ahmad.yusuf@example.com', '6281122334455', 'Jl. Diponegoro No. 21, Surabaya', '1992-03-15', 'active'),
-- 4
(gen_random_uuid(), 'Maria Clara', 'mariaclara', 'maria.clara@example.com', '6289988776655', 'Jl. Gajah Mada No. 10, Yogyakarta', '1995-07-01', 'suspended'),
-- 5
(gen_random_uuid(), 'Budi Santoso', 'budisantoso', 'budi.santoso@example.com', '6285566778899', 'Jl. Cihampelas No. 7, Bandung', '1985-02-28', 'active'),
-- 6
(gen_random_uuid(), 'Citra Lestari', 'citralestari', 'citra.lestari@example.com', '6286655443322', 'Jl. Malioboro No. 4, Yogyakarta', '1991-12-12', 'inactive'),
-- 7
(gen_random_uuid(), 'Kevin Pratama', 'kevinpratama', 'kevin.pratama@example.com', '6281346798200', 'Jl. Asia Afrika No. 33, Jakarta', '1993-09-30', 'active'),
-- 8
(gen_random_uuid(), 'Lina Hartati', 'linahartati', 'lina.hartati@example.com', '6287723456789', 'Jl. Braga No. 55, Bandung', '1994-04-18', 'suspended'),
-- 9
(gen_random_uuid(), 'Fajar Nugroho', 'fajarnugroho', 'fajar.nugroho@example.com', '6289001122334', 'Jl. Ahmad Yani No. 9, Semarang', '1987-08-22', 'active'),
-- 10
(gen_random_uuid(), 'Sinta Dewi', 'sintadewi', 'sinta.dewi@example.com', '6283234567890', 'Jl. Riau No. 14, Medan', '1996-06-06', 'active');
