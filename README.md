Persyaratan
- Docker
- Docker Compose
- Insomnia

Cara Running
- change .env.example to .env
- docker-compose up --build

Cara Stop
- docker-compose down

Titik Akhir yang telah dikerjakan
- Login
    - password semua semua user -> 'test'
- Find All Stock
- Unactive Warehouse 
- Transfer Stock 
    - data yang di request harus sesuai seeder di database seperti warehouse id, item id, item code, etc
- Order 
- Scheduler Lock and Release Stock (no fully worked on)

Catatan 
- Tidak ada proses membuat user, disini saya menggunakan seeder 
- Tidak ada proses membuat warehouse, disini saya menggunakan seeder 
- Tidak ada proses membuat product, disini saya menggunakan seeder 
- Semua titik akhir di forward apigw