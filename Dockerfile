FROM maven:3.6.3-openjdk-14-slim AS build

COPY settings.xml /usr/share/maven/conf/

COPY pom.xml pom.xml
COPY wrg-api/pom.xml wrg-api/pom.xml
COPY wrg-model/pom.xml wrg-model/pom.xml
COPY wrg-base/pom.xml wrg-base/pom.xml

RUN mvn dependency:go-offline package -B

## copy the pom and src code to the container
COPY wrg-api/src wrg-api/src
COPY wrg-model/src wrg-model/src
COPY wrg-base/src wrg-base/src

RUN mvn install

FROM openjdk:14-ea-jdk-alpine
USER root

RUN mkdir service

COPY --from=build /wrg-base/target/ /service/

ADD https://github.com/ufoscout/docker-compose-wait/releases/download/2.5.0/wait /wait

RUN chmod +x /wait

ENV JAVA_TOOL_OPTIONS -agentlib:jdwp=transport=dt_socket,server=y,suspend=n,address=*:5005

EXPOSE 5005

CMD /wait && java --enable-preview -jar /service/wrg-base-1.0-SNAPSHOT.jar -Xdebug