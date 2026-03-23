Project Structure
├── Register/       # Logic สำหรับการสมัครสมาชิก
├── config/         # การตั้งค่า Environment และ Database
├── database/       # ไฟล์การเชื่อมต่อและ Schema ฐานข้อมูล
├── login/          # Logic สำหรับการเข้าสู่ระบบ
├── middleware/     # Auth middleware และระบบตรวจสอบสิทธิ์
├── products/       # Logic การจัดการสินค้า
├── public/         # ไฟล์ Static (HTML, CSS, JS)
├── docker-compose.yaml
├── main.go         # จุดเริ่มต้นของ Application
└── go.mod          # ไฟล์จัดการ Dependencies

1.clone repository 
    git clone https://github.com/ppxlm22/market_byGolang.git
    cd market_byGolang
2.Install Dependencies
    go mod tidy
3.สร้างไฟล์ .env 
    DB_HOST=localhost
    DB_PORT=5432
    DB_USER=myuser      
    DB_PASSWORD=mypassword 
    DB_NAME=mydatabase  
    JWT_SECRET=victoriaSecretKey_12345
4.รันคำสั่ง
    docker-compose up -d
5.รันคำสั่ง
    go run .
  

