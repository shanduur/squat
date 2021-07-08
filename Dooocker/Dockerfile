# Stage used for building the app
FROM golang:latest AS build

# install build dependencies
RUN apt-get update && apt-get install -y \
    unixodbc \
    unixodbc-dev \
    git

# clone the repository
RUN cd /tmp && git clone https://github.com/shanduur/squat.git

# go to the dir containing cloned repository
WORKDIR /tmp/squat

# build app and tools
RUN make && make tools

# Stage used for running the app
FROM ubuntu:latest AS runtime

# install runtime dependencies
RUN apt-get update && apt-get install -y \
    unixodbc

# create config and data locations
RUN mkdir -p /etc/squat/config && mkdir -p /etc/squat/data

# set data and config locations
ENV CONFIG_LOCATION=/etc/squat/config/
ENV DATA_LOCATION=/etc/squat/data/

# copy executable files
COPY --from=build /tmp/squat/build/squat            /bin/squat
COPY --from=build /tmp/squat/build/gob-generator    /bin/gob-generator
COPY --from=build /tmp/squat/bin/squat-run.sh       /bin/squat-run.sh
COPY --from=build /tmp/squat/bin/data               /etc/squat/data
COPY --from=build /tmp/squat/bin/config             /etc/squat/config

# make them executable
RUN chmod +x /bin/squat /bin/gob-generator /bin/squat-run.sh

# main command of the app
CMD [ "/bin/squat-run.sh" ]
