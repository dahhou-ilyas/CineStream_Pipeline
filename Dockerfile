FROM selenium/standalone-chrome:latest

EXPOSE 4444 7900

CMD ["selenium-server"]