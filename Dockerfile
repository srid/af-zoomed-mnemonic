FROM golang:1.4-onbuild

# For use with jwilder/nginx-proxy
ENV VIRTUAL_HOST afmnemonic.happyandharmless.com

EXPOSE 8080
