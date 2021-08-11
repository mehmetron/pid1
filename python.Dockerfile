
# docker build -t python-build . -f python.Dockerfile
# docker run --name python-build -p 3000:3000 python-build
# docker exec -it python-build sh
FROM python:3.8

RUN apt-get -y update && \
    apt-get -y install nano && \
    apt-get -y install curl

#RUN curl -sSL https://raw.githubusercontent.com/python-poetry/poetry/master/install-poetry.py | python -
#ENV PATH /root/.local/bin:$PATH

RUN python -m pip install poetry

WORKDIR /app/playground


COPY main /app/

#CMD ["sleep", "10000"]
CMD [ "../main" ]