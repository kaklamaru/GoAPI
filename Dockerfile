# Start from the latest golang base image
FROM golang:latest

WORKDIR /app

# คัดลอกไฟล์ go.mod และ go.sum ไปยัง container
COPY go.mod go.sum ./

# ติดตั้ง dependencies
RUN go mod tidy

# คัดลอกโค้ดทั้งหมดไปยัง container
COPY . .

# สร้างแอปพลิเคชัน Go
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# สร้าง image ที่จะใช้รัน
FROM scratch

# คัดลอกไฟล์ executable ไปยัง container
COPY --from=build /app/main /main

# เปิดพอร์ตที่แอปพลิเคชันจะรัน
EXPOSE 8080

# คำสั่งที่ให้ container รันเมื่อเริ่ม
CMD ["/main"]