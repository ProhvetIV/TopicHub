FROM node:18.17.1-alpine AS build-stage

WORKDIR /frontend

COPY package*.json ./
RUN npm install 

COPY . .
RUN npm run build

FROM nginx:alpine AS production-stage
COPY --from=build-stage /frontend/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 5173
EXPOSE 8080
EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
