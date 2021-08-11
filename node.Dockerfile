
# docker build -t node-build . -f node.Dockerfile
# docker run --name node-build -p 3000:3000 node-build
# docker exec -it node-build sh
FROM node:16

RUN apt-get -y update && \
    apt-get -y install nano && \
    apt-get -y install curl

#WORKDIR /app
#RUN npm install -g degit
#RUN degit mehmetron/vite-react my-project
#RUN cd my-project && \
#    npm install

WORKDIR /app/playground


COPY main /app/

#CMD ["sleep", "10000"]
CMD [ "../main" ]