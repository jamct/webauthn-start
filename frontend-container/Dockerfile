FROM nginx:alpine
COPY ./nginx.conf /etc/nginx/conf.d/default.conf
COPY html /usr/share/nginx/html
ADD https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/css/bootstrap.min.css /usr/share/nginx/libs/css/bootstrap.min.css
ADD https://ajax.googleapis.com/ajax/libs/jquery/3.4.1/jquery.min.js /usr/share/nginx/libs/js/jquery.js
RUN chmod -R 555 /usr/share/nginx/libs/