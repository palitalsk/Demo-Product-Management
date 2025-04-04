FROM node:18-alpine as build

WORKDIR /app

COPY vue/ ./  

RUN npm install

COPY . .

RUN npm run build

# เริ่มต้น nginx สำหรับ serving แอปพลิเคชันที่ build แล้ว
FROM nginx:alpine

# คัดลอกไฟล์ build ของ Vue ไปยัง nginx
COPY --from=build /app/dist /usr/share/nginx/html

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
