FROM golang:1.21 as builder

WORKDIR /app

# Копируем исходники Go и компилируем
COPY server-wrapper.go .
RUN CGO_ENABLED=0 go build -o server-wrapper server-wrapper.go


FROM openjdk:8-jre

WORKDIR /minecraft

# Указываем версию Forge
ENV FORGE_VERSION=14.23.5.2860
ENV MINECRAFT_VERSION=1.12.2

# Устанавливаем wget и скачиваем Forge
RUN apt-get update && apt-get install -y wget \
    && wget -O forge-installer.jar "https://maven.minecraftforge.net/net/minecraftforge/forge/${MINECRAFT_VERSION}-${FORGE_VERSION}/forge-${MINECRAFT_VERSION}-${FORGE_VERSION}-installer.jar" \
    && java -jar forge-installer.jar --installServer \
    && rm forge-installer.jar 

COPY mods /minecraft/mods

EXPOSE 25565 8080

RUN echo "eula=true" > eula.txt
RUN echo "online-mode=false" > server.properties

COPY --from=builder app/server-wrapper ./server-wrapper

CMD ["./server-wrapper"]
