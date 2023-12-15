FROM ubuntu:latest

RUN apt-get update && \
    apt-get install -y software-properties-common && \
    rm -rf /var/lib/apt/lists/*

RUN add-apt-repository -y ppa:libreoffice/ppa && apt-get update

RUN apt-get install -y python3-pip && apt-get install -y libreoffice 

RUN pip install unoserver

COPY ./organizze-entries.xlsx organizze-entries.xlsx

CMD ["unoserver"]